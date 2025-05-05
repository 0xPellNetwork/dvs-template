package server

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
	apptypes "github.com/0xPellNetwork/dvs-template/types"
)

func (s *Server) RequestNumberSquared(ctx context.Context, request *types.RequestNumberSquaredIn) (*types.RequestNumberSquaredOut, error) {
	num, ok := math.NewIntFromString(request.Task.Squared)
	if !ok {
		return nil, fmt.Errorf("failed to convert string to int for %v", request.Task.Squared)
	}

	sq := num.Mul(num)

	key := []byte(apptypes.GenItemKey(request.Task.TaskIndex))
	result := types.TaskResult{
		TaskRequest: request.Task,
		IsOnChain:   false,
	}

	bresult, _ := result.Marshal()
	commitID, err := s.txMgr.Set(ctx, s.storeKey, key, bresult)
	if err != nil {
		return nil, fmt.Errorf("failed to set value in store: %w", err)
	}

	s.logger.Info("Calculated square",
		"input", num.String(),
		"result", sq.String(),
		"store-key-str", string(key),
		"store-key-bytes", fmt.Sprintf("%+v", key),
		"store-value-raw", result,
		"store-value-bytes", fmt.Sprintf("%+v", bresult),
		"store-commit-id", commitID,
	)

	return &types.RequestNumberSquaredOut{
		TaskIndex: request.Task.TaskIndex,
		Squared:   sq.String(),
	}, nil
}
