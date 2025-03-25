package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cosmossdk.io/log"
	csquaringmanager "github.com/0xPellNetwork/dvs-contracts-template/bindings/IncredibleSquaringServiceManager"
	"github.com/0xPellNetwork/pellapp-sdk/service/tx"
	interactorconfig "github.com/0xPellNetwork/pelldvs-interactor/config"
	rpclocal "github.com/0xPellNetwork/pelldvs/rpc/client/local"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"

	"github.com/0xPellNetwork/dvs-template/app"
	chainconnector "github.com/0xPellNetwork/dvs-template/chain_connector"
	dvsappcfg "github.com/0xPellNetwork/dvs-template/config"
	sqserver "github.com/0xPellNetwork/dvs-template/dvs/squared/server"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

// StartOperatorCmd defines the command to start the Operator
var StartOperatorCmd = &cobra.Command{
	Use:   "start-operator",
	Short: "Start the Operator",
	Long:  "Start the task operator service Example:\n squaringd start-operator --home",
	RunE:  startOperator,
}

func startOperator(cmd *cobra.Command, args []string) error {
	return runOperator(cmd)
}

func runOperator(cmd *cobra.Command) error {
	serverCtx := server.GetServerContextFromCmd(cmd)
	squaringConfig, err := dvsappcfg.LoadAppConfig(fmt.Sprintf("%s/%s", config.RootDir, "config/squaring.config.json"))
	if err != nil {
		logger.Error("Failed to load squaring config", "error", err)
		return fmt.Errorf("failed to load squaring config: %w", err)
	}

	sqserver.ChainConnector, err = chainconnector.NewClient(squaringConfig.GatewayRPCClientURL)
	if err != nil {
		logger.Error("Failed to create Chain Connector client", "error", err)
		return fmt.Errorf("failed to create Chain Connector client: %v", err)
	}

	node := app.NewApp(codectypes.NewInterfaceRegistry(), logger, config, squaringConfig)

	td, err := NewTaskDispatcher(
		serverCtx.Logger,
		node.DvsNode.GetLocalClient(),
		config.Pell.InteractorConfigPath,
		squaringConfig.ChainServiceManagerAddress,
	)
	if err != nil {
		return err
	}

	if err := td.Start(); err != nil {
		logger.Error("Failed to start task dispatcher", "error", err)
		return fmt.Errorf("failed to start TaskDispatcher: %w", err)
	}

	// Start node goroutine
	go func() {
		if err = node.Start(); err != nil {
			logger.Error("failed to start Node", "error", err.Error())
		}
	}()

	logger.Info("All components started successfully")

	// Block main thread for App
	select {}
}

// ChainWatcher watches a specific chain for new tasks
type ChainWatcher struct {
	chainID        int64
	rpcURL         string
	wsURL          string
	serviceManager *csquaringmanager.ContractIncredibleSquaringServiceManager
	taskChan       chan *csquaringmanager.ContractIncredibleSquaringServiceManagerNewTaskCreated
}

// TaskDispatcher manages multiple chain watchers and dispatches tasks
type TaskDispatcher struct {
	chains        []*ChainWatcher
	pellDVSClient *rpclocal.Local
	msgEncoder    tx.MsgEncoder
	logger        log.Logger
}

