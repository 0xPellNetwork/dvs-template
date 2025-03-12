package server

import (
	"fmt"

	"github.com/0xPellNetwork/pelldvs-libs/log"

	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

var _ types.SquaredMsgServerServer = Server{}

// Server struct represents the server with a logger and a chain connector client.
type Server struct {
	logger log.Logger // Logger for logging messages.
}

// NewServer creates a new Server instance with the provided logger and gateway RPC client URL.
func NewServer(
	logger log.Logger, // Logger for the server.
) (Server, error) {
	return Server{
		logger: logger, // Initialize the server with the provided logger.
	}, nil // Return the initialized server instance.
}

// Logger returns a module-specific logger.
func (k Server) Logger() log.Logger {
	// Add module-specific information to the logger.
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
