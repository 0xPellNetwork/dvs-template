package server

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"

	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
	apptypes "github.com/0xPellNetwork/dvs-template/types"
)

var _ types.SquaredMsgServerServer = &Server{}

// Server struct represents the server with a logger and a chain connector client.
type Server struct {
	logger   log.Logger               // Logger for logging messages.
	app      apptypes.AppCommitStorer // Application commit store reference.
	storeKey storetypes.StoreKey
}

// NewServer creates a new Server instance with the provided logger and gateway RPC client URL.
func NewServer(
	logger log.Logger, // Logger for the server.
	storeKey storetypes.StoreKey,
) (*Server, error) {
	return &Server{
		logger:   logger, // Initialize the server with the provided logger.
		storeKey: storeKey,
	}, nil // Return the initialized server instance.
}

// Logger returns a module-specific logger.
func (s *Server) Logger() log.Logger {
	// Add module-specific information to the logger.
	return s.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetApp sets the app reference for the server.
func (s *Server) SetApp(app apptypes.AppCommitStorer) {
	s.logger.Info("SetApp", "app", app) // Log the app being set.
	s.app = app                         // Set the app reference for the server.
}
