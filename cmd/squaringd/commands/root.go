package commands

import (
	"os"

	"github.com/0xPellNetwork/pelldvs-libs/log"
	"github.com/0xPellNetwork/pelldvs/libs/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/0xPellNetwork/dvs-template/tools"
)

var (
	logger = log.NewLogger(os.Stdout)
	config = tools.GenerateDefaultNodeConfig()
)

// RootCmd is the root command for squaringd server.
var RootCmd = &cobra.Command{
	Use:   "squaringd",
	Short: "Square number application developed based on PellDvs",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// Initialize config with home path
		config = tools.GenerateNodeConfig(viper.GetString("home"))
		if viper.GetBool(cli.TraceFlag) {
			logger = log.NewTracingLogger(logger)
		}
		logger = logger.With("module", "main")
		return nil
	},
}
