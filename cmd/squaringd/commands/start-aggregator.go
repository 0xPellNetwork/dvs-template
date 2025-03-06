package commands

import (
	"context"
	"fmt"

	"github.com/0xPellNetwork/pelldvs/aggregator"
	"github.com/0xPellNetwork/pelldvs/aggregator/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagRPCAddress = "address"
	flagTimeout    = "timeout"
)

// StartAggregatorCmd defines the command to start the aggregator
var StartAggregatorCmd = &cobra.Command{
	Use:   "start-aggregator",
	Short: "Start the aggregator",
	RunE:  startAggregator,
}

func init() {
	StartAggregatorCmd.Flags().String(flagRPCAddress, "0.0.0.0:26653", "RPC server listen address")
	StartAggregatorCmd.Flags().String(flagTimeout, "3s", "Aggregation operation timeout")
}

// startAggregator implements the logic to start the aggregator
func startAggregator(cmd *cobra.Command, args []string) error {
	return runAggregatorService(cmd)
}

func runAggregatorService(cmd *cobra.Command) error {
	rpcAddress := viper.GetString(flagRPCAddress)
	timeout := viper.GetString(flagTimeout)

	aggregatorConfig, err := aggregator.LoadConfig(config.RootDir + "/config/aggregator.json")
	if err != nil {
		return fmt.Errorf("failed to load aggregator configuration: %v", err)
	}

	if rpcAddress != "" {
		aggregatorConfig.AggregatorRPCServer = rpcAddress
	}
	if timeout != "" {
		aggregatorConfig.OperatorResponseTimeout = timeout
	}

	rpcAggregator, err := rpc.NewRPCServerAggregator(context.Background(), config, aggregatorConfig, logger)
	if err != nil {
		return fmt.Errorf("failed to create RPCAggregator: %v", err)
	}

	if err = rpcAggregator.Start(); err != nil {
		return fmt.Errorf("failed to start aggregator: %v", err)
	}

	logger.Info("Aggregator started", "RPC address", rpcAddress)

	// Keep the service running until an interrupt signal is received
	select {}
}
