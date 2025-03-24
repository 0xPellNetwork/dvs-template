package server

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
	dvscommontypes "github.com/0xPellNetwork/dvs-template/dvs/types"
	apptypes "github.com/0xPellNetwork/dvs-template/types"
)

func (s *Server) RequestNumberSquared(ctx context.Context, request *types.RequestNumberSquaredIn) (*types.RequestNumberSquaredOut, error) {
	numInt := request.Task.Squared.Int64()
	s.logger.Info("ProcessRequestNumberSquared", "Number", fmt.Sprintf("%+v", numInt))

	// Calculate square
	squaredInt := numInt * numInt
	squared := math.NewInt(squaredInt)

	if s.app == nil {
		return nil, fmt.Errorf("app is not set")
	}

	store := s.app.GetCommitStore(s.storeKey)
	if store == nil {
		return nil, fmt.Errorf("store is not set")
	}

	key := []byte(apptypes.GenItemKey(request.Task.TaskIndex))
	result := dvscommontypes.TaskResult{
		TaskIndex:   request.Task.TaskIndex,
		TaskRequest: request.Task,
		PutOnChain:  false,
	}
	bresult, _ := result.Marshal()
	store.Set(key, bresult)
	commitID := s.app.GetCommitMultiStore().Commit()

	s.logger.Info("Calculated square", "input", numInt, "result", squared, "store-commit-id", commitID,
		"store-key-str", string(key),
		"store-key-bytes", fmt.Sprintf("%+v", key),
		"store-value-raw", result,
		"store-value-bytes", fmt.Sprintf("%+v", bresult),
	)

	return &types.RequestNumberSquaredOut{
		TaskIndex: request.Task.TaskIndex,
		Squared:   squared,
	}, nil
}
