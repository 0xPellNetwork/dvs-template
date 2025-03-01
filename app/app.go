package app

import (
	"github.com/0xPellNetwork/pellapp-sdk/baseapp"
	dvsservermanager "github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler"
	"github.com/0xPellNetwork/pellapp-sdk/pelldvs"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	"github.com/0xPellNetwork/pelldvs/config"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdktypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	dvs "github.com/0xPellNetwork/dvs-template/dvs/squared"
	msg_server "github.com/0xPellNetwork/dvs-template/dvs/squared/msg_server"
	dvstypes "github.com/0xPellNetwork/dvs-template/dvs/squared/types"
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
	logger            log.Logger
	appCodec          codec.Codec
	interfaceRegistry codectypes.InterfaceRegistry

	DvsNode   *pelldvs.Node
	DvsServer msg_server.Server
}

// NewApp initializes a new App instance
func NewApp(
	interfaceRegistry codectypes.InterfaceRegistry,
	logger log.Logger,
	cmtcfg *config.Config,
	gatewayRPCClientURL string,
) *App {
	var app = &App{
		BaseApp:           baseapp.NewBaseApp(logger),
		interfaceRegistry: interfaceRegistry,
		logger:            logger,
		appCodec:          codec.NewProtoCodec(interfaceRegistry),
	}

	// Register standard types interfaces
	std.RegisterInterfaces(app.interfaceRegistry)
	sdktypes.RegisterInterfaces(app.interfaceRegistry)

	var err error

	// Build dvs node
	app.DvsNode, err = pelldvs.NewNode(app.logger, app, cmtcfg)
	if err != nil {
		panic(err)
	}

	// Initialize DVS server
	app.DvsServer, err = msg_server.NewServer(app.logger, gatewayRPCClientURL)
	if err != nil {
		panic(err)
	}

	// Initialize DVS server manager
	dvsservermanager.InitDvsMsgHelper(app.appCodec)

	// Register DVS services
	dvs.NewAppModule(app.DvsServer).RegisterServices()
	dvstypes.RegisterInterfaces(app.interfaceRegistry)

	return app
}

// Start method starts the application
func (app *App) Start() error {
	app.logger.Info("App Start")
	if err := app.DvsNode.Start(); err != nil {
		app.logger.Error("DvsNode Start Failed", "error", err.Error())
		return err
	}

	// Block the main thread
	c := make(chan any)
	<-c
	return nil
}
