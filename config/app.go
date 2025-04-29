package config

import (
	"encoding/json"
	"os"
)

// AppConfig defines the configuration file structure
type AppConfig struct {
}

// LoadAppConfig loads configuration from the specified path
func LoadAppConfig(path string) (*AppConfig, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Finnalize initializes the configuration
func (ac *AppConfig) Finalize() {
}

// Validate checks if the configuration is valid
func (ac *AppConfig) Validate() error {
	return nil
}
