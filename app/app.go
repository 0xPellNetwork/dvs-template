package app

import (
	"fmt"
	"io"

	storetypes "cosmossdk.io/store/types"
	"github.com/0xPellNetwork/pellapp-sdk/baseapp"
	"github.com/0xPellNetwork/pellapp-sdk/pelldvs"
	"github.com/0xPellNetwork/pellapp-sdk/server/api"
	servercfg "github.com/0xPellNetwork/pellapp-sdk/server/config"
	servertypes "github.com/0xPellNetwork/pellapp-sdk/server/types"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	pelldvscfg "github.com/0xPellNetwork/pelldvs/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
	InterfaceRegistry codectypes.InterfaceRegistry
	ModuleManager     *sdktypes.ModuleManager

	DvsNode *pelldvs.Node
}

type DBContext struct {
	ID     string
	Config *pelldvscfg.Config
}

// DBProvider takes a DBContext and returns an instantiated DB.
type DBProvider func(*DBContext) (dbm.DB, error)

// DefaultDBProvider returns a database using the DBBackend and DBDir
// specified in the ctx.Config.
func DefaultDBProvider(ctx *DBContext) (dbm.DB, error) {
	dbType := dbm.BackendType(ctx.Config.DBBackend)
	return dbm.NewDB(ctx.ID, dbType, ctx.Config.DBDir())
}

// NewApp initializes a new App instance
func NewApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
	baseappOptions ...func(*baseapp.BaseApp),
) *App {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	app := &App{
		BaseApp:           baseapp.NewBaseApp(Name, logger, db, cdc),
		InterfaceRegistry: interfaceRegistry,
		logger:            logger,
		appCodec:          cdc,
	}

	// Register standard types interfaces
	std.RegisterInterfaces(app.InterfaceRegistry)
	authtypes.RegisterInterfaces(app.InterfaceRegistry)

	exampleStoreKey := storetypes.NewKVStoreKey("dvs-template")
	app.logger.Info("Mounting query store", "key", exampleStoreKey.String())

	// mount the store
	app.MountStore(exampleStoreKey, storetypes.StoreTypeIAVL)

	// load latest version
	app.logger.Info("Loading latest version")
	err := app.CommitMultiStore().LoadLatestVersion()
	if err != nil {
		app.logger.Error("Failed to load latest version", "error", err)
		panic(fmt.Sprintf("failed to load latest version: %v", err))
	}

	app.BaseApp.SetGRPCQueryRouter(baseapp.NewGRPCQueryRouter())
	app.BaseApp.GRPCQueryRouter().SetInterfaceRegistry(app.InterfaceRegistry)

	app.ModuleManager = sdktypes.NewManager()
	app.ModuleManager.RegisterInterfaces(app.InterfaceRegistry)
	app.ModuleManager.RegisterServices(app.GetMsgRouter().GetConfigurator())
	app.ModuleManager.RegisterResultMsgExtractors(app.GetMsgRouter().GetConfigurator())
	app.ModuleManager.RegisterQueryServices(app.BaseApp.GRPCQueryRouter())

	app.logger.Info("interface registry",
		"allInterfaces", app.InterfaceRegistry.ListAllInterfaces(),
		"InterfaceRegistry", app.InterfaceRegistry,
	)

	return app
}

func (app *App) RegisterAPIRoutes(apiSvr *api.Server, _ servercfg.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register grpc-gateway routes for all modules.
	app.ModuleManager.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
}
