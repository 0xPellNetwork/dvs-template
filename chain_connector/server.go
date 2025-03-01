package chain_connector

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/rpc"
	"os"
	"sync"

	csquaringManager "github.com/0xPellNetwork/dvs-contracts-template/bindings/IncredibleSquaringServiceManager"
	"github.com/0xPellNetwork/pelldvs-interactor/chainlibs/eth"
	"github.com/0xPellNetwork/pelldvs/libs/log"
	"github.com/0xPellNetwork/pelldvs/libs/service"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))

// Config configuration structure
type Config struct {
	TaskGatewayPrivateKeyStorePath string                 `json:"gateway_key_path"`
	ServerAddr                     string                 `json:"server_addr"`
	Chains                         map[uint64]ChainConfig `json:"chains"`
}

type ChainConfig struct {
	RPCURL          string `json:"rpc_url"`
	ContractAddress string `json:"contract_address"`
	ChainID         uint64 `json:"chain_id"`
	GasLimit        uint64 `json:"gas_limit"`
}

// Server wraps TaskGateway as RPC server
type Server struct {
	service.BaseService
	config   *Config
	server   *rpc.Server
	listener net.Listener
	address  string
	mu       sync.Mutex
	Chains   map[uint64]*ChainBinding
}

type ChainBinding struct {
	ChainID         uint64
	ContractAddress common.Address
	Client          eth.Client
	ServiceManager  *csquaringManager.ContractIncredibleSquaringServiceManager
	Auth            *bind.TransactOpts
}

// TaskRequest represents RPC request structure
type TaskRequest struct {
	ChainID                     uint64
	Task                        csquaringManager.IIncredibleSquaringServiceManagerTask
	TaskResponse                csquaringManager.IIncredibleSquaringServiceManagerTaskResponse
	NonSignerStakesAndSignature csquaringManager.IBLSSignatureVerifierNonSignerStakesAndSignature
}

// TaskResponse represents RPC response structure
type TaskResponse struct {
	Error string
}

// NewServer creates a new TaskGateway RPC server
func NewServer(configPath string) (*Server, error) {
	// Read config file
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		logger.Error("Failed to read config file", "error", err)
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		logger.Error("Failed to parse config", "error", err)
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	// Read private key
	keyJSON, err := os.ReadFile(config.TaskGatewayPrivateKeyStorePath)
	if err != nil {
		logger.Error("Failed to read private key file", "error", err)
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	key, err := keystore.DecryptKey(keyJSON, "")
	if err != nil {
		logger.Error("Failed to decrypt private key", "error", err)
		return nil, fmt.Errorf("failed to decrypt private key: %v", err)
	}

	chainBindings := make(map[uint64]*ChainBinding)
	for chainID, chainCfg := range config.Chains {
		logger.Info("Creating chain binding", "chainID", chainID, "chainCfg", chainCfg)
		ethClient, err := eth.NewClient(chainCfg.RPCURL)
		if err != nil {
			logger.Error("Failed to create Ethereum client", "chainID", chainID, "error", err)
			return nil, fmt.Errorf("failed to create Ethereum client for chainID: %d, error: %v", chainID, err)
		}
		contractAddress := common.HexToAddress(chainCfg.ContractAddress)
		serviceManager, err := csquaringManager.NewContractIncredibleSquaringServiceManager(contractAddress, ethClient)
		if err != nil {
			logger.Error("Failed to create service manager contract", "chainID", chainID, "error", err)
			return nil, fmt.Errorf("failed to create service manager contract for chainID: %d, error: %v", chainID, err)
		}

		// Create transaction authenticator
		auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, big.NewInt(int64(chainID)))
		if err != nil {
			logger.Error("Failed to create transaction authenticator", "error", err)
			return nil, fmt.Errorf("failed to create transaction authenticator: %v", err)
		}

		// set gas limit to force the transaction to be broadcasted
		if chainCfg.GasLimit > 0 {
			auth.GasLimit = chainCfg.GasLimit
		}

		chainBindings[chainID] = &ChainBinding{
			ChainID:         chainID,
			Client:          ethClient,
			ContractAddress: contractAddress,
			ServiceManager:  serviceManager,
			Auth:            auth,
		}
	}

	server := rpc.NewServer()
	ts := &Server{
		config:  &config,
		server:  server,
		address: config.ServerAddr,
		Chains:  chainBindings,
	}

	if err := server.Register(ts); err != nil {
		logger.Error("Failed to register RPC server", "error", err)
		return nil, fmt.Errorf("failed to register RPC server: %v", err)
	}

	ts.BaseService = *service.NewBaseService(nil, "TaskGatewayServer", ts)
	return ts, nil
}

// OnStart starts the RPC server
func (ts *Server) OnStart() error {
	listener, err := net.Listen("tcp", ts.address)
	if err != nil {
		logger.Error("Failed to start listener", "address", ts.address, "error", err)
		return fmt.Errorf("failed to start listener on %s: %v", ts.address, err)
	}

	ts.listener = listener
	logger.Info("Task Gateway RPC server started", "address", ts.address)

	go ts.server.Accept(listener)
	return nil
}

// OnStop stops the RPC server
func (ts *Server) OnStop() {
	if ts.listener != nil {
		logger.Info("Stopping TaskGateway RPC server")
		ts.listener.Close()
	}
}

var taskResponseState = map[uint64]map[uint32]int{}

// RespondToTask handles task response RPC method
func (ts *Server) RespondToTask(req *TaskRequest, resp *TaskResponse) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	logger.Info("Handling task response request",
		"taskID", req.TaskResponse.ReferenceTaskIndex,
		"req", req,
		"resp", resp,
	)
	if taskResponseState[req.ChainID] == nil {
		taskResponseState[req.ChainID] = make(map[uint32]int)
	}

	if _, ok := taskResponseState[req.ChainID][req.TaskResponse.ReferenceTaskIndex]; ok {
		logger.Info("Task response already processed", "taskID", req.TaskResponse.ReferenceTaskIndex)
		resp.Error = ""
		return nil
	}

	_, err := ts.sendResponseToTask(req)
	if err != nil {
		logger.Error("Failed to respond to task", "error", err)
		resp.Error = err.Error()
		return err
	}

	taskResponseState[req.ChainID][req.TaskResponse.ReferenceTaskIndex] = 1

	logger.Info("Successfully responded to task", "taskID", req.TaskResponse.ReferenceTaskIndex)
	resp.Error = ""
	return nil
}

// Config returns the server configuration
func (ts *Server) Config() *Config {
	return ts.config
}

// sendResponseToTask sends response to task,
// if error contains chainErrorInvalidReferenceBLock, it will retry after timeWaitForResendTaskResponse
func (ts *Server) sendResponseToTask(req *TaskRequest) (*gethtypes.Transaction, error) {
	chainBinding := ts.Chains[req.ChainID]
	if chainBinding == nil {
		return nil, errors.New(fmt.Sprintf("chainBinding not found for chainID: %d", req.ChainID))
	}
	serviceManager := chainBinding.ServiceManager
	if serviceManager == nil {
		return nil, errors.New(fmt.Sprintf("service manager not found for chainID: %d", req.ChainID))
	}
	tx, err := serviceManager.RespondToTask(chainBinding.Auth, req.Task, req.TaskResponse, req.NonSignerStakesAndSignature)
	if err == nil {
		return tx, nil
	}

	logger.Error("Failed to respond to task@sendResponseToTask",
		"error", err.Error(),
	)
	return nil, err
}