// NewTaskDispatcher creates a new task dispatcher
func NewTaskDispatcher(logger log.Logger, pellDVSClient *rpclocal.Local, interacotrCfgPath string, chainServiceManagerAddress map[int64]string) (*TaskDispatcher, error) {
	var interacotrConfig = interactorconfig.Config{}

	configBytes, err := os.ReadFile(interacotrCfgPath)
	if err != nil {
		logger.Error("Failed to read chain config", "error", err)
		return nil, fmt.Errorf("failed to read chain config: %w", err)
	}

	if err := json.Unmarshal(configBytes, &interacotrConfig); err != nil {
		logger.Error("Failed to parse chain config", "error", err)
		return nil, fmt.Errorf("failed to parse chain config: %w", err)
	}

	var chains []*ChainWatcher

	for chainID, detail := range interacotrConfig.ContractConfig.DVSConfigs {
		serviceManagerAddr, has := chainServiceManagerAddress[int64(chainID)]
		if !has {
			logger.Error("Chain service manager address not found", "chainID", chainID)
			continue
		}

		wsCLient, err := ethclient.Dial(detail.WSURL)
		if err != nil {
			logger.Error("Failed to connect to Ethereum client", "error", err)
			return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
		}

		serviceManager, err := csquaringmanager.NewContractIncredibleSquaringServiceManager(common.HexToAddress(serviceManagerAddr), wsCLient)
		if err != nil {
			logger.Error("Failed to create contract filter", "error", err)
			return nil, fmt.Errorf("failed to create contract filter: %w", err)
		}

		chains = append(chains, &ChainWatcher{
			chainID:        int64(chainID),
			rpcURL:         detail.RPCURL,
			wsURL:          detail.WSURL,
			serviceManager: serviceManager,
			taskChan:       make(chan *csquaringmanager.ContractIncredibleSquaringServiceManagerNewTaskCreated),
		})
	}

	return &TaskDispatcher{
		logger:        logger,
		chains:        chains,
		pellDVSClient: pellDVSClient,
		msgEncoder:    tx.NewDefaultDecoder(codec.NewProtoCodec(codectypes.NewInterfaceRegistry())),
	}, nil
}

// Start starts the task dispatcher
func (td *TaskDispatcher) Start() error {
	for _, chain := range td.chains {
		go td.listenForNewTasks(chain)
	}
	return nil
}

// listenForNewTasks listens for new tasks on a specific chain
func (td *TaskDispatcher) listenForNewTasks(chain *ChainWatcher) {
	td.logger.Info("start listening for new tasks",
		"chainID", chain.chainID,
		"serviceManager", chain.serviceManager,
		"rpcURL", chain.rpcURL,
		"wsURL", chain.wsURL,
	)
	filterOpts := &bind.WatchOpts{}

	taskChan := make(chan *csquaringmanager.ContractIncredibleSquaringServiceManagerNewTaskCreated, 1000)
	sub, err := chain.serviceManager.WatchNewTaskCreated(filterOpts, taskChan, nil)
	if err != nil {
		td.logger.Error("Failed to create task listener", "error", err)
		td.logger.Error("Failed to create task listener", "error", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case newTask := <-taskChan:
			td.logger.Info("New task created", "taskID", newTask.TaskIndex, "chainID", chain.chainID)
			groupNumbers := make([]uint32, len(newTask.Task.GroupNumbers))
			for i, b := range newTask.Task.GroupNumbers {
				groupNumbers[i] = uint32(b)
			}

			taskData, err := td.serializeTask(uint64(chain.chainID), newTask)
			if err != nil {
				td.logger.Error("Failed to serialize task", "chainID", chain.chainID, "error", err)
				return
			}

			_, err = td.pellDVSClient.RequestDVS(
				context.Background(),
				taskData,
				int64(newTask.Task.TaskCreatedBlock),
				chain.chainID,
				groupNumbers,
				[]uint32{newTask.Task.GroupThresholdPercentage},
			)
			if err != nil {
				td.logger.Error("Failed to send task", "error", err)
			}

		case err := <-sub.Err():
			td.logger.Error("Task monitoring error", "error", err)
			return
		}
	}
}

func (td *TaskDispatcher) serializeTask(chainID uint64, newTask *csquaringmanager.ContractIncredibleSquaringServiceManagerNewTaskCreated) ([]byte, error) {
	td.logger.Info("serializeTask",
		"chainID", chainID,
		"taskIndex", newTask.TaskIndex,
		"task", fmt.Sprintf("%+v", newTask.Task),
	)

	task := newTask.Task
	taskRequest := &types.RequestNumberSquaredIn{
		Task: &types.TaskRequest{
			TaskIndex:                newTask.TaskIndex,
			Height:                   task.TaskCreatedBlock,
			ChainId:                  chainID,
			Squared:                  task.NumberToBeSquared.String(),
			GroupNumbers:             task.GroupNumbers,
			GroupThresholdPercentage: task.GroupThresholdPercentage,
		},
	}

	return td.msgEncoder.EncodeMsgs(taskRequest)
}
