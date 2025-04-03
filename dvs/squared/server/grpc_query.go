package server

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	storetypes "cosmossdk.io/store/types"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	"github.com/golang/protobuf/proto" //nolint:staticcheck

	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
	apptypes "github.com/0xPellNetwork/dvs-template/types"
)

// make sure Server implements the QueryServiceServer interface
var _ types.QueryServer = &Querier{}

type Querier struct {
	types.UnimplementedQueryServer
	logger   log.Logger
	storeKey storetypes.StoreKey
	queryMgr sdktypes.QueryManager
}

// NewServer creates a new Server instance
func NewQuerier(logger log.Logger, storeKey storetypes.StoreKey, queryMgr sdktypes.QueryManager) (*Querier, error) {
	return &Querier{
		logger:   logger,
		storeKey: storeKey,
		queryMgr: queryMgr,
	}, nil
}

// Task retrieves data for a given task index, like 0
func (s *Querier) Task(ctx context.Context, req *types.QueryTaskRequest) (*types.QueryTaskResponse, error) {
	s.logger.Info("Task One  request",
		"key", req.TaskIndex,
		"store_key", s.storeKey.String(),
	)
	key := []byte(req.TaskIndex)
	value, _ := s.queryMgr.Get(ctx, s.storeKey, key)
	if len(value) == 0 {
		s.logger.Error("failed to get value from store", "key", req.TaskIndex)
		return nil, fmt.Errorf("failed to get value for key: %s", req.TaskIndex)
	}

	var result = types.TaskResult{}
	err := proto.Unmarshal(value, &result)
	if err != nil {
		s.logger.Error("failed to unmarshal task result", "error", err)
	}

	s.logger.Info("Task One request",
		"store-key-str", req.TaskIndex,
		"store-key-bytes", fmt.Sprintf("%+v", key),
		"store-value-raw", result,
		"store-value-bytes", value,
	)
	return &types.QueryTaskResponse{
		Value: &result,
	}, nil
}

// Tasks lists all data with a task indexs like "0,1,2,3,4"
func (s *Querier) Tasks(ctx context.Context, req *types.QueryTasksRequest) (*types.QueryTasksResponse, error) {
	s.logger.Info("Tasks List request", "keys", req.TaskIndexes)

	keyStrList := strings.Split(req.TaskIndexes, ",")
	if len(keyStrList) == 0 {
		return &types.QueryTasksResponse{}, nil
	}

	// convert keyStrList to []byte
	keyList := make([][]byte, len(keyStrList))
	for _, keyStr := range keyStrList {
		trimmedKey := strings.TrimSpace(keyStr)
		if len(trimmedKey) == 0 {
			continue
		}

		taskID, ok := big.NewInt(0).SetString(trimmedKey, 10)
		if !ok {
			s.logger.Error("invalid task ID", "key", trimmedKey)
			continue
		}

		keyByte := apptypes.GenItemKey(uint32(taskID.Uint64()))
		keyList = append(keyList, []byte(keyByte))
	}

	if len(keyList) == 0 {
		return &types.QueryTasksResponse{}, nil
	}

	var result = types.QueryTasksResponse{}
	for _, key := range keyList {
		if len(key) == 0 {
			continue
		}
		s.logger.Debug("Tasks List get item", "key-str", string(key))

		value, err := s.queryMgr.Get(ctx, s.storeKey, key)
		if err != nil {
			s.logger.Error("failed to get value from store", "key", string(key), "error", err)
			continue
		}
		if len(value) == 0 {
			continue
		}

		var taskResult types.TaskResult
		err = proto.Unmarshal(value, &taskResult)
		if err != nil {
			s.logger.Error("failed to unmarshal task result", "error", err)
			continue
		}

		result.Items = append(result.Items, &taskResult)
	}
	return &result, nil
}
