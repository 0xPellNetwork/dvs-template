package commands

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/spf13/cobra"

	"github.com/0xPellNetwork/dvs-template/app"
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
	node := app.NewApp(codectypes.NewInterfaceRegistry(), logger, config)

	// Start node goroutine
	go func() {
		if err := node.Start(); err != nil {
			logger.Error("failed to start Node", "error", err.Error())
		}
	}()

	logger.Info("All components started successfully")

	// Block main thread for App
	select {}
}
