package types

import (
	"context"

	storetypes "cosmossdk.io/store/types"
)

type DataManager interface {
	Get(ctx context.Context, storeKey storetypes.StoreKey, key []byte) ([]byte, error)
	Set(ctx context.Context, storeKey storetypes.StoreKey, key, value []byte) error
	Delete(ctx context.Context, storeKey storetypes.StoreKey, key []byte) error
}

type TxManager interface {
	DataManager
	Commit() (storetypes.CommitID, error)
}

type QueryManager interface {
	DataManager
}
