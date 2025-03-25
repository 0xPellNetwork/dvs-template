package types

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
)

func GenItemKey(taskId uint32) string {
	return fmt.Sprintf("%d", taskId)
}

type StoreProvider interface {
	CommitMultiStore() storetypes.CommitMultiStore
	QueryMultiStore() storetypes.MultiStore
}
