package dvs

import (
	storetypes "cosmossdk.io/store/types"
	sdkservice "github.com/0xPellNetwork/pellapp-sdk/service"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"google.golang.org/grpc"

	resulthandlers "github.com/0xPellNetwork/dvs-template/dvs/squared/result"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/server"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
	apptypes "github.com/0xPellNetwork/dvs-template/types"
)

// AppModule implements an application module for the dvs module.
type AppModule struct {
	logger      log.Logger
	server      *server.Server
	queryServer *server.Querier

	txMgr    apptypes.TxManager
	queryMgr apptypes.DataManager
}

// NewAppModule creates a new AppModule object
func NewAppModule(logger log.Logger,
	storeKey storetypes.StoreKey,
	txMgr apptypes.TxManager,
	queryMgr apptypes.DataManager,
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

// RegisterServices registers module services.
func (am *AppModule) RegisterServices(router *sdkservice.MsgRouter) {
	configurator := router.GetConfigurator()
	// register dvs-msg handler server
	types.RegisterSquaredMsgServerServer(configurator, am.server)
	types.RegisterQueryServer(configurator, am.queryServer)

	// register dvs-msg result handler
	configurator.RegisterResultMsgExtractor(
		&types.RequestNumberSquaredOut{}, resulthandlers.NewResultHandler(),
	)
}

func (am *AppModule) RegisterGRPCServer(srv *grpc.Server) {
	// Register the module's gRPC server
	//types.RegisterSquaredMsgServerServer(srv, am.server)
	types.RegisterQueryServer(srv, am.queryServer)

}

func (am *AppModule) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}
