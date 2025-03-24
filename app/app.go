package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	storetypes "cosmossdk.io/store/types"
	"github.com/0xPellNetwork/pellapp-sdk/baseapp"
	"github.com/0xPellNetwork/pellapp-sdk/pelldvs"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	"github.com/0xPellNetwork/pelldvs/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdktypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	dvsappcfg "github.com/0xPellNetwork/dvs-template/config"
	"github.com/0xPellNetwork/dvs-template/dvs/query"
	"github.com/0xPellNetwork/dvs-template/dvs/query/types"
	sq "github.com/0xPellNetwork/dvs-template/dvs/squared"
	"google.golang.org/grpc/credentials/insecure"
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
	dvsAppConfig *dvsappcfg.AppConfig
	*baseapp.BaseApp

	appCodec          codec.Codec
	logger            log.Logger
	interfaceRegistry codectypes.InterfaceRegistry

	DvsNode     *pelldvs.Node
	queryModule *query.AppModule
}

type DBContext struct {
	ID     string
	Config *config.Config
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
	interfaceRegistry codectypes.InterfaceRegistry,
	logger log.Logger,
	cfg *config.Config,
	appConfig *dvsappcfg.AppConfig,
) *App {
	cdc := codec.NewProtoCodec(interfaceRegistry)
	db, err := DefaultDBProvider(&DBContext{
		ID:     Name + "-db",
		Config: cfg,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to initialize db: %v", err))
	}

	app := &App{
		BaseApp:           baseapp.NewBaseApp(Name, logger, db, cdc),
		interfaceRegistry: interfaceRegistry,
		logger:            logger,
		appCodec:          cdc,
		dvsAppConfig:      appConfig,
	}

	// Register standard types interfaces
	std.RegisterInterfaces(app.interfaceRegistry)
	sdktypes.RegisterInterfaces(app.interfaceRegistry)

	// Build dvs node
	app.DvsNode, err = pelldvs.NewNode(app.logger, app, cfg)
	if err != nil {
		panic(err)
	}

	storeKey := storetypes.NewKVStoreKey(appConfig.QueryStoreKey)
	app.logger.Info("Mounting query store", "key", storeKey.String())
	// mount the store
	app.MountStore(storeKey, storetypes.StoreTypeIAVL)

	// load latest version
	app.logger.Info("Loading latest version")
	err = app.GetCommitMultiStore().LoadLatestVersion()
	if err != nil {
		app.logger.Error("Failed to load latest version", "error", err)
		panic(fmt.Sprintf("failed to load latest version: %v", err))
	}

	// set up default query store
	err = app.SetupDefaultQueryStore()
	if err != nil {
		panic(err)
	}

	// Register DVS module services
	sqModule := sq.NewAppModule(app.logger, storeKey)
	sqModule.RegisterServices(app.GetMsgRouter())
	sqModule.RegisterInterfaces(app.interfaceRegistry)
	sqModule.SetAppCommitStore(app)

	// Register query services
	app.queryModule = query.NewAppModule(app.logger, storeKey)
	app.queryModule.SetAppQuerier(app)

	return app
}

// Start method starts the application
func (app *App) Start() error {
	app.logger.Info("App Start")
	if err := app.DvsNode.Start(); err != nil {
		app.logger.Error("DvsNode Start Failed", "error", err.Error())
		return err
	}

	// start query http server
	if err := app.SetupQueryHTTPServer(
		app.dvsAppConfig.QueryRPCServerAddress,
		app.dvsAppConfig.QueryHTTPServerAddress); err != nil {
		app.logger.Error("Failed to setup HTTP server", "error", err)
		return err
	}

	// start query grpc server
	if err := app.SetupQueryGRPCServer(); err != nil {
		app.logger.Error("Failed to setup gRPC server", "error", err)
		return err
	}

	// Block the main thread
	c := make(chan any)
	<-c
	return nil
}

// setup gGRPC server and start listening
func (app *App) SetupQueryGRPCServer() error {
	// create gRPC server
	grpcServer := grpc.NewServer()

	app.queryModule.RegisterGRPCServices(grpcServer)

	go func() {
		lis, err := net.Listen("tcp", app.dvsAppConfig.QueryRPCServerAddress)
		if err != nil {
			app.logger.Error("Failed to listen", "error", err)
			return
		}
		app.logger.Info("Starting Query gRPC server", "address", app.dvsAppConfig.QueryRPCServerAddress)
		if err := grpcServer.Serve(lis); err != nil {
			app.logger.Error("Failed to serve", "error", err)
		}
	}()

	return nil
}

func (app *App) SetupQueryHTTPServer(grpcAddr, httpAddr string) error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := types.RegisterQueryServiceHandlerFromEndpoint(
		ctx, mux, grpcAddr, opts,
	)
	if err != nil {
		return err
	}

	go func() {
		app.logger.Info("Starting Query HTTP server", "address", httpAddr)
		err := http.ListenAndServe(httpAddr, mux)
		if err != nil {
			app.logger.Error("Failed to start HTTP server", "error", err)
		}
	}()
	return nil
}
