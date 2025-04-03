package commands

import (
	"fmt"

	"github.com/0xPellNetwork/pellapp-sdk/server"
	"github.com/spf13/cobra"

	chainconnector "github.com/0xPellNetwork/dvs-template/chain_connector"
	dvsappcfg "github.com/0xPellNetwork/dvs-template/config"
	sqserver "github.com/0xPellNetwork/dvs-template/dvs/squared/server"
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
	squaringConfig, err := dvsappcfg.LoadAppConfig(config.RootDir + "/config/squaring.config.json")
	if err != nil {
		logger.Error("Failed to load squaring config", "error", err)
		return fmt.Errorf("failed to load squaring config: %w", err)
	}

	logger.Info("Starting Operator", "squaringConfig", fmt.Sprintf("%+v", squaringConfig))

	sqserver.ChainConnector, err = chainconnector.NewClient(squaringConfig.GatewayRPCClientURL)
	if err != nil {
		logger.Error("Failed to create Chain Connector client", "error", err)
		return fmt.Errorf("failed to create Chain Connector client: %v", err)
	}

	td, err := NewTaskDispatcher(
		serverCtx.Logger,
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
	//go func() {
	//	if err = node.Start(); err != nil {
	//		logger.Error("failed to start Node", "error", err.Error())
	//	}
	//}()

	logger.Info("All components started successfully")

	// Block main thread for App
	select {}
}
