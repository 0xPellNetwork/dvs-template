package dvs

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	resulthandlers "github.com/0xPellNetwork/dvs-template/dvs/squared/result"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/server"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

var _ sdktypes.AppModule = &AppModule{}

var _ sdktypes.BasicModule = &AppModule{}

var _ sdktypes.MsgResultExtractor = &AppModule{}

// AppModule implements an application module for the dvs module.
type AppModule struct {
	logger      log.Logger
	server      *server.Server
	queryServer *server.Querier

	txMgr    sdktypes.TxManager
	queryMgr sdktypes.QueryManager
}

// NewAppModule creates a new AppModule object
func NewAppModule(logger log.Logger,
	storeKey storetypes.StoreKey,
	txMgr sdktypes.TxManager,
	queryMgr sdktypes.QueryManager,
) *AppModule {
	s, err := server.NewServer(logger, storeKey, txMgr)
	if err != nil {
		panic(err)
	}

	qs, err := server.NewQuerier(logger, storeKey, queryMgr)
	if err != nil {
		panic(err)
	}

	return &AppModule{
		logger:      logger.With("module", types.ModuleName),
		server:      s,
		queryServer: qs,
		txMgr:       txMgr,
		queryMgr:    queryMgr,
	}
}

func (this *AppModule) Name() string {
	return types.ModuleName
}

func (this *AppModule) IsAppModule() {}

// RegisterServices registers module services.
func (am *AppModule) RegisterServices(configurator sdktypes.Configurator) {
	//configurator := router.GetConfigurator()
	// register dvs-msg handler server
	types.RegisterSquaredMsgServerServer(configurator, am.server)
	types.RegisterQueryServer(configurator, am.queryServer)
}

func (am *AppModule) RegisterGRPCServer(srv *grpc.Server) {
	// Register the module's gRPC server
	types.RegisterSquaredMsgServerServer(srv, am.server)
	types.RegisterQueryServer(srv, am.queryServer)

}

func (am *AppModule) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the bank module.
func (am *AppModule) RegisterGRPCGatewayRoutes(clientCtx gogogrpc.ClientConn, mux *gwruntime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (am *AppModule) RegisterResultMsgExtractors(configurator sdktypes.Configurator) {
	// register dvs-msg result handler
	configurator.RegisterResultMsgExtractor(
		&types.RequestNumberSquaredOut{}, resulthandlers.NewResultHandler(),
	)
}

// RegisterQueryServices
func (am *AppModule) RegisterQueryServices(router gogogrpc.Server) {
	// Register the module's query server
	types.RegisterQueryServer(router, am.queryServer)
}
