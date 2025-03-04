package app

import (
	"github.com/0xPellNetwork/pellapp-sdk/baseapp"
	"github.com/0xPellNetwork/pellapp-sdk/pelldvs"
	dsm "github.com/0xPellNetwork/pellapp-sdk/service"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	"github.com/0xPellNetwork/pelldvs/config"

	rpclocal "github.com/0xPellNetwork/pelldvs/rpc/client/local"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdktypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sq "github.com/0xPellNetwork/dvs-template/dvs/squared"
	sqserver "github.com/0xPellNetwork/dvs-template/dvs/squared/server"
	sqtypes "github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

const (
	// Application name
	Name = "dvs-template"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
)

// App struct represents the application
type App struct {
	*baseapp.BaseApp

	appCodec          codec.Codec
	logger            log.Logger
	interfaceRegistry codectypes.InterfaceRegistry

	dvsNode   *pelldvs.Node
	DvsServer sqserver.Server
	DVSClient *rpclocal.Local
}

// Start method starts the application
func (app *App) Start() error {
	app.logger.Info("App Start")
	if err := app.dvsNode.Start(); err != nil {
		app.logger.Error("DvsNode Start Failed", "error", err.Error())
		return err
	}

	// Block the main thread
	c := make(chan any)
	<-c
	return nil
}

// NewApp initializes a new App instance
func NewApp(
	interfaceRegistry codectypes.InterfaceRegistry,
	logger log.Logger,
	cmtcfg *config.Config,
	gatewayRPCClientURL string,
) *App {
	cdc := codec.NewProtoCodec(interfaceRegistry)

	app := &App{
		BaseApp:           baseapp.NewBaseApp(Name, logger, cdc),
		interfaceRegistry: interfaceRegistry,
		logger:            logger,
		appCodec:          cdc,
	}

	// Register standard types interfaces
	std.RegisterInterfaces(app.interfaceRegistry)
	sdktypes.RegisterInterfaces(app.interfaceRegistry)

	var err error

	// Build dvs node
	app.dvsNode, err = pelldvs.NewNode(app.logger, app, cmtcfg)
	if err != nil {
		panic(err)
	}

	// Initialize DVS client
	app.DVSClient = app.dvsNode.GetLocalClient()

	// Initialize DVS server
	app.DvsServer, err = sqserver.NewServer(app.logger, gatewayRPCClientURL)
	if err != nil {
		panic(err)
	}

	// Initialize DVS server manager
	handler := dsm.NewDvsMsgHandlers(app.appCodec)

	// Register DVS services
	sq.NewAppModule(app.DvsServer).RegisterServices(handler.GetProcessor())
	sqtypes.RegisterInterfaces(app.interfaceRegistry)

	return app
}
