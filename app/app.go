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
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	dvsappcfg "github.com/0xPellNetwork/dvs-template/config"
	sq "github.com/0xPellNetwork/dvs-template/dvs/squared"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
	"google.golang.org/grpc/reflection"
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
	sqModule          *sq.AppModule

	DvsNode *pelldvs.Node
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
	err = app.CommitMultiStore().LoadLatestVersion()
	if err != nil {
		app.logger.Error("Failed to load latest version", "error", err)
		panic(fmt.Sprintf("failed to load latest version: %v", err))
	}

	// set up default query store
	err = app.SetupDefaultQueryStore()
	if err != nil {
		panic(err)
	}
	txMgr := NewAppTxManager(app.BaseApp)
	queryMgr := NewAppQueryManager(app.BaseApp)

	// Register DVS module services
	app.sqModule = sq.NewAppModule(app.logger, storeKey, txMgr, queryMgr)
	app.sqModule.RegisterServices(app.GetMsgRouter())
	app.sqModule.RegisterInterfaces(app.interfaceRegistry)

	app.BaseApp.SetGRPCQueryRouter(baseapp.NewGRPCQueryRouter())
	app.BaseApp.GRPCQueryRouter().SetInterfaceRegistry(app.interfaceRegistry)
	app.sqModule.RegisterQueryServer(app.BaseApp.GRPCQueryRouter())

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

func (app *App) createQueryContextInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//// get metadata from the context
		//// if metadata is not found, return an error
		//// if metadata is found, create a query context and inject it into the context.
		//md, ok := metadata.FromIncomingContext(ctx)
		//if !ok {
		//	return nil, status.Error(codes.Internal, "failed to read metadata")
		//}
		//
		//// create a query context
		////queryCtx, err := app.CreateQueryContext(md.)
		//
		//// inject the query context into the context
		//ctx = context.WithValue(ctx, "app_query_context", queryCtx)

		app.logger.Info("GRPC CreateQueryContextInterceptor")

		return handler(ctx, req)
	}
}

// setup gGRPC server and start listening
func (app *App) SetupQueryGRPCServer() error {
	// create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				grpcrecovery.UnaryServerInterceptor(),
				app.createQueryContextInterceptor(),
			),
		),
	)

	reflection.Register(grpcServer)

	app.BaseApp.RegisterGRPCServer(grpcServer)

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

	err := types.RegisterQueryHandlerFromEndpoint(
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
