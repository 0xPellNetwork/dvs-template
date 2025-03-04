package server

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

func (s Server) RequestNumberSquared(ctx context.Context, request *types.RequestNumberSquaredIn) (*types.RequestNumberSquaredOut, error) {
	numInt := request.Task.Squared.Int64()
	s.logger.Info("ProcessRequestNumberSquared", "Number", fmt.Sprintf("%+v", numInt))

	// Calculate square
	squaredInt := numInt * numInt
	squared := math.NewInt(squaredInt)

	s.logger.Info("Calculated square", "input", numInt, "result", squared)
	return &types.RequestNumberSquaredOut{
		TaskIndex: request.Task.TaskIndex,
		Squared:   squared,
	}, nil
}
