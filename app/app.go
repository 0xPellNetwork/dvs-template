package app

import (
	"github.com/0xPellNetwork/pellapp-sdk/baseapp"
	dvsservermanager "github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler"
	"github.com/0xPellNetwork/pellapp-sdk/pelldvs"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	"github.com/0xPellNetwork/pelldvs/config"
	rpclocal "github.com/0xPellNetwork/pelldvs/rpc/client/local"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdktypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"

	dvs "github.com/0xPellNetwork/dvs-template/dvs/squared"
	dvsserver "github.com/0xPellNetwork/dvs-template/dvs/squared/server"
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

	appCodec          codec.Codec
	interfaceRegistry codectypes.InterfaceRegistry

	dvsNode *pelldvs.Node

	DvsServer                dvsserver.Server
	ProcessRequestServer     grpc1.Server
	PostProcessRequestServer grpc1.Server

	logger log.Logger

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
	app.dvsNode, err = pelldvs.NewNode(app.logger, app, cmtcfg)
	if err != nil {
		panic(err)
	}

	// Initialize DVS client
	app.DVSClient = app.dvsNode.GetLocalClient()

	// Initialize DVS server
	app.DvsServer, err = dvsserver.NewServer(app.logger, gatewayRPCClientURL)
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
