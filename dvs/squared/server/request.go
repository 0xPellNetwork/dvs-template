package server

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"

	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
	apptypes "github.com/0xPellNetwork/dvs-template/types"
)

func (s *Server) RequestNumberSquared(ctx context.Context, request *types.RequestNumberSquaredIn) (*types.RequestNumberSquaredOut, error) {
	num, ok := math.NewIntFromString(request.Task.Squared)
	if !ok {
		return nil, fmt.Errorf("failed to convert string to int for %v", request.Task.Squared)
	}

	numInt := num.Int64()
	s.logger.Info("ProcessRequestNumberSquared", "Number", fmt.Sprintf("%+v", numInt))

	// Calculate square
	squaredInt := numInt * numInt
	squared := math.NewInt(squaredInt)

	pkgCtx := sdktypes.UnwrapContext(ctx)
	store := pkgCtx.KVStore(s.storeKey)
	if store == nil {
		return nil, fmt.Errorf("store is not set")
	}

	key := []byte(apptypes.GenItemKey(request.Task.TaskIndex))
	result := types.TaskResult{
		TaskIndex:   request.Task.TaskIndex,
		TaskRequest: request.Task,
		PutOnChain:  false,
	}
	bresult, _ := result.Marshal()
	store.Set(key, bresult)
	commitID, err := s.txMgr.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

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
