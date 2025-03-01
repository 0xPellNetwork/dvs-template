package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/0xPellNetwork/dvs-template/chain_connector"
)

// StartTaskGatewayCmd defines the command to start the TaskGateway
var StartTaskGatewayCmd = &cobra.Command{
	Use:   "start-chain-connector",
	Short: "Start the chain connector",
	Long:  "Start the chain connector service Example:\n squaringd start-chain-connector --home=",
	RunE:  startChainConnector,
}

func init() {
	StartTaskGatewayCmd.Flags().String("home", "", "home_dir")
}

func startChainConnector(cmd *cobra.Command, args []string) error {
	return runChainConnector(cmd)
}

func runChainConnector(cmd *cobra.Command) error {
	logger.Info("Starting Chain Connector", "home", config.RootDir)
	// Create RPC server
	rpcServer, err := chain_connector.NewServer(fmt.Sprintf("%s/%s", config.RootDir, "config/task_gateway.config.json"))
	if err != nil {
		logger.Error("Failed to create Task Gateway RPC server", "error", err)
		return fmt.Errorf("Failed to create Task Gateway RPC server: %v", err)
	}

	// Start server
	if err := rpcServer.Start(); err != nil {
		logger.Error("Failed to start Task Gateway RPC server", "error", err)
		return fmt.Errorf("Failed to start Task Gateway RPC server: %v", err)
	}

	logger.Info("Task Gateway RPC server started successfully", "address", rpcServer.Config().ServerAddr)

	// Block main thread for TaskGateway
	select {}

}
