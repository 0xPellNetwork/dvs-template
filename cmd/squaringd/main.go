package main

import (
	svrcmd "github.com/0xPellNetwork/pellapp-sdk/server/cmd"

	"github.com/0xPellNetwork/dvs-template/cmd/squaringd/commands"
)

func main() {
	rootCmd := commands.RootCmd
	rootCmd.AddCommand(
		commands.StartTaskGatewayCmd,
		commands.StartOperatorCmd,
		commands.StartAggregatorCmd,
		commands.StartAggregatorCmd,
	)
	commands.InitRunOperatorCommand(rootCmd)
	if err := svrcmd.Execute(rootCmd, "", commands.DefaultNodeHome); err != nil {
		panic(err)
	}
}
