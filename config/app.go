package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	defaultGatewayRPCURL = "http://localhost:8949"
)

// AppConfig defines the configuration file structure
type AppConfig struct {
	GatewayRPCClientURL        string            `json:"gateway_rpc_client_url"`
	ChainServiceManagerAddress map[uint64]string `json:"chain_service_manager_address"`
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
	if ac.GatewayRPCClientURL == "" {
		ac.GatewayRPCClientURL = defaultGatewayRPCURL
	}

	if ac.ChainServiceManagerAddress == nil {
		ac.ChainServiceManagerAddress = make(map[uint64]string)
	}
}

// Validate checks if the configuration is valid
func (ac *AppConfig) Validate() error {
	if ac.GatewayRPCClientURL == "" {
		return fmt.Errorf("GatewayRPCClientURL is required")
	}
	if len(ac.ChainServiceManagerAddress) == 0 {
		return fmt.Errorf("ChainServiceManagerAddress is required")
	}
	return nil
}
