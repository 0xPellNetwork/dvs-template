package tools

import (
	"fmt"

	"github.com/0xPellNetwork/pelldvs/config"
	"github.com/spf13/viper"
)

func GenerateDefaultNodeConfig() *config.Config {
	return config.DefaultConfig()
}

func GenerateNodeConfig(homeFlag string) *config.Config {
	// Create default config
	cmtcfg := config.DefaultConfig()

	// If home flag is provided, set it as root directory
	if homeFlag != "" {
		cmtcfg.RootDir = homeFlag
		cmtcfg.SetRoot(homeFlag)
		config.EnsureRoot(homeFlag)
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
