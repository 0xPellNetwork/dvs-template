package server

import (
	"fmt"

	"github.com/0xPellNetwork/pelldvs-libs/log"

	chainConnector "github.com/0xPellNetwork/dvs-template/chain_connector"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

var _ types.SquaredMsgServerServer = Server{}

// Server struct represents the server with a logger and a chain connector client.
type Server struct {
	logger log.Logger // Logger for logging messages.
	// NOTE: chain connector client. In this case, we need to write the computed results to the chain,
	// so this connector is required. However, it is not mandatory and depends on the specific
	// business requirements of dvs.
	tg *chainConnector.Client // Chain connector client for interacting with the blockchain.
}

// NewServer creates a new Server instance with the provided logger and gateway RPC client URL.
func NewServer(
	logger log.Logger, // Logger for the server.
	gatewayRPCClientURL string, // URL for the gateway RPC client.
) (Server, error) {
	k := Server{
		logger: logger, // Initialize the server with the provided logger.
	}
	var err error
	// Create a new chain connector client using the provided URL.
	k.tg, err = chainConnector.NewClient(gatewayRPCClientURL)
	if err != nil {
		// Log an error message if the client creation fails.
		logger.Error("Failed to create Chain Connector client", "error", err)
		return k, fmt.Errorf("Failed to create Chain Connector client :%v", err)
	}

	return k, nil // Return the initialized server instance.
}

// Logger returns a module-specific logger.
func (k *Server) Logger() log.Logger {
	// Add module-specific information to the logger.
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
