package commands

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/0xPellNetwork/pellapp-sdk/client"
	"github.com/0xPellNetwork/pellapp-sdk/server"
	servertypes "github.com/0xPellNetwork/pellapp-sdk/server/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	pelldvscfg "github.com/0xPellNetwork/pelldvs/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/0xPellNetwork/dvs-template/app"
)

func addFlags(cmd *cobra.Command) {
	cmd.Flags().String("config", "", "Path to the dvstemplated config file")
}

// execPostSetup initializes and configures various components of the server after setup, such as loading configurations,
// creating chain connector clients, and starting the task dispatcher. It runs concurrent tasks using an error group.
func execPostSetup(svrCtx *server.Context, clientCtx client.Context, ctx context.Context, app servertypes.Application, g *errgroup.Group) error {
	var configPath = DefaultNodeHome + "/config/config.toml"
	var pellDVSConfig pelldvscfg.Config
	vp := viper.New()
	vp.SetConfigFile(configPath)
	if err := vp.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("failed to read in config: %v", err))
	}
	err := vp.Unmarshal(&pellDVSConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	pellDVSConfig.SetRoot(DefaultNodeHome)

	svrCtx.Logger.Info("Dispatcher started")
	return nil
}

// postSetup is called after the application is initialized
func postSetup(svrCtx *server.Context, clientCtx client.Context, ctx context.Context, app servertypes.Application, g *errgroup.Group) error {
	// start dispatcher here
	g.Go(func() error {
		return execPostSetup(svrCtx, clientCtx, ctx, app, g)
	})
	return nil
}

// openDB opens the database
func openDB(rootDir string, backendType dbm.BackendType) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return dbm.NewDB("dvstemplatedapp", backendType, dataDir)
}

// InitRunOperatorCommand initializes the run operator command
func InitRunOperatorCommand(rootCmd *cobra.Command) {
	server.AddCommandsWithStartCmdOptions(
		rootCmd,
		DefaultNodeHome,
		newApp,
		server.StartCmdOptions{
			DBOpener:            openDB,
			PostSetup:           postSetup,
			AddFlags:            addFlags,
			StartCommandHandler: nil,
		},
	)
}

// newApp creates the application
func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)
	return app.NewApp(
		logger,
		db,
		traceStore,
		appOpts,
		baseappOptions...,
	)
}
