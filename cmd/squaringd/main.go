package main

import (
	"os"

	"github.com/0xPellNetwork/pelldvs/libs/cli"

	"github.com/0xPellNetwork/dvs-template/cmd/squaringd/commands"
)

func main() {
	rootCmd := commands.RootCmd
	rootCmd.AddCommand(
		commands.StartOperatorCmd,
		commands.StartAggregatorCmd,
	)
	cmd := cli.PrepareBaseCmd(rootCmd, "", os.ExpandEnv("$HOME"))
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
