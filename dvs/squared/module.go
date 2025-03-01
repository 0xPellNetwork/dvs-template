package dvs

import (
	dsm "github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler"
	grpc1 "github.com/cosmos/gogoproto/grpc"

	resulthandlers "github.com/0xPellNetwork/dvs-template/dvs/squared/result"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/server"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

// AppModule implements an application module for the dvs module.
type AppModule struct {
	server         server.Server
	RequestServer  grpc1.Server
	ResponseServer grpc1.Server
}

// NewAppModule creates a new AppModule object
func NewAppModule(s server.Server) AppModule {
	return AppModule{
		server:         s,
		RequestServer:  dsm.GetRequestHandler(),
		ResponseServer: dsm.GetResponseHandler(),
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices() {
	dvsProcessRequestServer := server.NewRequestServer(am.server)
	dvsPostProcessRequestServer := server.NewResponseServer(am.server)

	// register dvs-msg handler server
	types.RegisterDVSRequestServer(am.RequestServer, dvsProcessRequestServer)
	types.RegisterDVSResponseServer(am.ResponseServer, dvsPostProcessRequestServer)

	// register dvs-msg result handler
	if r, ok := am.RequestServer.(*dsm.RequestHandler); ok {
		r.RegisterResultHandler(
			&types.RequestNumberSquaredOut{}, resulthandlers.NewResultHandler(),
		)
	}
}
