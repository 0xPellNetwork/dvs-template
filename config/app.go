package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	defaultQueryHTTPServerAddr = "0.0.0.0:8123"
	defaultQueryGRPCServerAddr = "0.0.0.0:9123"
	defaultQueryStoreKey       = "query"
	defaultGatewayRPCURL       = "http://localhost:8949"
)

// AppConfig defines the configuration file structure
type AppConfig struct {
	QueryRPCServerAddress      string           `json:"query_rpc_server_address"`
	QueryHTTPServerAddress     string           `json:"query_http_server_address"`
	GatewayRPCClientURL        string           `json:"gateway_rpc_client_url"`
	ChainServiceManagerAddress map[int64]string `json:"chain_service_manager_address"`
	QueryStoreKey              string           `json:"query_store_key"`
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

	config.Finalize()

	return &config, nil
}

// Finnalize initializes the configuration
func (ac *AppConfig) Finalize() {
	if ac.QueryRPCServerAddress == "" {
		ac.QueryRPCServerAddress = defaultQueryGRPCServerAddr
	}
	if ac.QueryHTTPServerAddress == "" {
		ac.QueryHTTPServerAddress = defaultQueryHTTPServerAddr
	}
	if ac.GatewayRPCClientURL == "" {
		ac.GatewayRPCClientURL = defaultGatewayRPCURL
	}

	if ac.QueryStoreKey == "" {
		ac.QueryStoreKey = defaultQueryStoreKey
	}

	if ac.ChainServiceManagerAddress == nil {
		ac.ChainServiceManagerAddress = make(map[int64]string)
	}
}

// Validate checks if the configuration is valid
func (ac *AppConfig) Validate() error {
	if ac.QueryRPCServerAddress == "" {
		return fmt.Errorf("QueryRPCServerAddress is required")
	}
	if ac.QueryHTTPServerAddress == "" {
		return fmt.Errorf("QueryHTTPServerAddress is required")
	}
	if ac.GatewayRPCClientURL == "" {
		return fmt.Errorf("GatewayRPCClientURL is required")
	}
	if len(ac.ChainServiceManagerAddress) == 0 {
		return fmt.Errorf("ChainServiceManagerAddress is required")
	}
	return nil
}
