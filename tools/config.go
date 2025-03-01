package tools

import (
	"encoding/json"
	"math/big"
	"os"
)

// Config defines the configuration file structure
type TaskGatewayConfig struct {
	TaskGatewayAddress        string   `json:"task_gateway_address"`
	ServiceManagerAddress     string   `json:"service_manager_address"`
	TaskGatewayPrivateKeyPath string   `json:"gateway_key_path"`
	RPCURL                    string   `json:"rpc_url"`
	ChainID                   *big.Int `json:"chain_id"`
}

// LoadConfig loads configuration from the specified path
func LoadTaskGatewayConfig(path string) (*TaskGatewayConfig, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config TaskGatewayConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SquaringConfig defines the configuration file structure
type SquaringConfig struct {
	GatewayRPCClientURL        string           `json:"gateway_rpc_client_url"`
	ServiceManagerAddress      string           `json:"service_manager_address"`
	ChainServiceManagerAddress map[int64]string `json:"chain_service_manager_address"`
}

// LoadSquaringConfig loads configuration from the specified path
func LoadSquaringConfig(path string) (*SquaringConfig, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config SquaringConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
