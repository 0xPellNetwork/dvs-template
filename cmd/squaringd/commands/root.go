package commands

import (
	"fmt"
	"os"

	"github.com/0xPellNetwork/pelldvs-libs/log"
	dvsconfig "github.com/0xPellNetwork/pelldvs/config"
	"github.com/0xPellNetwork/pelldvs/libs/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logger = log.NewLogger(os.Stdout)
	config = dvsconfig.DefaultConfig()
)

// RootCmd is the root command for squaringd server.
var RootCmd = &cobra.Command{
	Use:   "squaringd",
	Short: "Square number application developed based on PellDvs",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// Initialize config with home path
		config = loadNodeConfig(viper.GetString("home"))
		if viper.GetBool(cli.TraceFlag) {
			logger = log.NewTracingLogger(logger)
		}
		logger = logger.With("module", "main")
		return nil
	},
}

func loadNodeConfig(homeFlag string) *dvsconfig.Config {
	// Create default config
	cmtcfg := dvsconfig.DefaultConfig()

	// If home flag is provided, set it as root directory
	if homeFlag != "" {
		cmtcfg.RootDir = homeFlag
		cmtcfg.SetRoot(homeFlag)
		dvsconfig.EnsureRoot(homeFlag)
	}

	// Set config file path
	viper.SetConfigFile(fmt.Sprintf("%s/%s", cmtcfg.RootDir, "config/config.toml"))

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %v", err))
	}

	// Parse config into struct
	if err := viper.Unmarshal(cmtcfg); err != nil {
		panic(fmt.Errorf("error parsing config file: %v", err))
	}

	return cmtcfg
}
