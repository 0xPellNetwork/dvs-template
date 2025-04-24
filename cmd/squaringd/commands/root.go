package commands

import (
	"fmt"
	"os"

	"github.com/0xPellNetwork/pellapp-sdk/client"
	clienthelpers "github.com/0xPellNetwork/pellapp-sdk/client/helpers"
	"github.com/0xPellNetwork/pellapp-sdk/server"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	dvsconfig "github.com/0xPellNetwork/pelldvs/config"
	"github.com/0xPellNetwork/pelldvs/libs/cli"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logger = log.NewLogger(os.Stdout)
	config = dvsconfig.DefaultConfig()
)

var DefaultNodeHome string

func init() {
	var err error
	DefaultNodeHome, err = clienthelpers.GetNodeHomeDirectory(".pelldvs")
	if err != nil {
		panic(err)
	}
}

// RootCmd creates the root command for the application
func RootCmd() *cobra.Command {
	initClientCtx := client.Context{}.
		WithHomeDir(DefaultNodeHome).
		WithInterfaceRegistry(codectypes.NewInterfaceRegistry())

	cmd := &cobra.Command{
		Use:   "squared",
		Short: "Square number application developed based on PellDvs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			config = loadNodeConfig(viper.GetString("home"))
			if viper.GetBool(cli.TraceFlag) {
				logger = log.NewTracingLogger(logger)
			}
			logger = logger.With("module", "main")

			initClientCtx = initClientCtx.WithCmdContext(cmd.Context())

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customPellDVSConfig := initPellDVSConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customPellDVSConfig)
		},
	}

	return cmd
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
