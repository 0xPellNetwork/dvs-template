package server

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

type RequestServer struct {
	Server
}

// NewDvsProcessRequestServer returns an implementation of the DvsProcessRequestServer interface
// for the provided Server.
func NewRequestServer(server Server) types.DVSRequestServer {
	return &RequestServer{
		Server: server,
	}
}

var _ types.DVSRequestServer = RequestServer{}

func (server RequestServer) RequestNumberSquared(ctx context.Context, request *types.RequestNumberSquaredIn) (*types.RequestNumberSquaredOut, error) {
	numInt := request.Task.Squared.Int64()
	server.logger.Info("ProcessRequestNumberSquared", "Number", fmt.Sprintf("%+v", numInt))

	// Calculate square
	squaredInt := numInt * numInt
	squared := math.NewInt(squaredInt)

	server.logger.Info("Calculated square", "input", numInt, "result", squared)
	return &types.RequestNumberSquaredOut{
		TaskIndex: request.Task.TaskIndex,
		Squared:   squared,
	}, nil
}
