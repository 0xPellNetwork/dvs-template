package types

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
)

func GenItemKey(taskId uint32) string {
	return fmt.Sprintf("task-%d", taskId)
}

// AppQuerier defines the interface for querying app data
type AppCommitStorer interface {
	GetCommitStore(key storetypes.StoreKey) storetypes.KVStore
	GetCommitMultiStore() storetypes.CommitMultiStore
}

type AppQueryStorer interface {
	GetQueryStore(key storetypes.StoreKey) storetypes.KVStore
	GetQueryMultiStore() storetypes.MultiStore
}
