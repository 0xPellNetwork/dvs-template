package app

import (
	"github.com/0xPellNetwork/pellapp-sdk/baseapp"
	"github.com/0xPellNetwork/pellapp-sdk/pelldvs"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	"github.com/0xPellNetwork/pelldvs/config"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdktypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sq "github.com/0xPellNetwork/dvs-template/dvs/squared"
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

	DvsNode *pelldvs.Node
}

// NewApp initializes a new App instance
func NewApp(
	interfaceRegistry codectypes.InterfaceRegistry,
	logger log.Logger,
	cfg *config.Config,
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
	app.DvsNode, err = pelldvs.NewNode(app.logger, app, cfg)
	if err != nil {
		panic(err)
	}

	// Register DVS module services
	sqModule := sq.NewAppModule(app.logger)
	sqModule.RegisterServices(app.GetMsgRouter())
	sqModule.RegisterInterfaces(app.interfaceRegistry)

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
