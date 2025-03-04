package dvs

import (
	dsm "github.com/0xPellNetwork/pellapp-sdk/service"
	grpc1 "github.com/cosmos/gogoproto/grpc"

	resulthandlers "github.com/0xPellNetwork/dvs-template/dvs/squared/result"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/server"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

// AppModule implements an application module for the dvs module.
type AppModule struct {
	server server.Server
}

// NewAppModule creates a new AppModule object
func NewAppModule(s server.Server) AppModule {
	return AppModule{
		server: s,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(rootServer grpc1.Server) {
	// register dvs-msg handler server
	types.RegisterSquaredMsgServerServer(rootServer, am.server)

	// register dvs-msg result handler
	if r, ok := rootServer.(*dsm.Processor); ok {
		r.RegisterResultMsgExtractor(
			&types.RequestNumberSquaredOut{}, resulthandlers.NewResultHandler(),
		)
	}
}
