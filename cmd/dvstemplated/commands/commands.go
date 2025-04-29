package commands

import (
	serverconfig "github.com/0xPellNetwork/pellapp-sdk/server/config"
	pelldvscfg "github.com/0xPellNetwork/pelldvs/config"
)

// initPellDVSConfig helps to override default CometBFT Config values.
// return cmtcfg.DefaultConfig if no custom configuration is required for the application.
func initPellDVSConfig() *pelldvscfg.Config {
	cfg := pelldvscfg.DefaultConfig()

	// these values put a higher strain on node memory
	// cfg.P2P.MaxNumInboundPeers = 100
	// cfg.P2P.MaxNumOutboundPeers = 40

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// The following code snippet is just for reference.

	// DvsTemplatedConfig defines an arbitrary custom config to extend app.toml.
	// If you don't need it, you can remove it.
	// If you wish to add fields that correspond to flags that aren't in the SDK server config,
	// this custom config can as well help.
	type DvsTemplatedConfig struct {
		Filepath string `mapstructure:"filepath"`
	}

	type CustomAppConfig struct {
		serverconfig.Config `mapstructure:",squash"`

		DvsTemplated DvsTemplatedConfig `mapstructure:"dvstemplated"`
	}

	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()

	// Now we set the custom config default values.
	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
		DvsTemplated: DvsTemplatedConfig{
			Filepath: "anything",
		},
	}

	// The default SDK app template is defined in serverconfig.DefaultConfigTemplate.
	// We append the custom config template to the default one.
	// And we set the default config to the custom app template.
	customAppTemplate := serverconfig.DefaultConfigTemplate + `
[dvstemplated]
# That field will be parsed by server.InterceptConfigsPreRunHandler and held by viper.
# Do not forget to add quotes around the value if it is a string.
filepath = "{{ .DvsTemplated.Filepath }}"

# end for custom config
`

	return customAppTemplate, customAppConfig
}
