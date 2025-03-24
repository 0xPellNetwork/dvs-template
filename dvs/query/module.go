package query

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	"google.golang.org/grpc"

	"github.com/0xPellNetwork/dvs-template/dvs/query/server"
	"github.com/0xPellNetwork/dvs-template/dvs/query/types"
	apptypes "github.com/0xPellNetwork/dvs-template/types"
)

type AppModule struct {
	logger log.Logger
	server *server.Server
}

// NewAppModule creates a new AppModule instance
func NewAppModule(logger log.Logger, storeKey storetypes.StoreKey) *AppModule {
	s := server.NewServer(logger.With("module", types.ModuleName), storeKey)
	return &AppModule{
		logger: logger,
		server: s,
	}
}

// RegisterGRPCServices registers the gRPC services for the module
func (am *AppModule) RegisterGRPCServices(server *grpc.Server) {
	types.RegisterQueryServiceServer(server, am.server)
}

// SetAppQuerier sets the app querier for the server
func (am *AppModule) SetAppQuerier(app apptypes.AppQueryStorer) {
	am.server.SetApp(app)
}
