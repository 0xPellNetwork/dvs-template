package server

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	storetypes "cosmossdk.io/store/types"
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
	queryMgr apptypes.QueryManager
}

// NewServer creates a new Server instance
func NewQuerier(logger log.Logger, storeKey storetypes.StoreKey, queryMgr apptypes.QueryManager) (*Querier, error) {
	return &Querier{
		logger:   logger,
		storeKey: storeKey,
		queryMgr: queryMgr,
	}, nil
}

// GetData retrieves data for a given key
func (s *Querier) GetData(ctx context.Context, req *types.GetDataRequest) (*types.GetDataResponse, error) {
	s.logger.Info("GetData request",
		"key", req.Key,
		"store_key", s.storeKey.String(),
	)

	key := []byte(req.Key)
	value, err := s.queryMgr.Get(ctx, s.storeKey, []byte(req.Key))
	if len(value) == 0 {
		s.logger.Error("failed to get value from store", "key", req.Key)
		return nil, fmt.Errorf("failed to get value for key: %s", req.Key)
	}

	var result = types.TaskResult{}
	err = proto.Unmarshal(value, &result)
	if err != nil {
		s.logger.Error("failed to unmarshal task result", "error", err)
	}

	s.logger.Info("GetData request",
		"store-key-str", req.Key,
		"store-key-bytes", fmt.Sprintf("%+v", key),
		"store-value-raw", result,
		"store-value-bytes", value,
	)
	return &types.GetDataResponse{
		Value: &result,
	}, nil
}

// ListData lists all data with a key list like "task-01,task-02,key3"
func (s *Querier) ListData(ctx context.Context, req *types.ListDataRequest) (*types.ListDataResponse, error) {
	s.logger.Debug("ListData request", "keys", req.Keys)

	keyStrList := strings.Split(req.Keys, ",")
	if len(keyStrList) == 0 {
		return &types.ListDataResponse{}, nil
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
		return &types.ListDataResponse{}, nil
	}

	var result = types.ListDataResponse{}
	for _, key := range keyList {
		if len(key) == 0 {
			continue
		}
		s.logger.Info("ListData get item", "key-str", string(key))

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

		result.Items = append(result.Items, &types.ListItem{
			Key:   string(key),
			Value: &taskResult,
		})
	}
	return &result, nil
}
