package server

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"

	"github.com/0xPellNetwork/dvs-template/dvs/query/types"
	apptypes "github.com/0xPellNetwork/dvs-template/types"
)

// make sure Server implements the QueryServiceServer interface
var _ types.QueryServiceServer = &Server{}

type Server struct {
	logger   log.Logger
	app      apptypes.AppQueryStorer
	storeKey storetypes.StoreKey
}

// NewServer creates a new Server instance
func NewServer(logger log.Logger, storeKey storetypes.StoreKey) *Server {
	return &Server{
		logger:   logger,
		storeKey: storeKey,
	}
}

// SetApp sets the app reference
func (s *Server) SetApp(app apptypes.AppQueryStorer) {
	s.app = app
}

// AppQuerier defines the interface for querying app data
type AppQuerier interface {
	GetCommitStore(key storetypes.StoreKey) storetypes.KVStore
	GetCommitMultiStore() storetypes.CommitMultiStore
}
