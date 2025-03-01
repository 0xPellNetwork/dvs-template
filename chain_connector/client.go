package chain_connector

import (
	"fmt"
	"net/rpc"

	csquaringManager "github.com/0xPellNetwork/dvs-contracts-template/bindings/IncredibleSquaringServiceManager"
)

// Client represents RPC client
type Client struct {
	client *rpc.Client
}

// NewClient creates a new TaskGateway RPC client
func NewClient(address string) (*Client, error) {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		logger.Error("Failed to connect to RPC server", "error", err)
		return nil, fmt.Errorf("failed to connect to RPC server: %v", err)
	}

	logger.Info("Connected to RPC server", "address", address)
	return &Client{client: client}, nil
}

// RespondToTask handles client RPC method call
func (tc *Client) RespondToTask(
	chainID uint64,
	task csquaringManager.IIncredibleSquaringServiceManagerTask,
	taskResponse csquaringManager.IIncredibleSquaringServiceManagerTaskResponse,
	nonSignerStakesAndSignature csquaringManager.IBLSSignatureVerifierNonSignerStakesAndSignature,
) error {
	req := &TaskRequest{
		ChainID:                     chainID,
		Task:                        task,
		TaskResponse:                taskResponse,
		NonSignerStakesAndSignature: nonSignerStakesAndSignature,
	}
	resp := &TaskResponse{}

	logger.Info("Sending task response request", "taskID", taskResponse.ReferenceTaskIndex)

	err := tc.client.Call("Server.RespondToTask", req, resp)
	if err != nil {
		logger.Error("RPC call failed", "error", err)
		return err
	}

	if resp.Error != "" {
		logger.Error("task RespondToTask failed", "error", resp.Error)
		return fmt.Errorf("task RespondToTask failed: %s", resp.Error)
	}

	logger.Info("Task response sent successfully", "taskID", taskResponse.ReferenceTaskIndex)
	return nil
}
