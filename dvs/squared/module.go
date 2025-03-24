package dvs

import (
	storetypes "cosmossdk.io/store/types"
	sdkservice "github.com/0xPellNetwork/pellapp-sdk/service"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	resulthandlers "github.com/0xPellNetwork/dvs-template/dvs/squared/result"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/server"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
	apptypes "github.com/0xPellNetwork/dvs-template/types"
)

// AppModule implements an application module for the dvs module.
type AppModule struct {
	logger log.Logger
	server *server.Server
	app    apptypes.AppCommitStorer
}

// NewAppModule creates a new AppModule object
func NewAppModule(logger log.Logger, storeKey storetypes.StoreKey) *AppModule {
	s, err := server.NewServer(logger, storeKey)
	if err != nil {
		panic(err)
	}

	return &AppModule{
		logger: logger.With("module", types.ModuleName),
		server: s,
	}
}

// RegisterServices registers module services.
func (am *AppModule) RegisterServices(router *sdkservice.MsgRouter) {
	configurator := router.GetConfigurator()
	// register dvs-msg handler server
	types.RegisterSquaredMsgServerServer(configurator, am.server)

	// register dvs-msg result handler
	configurator.RegisterResultMsgExtractor(
		&types.RequestNumberSquaredOut{}, resulthandlers.NewResultHandler(),
	)
}

func (am *AppModule) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (am *AppModule) SetAppCommitStore(app apptypes.AppCommitStorer) {
	am.logger.Info("SetAppCommitStore", "app", app)
	am.app = app
	am.server.SetApp(app)
}
