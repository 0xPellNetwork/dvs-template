package types

import (
	"context"

	storetypes "cosmossdk.io/store/types"
)

type QueryManager interface {
	Get(ctx context.Context, storeKey storetypes.StoreKey, key []byte) ([]byte, error)
}

type TxManager interface {
	QueryManager
	Set(ctx context.Context, storeKey storetypes.StoreKey, key, value []byte) (storetypes.CommitID, error)
	Delete(ctx context.Context, storeKey storetypes.StoreKey, key []byte) (storetypes.CommitID, error)
}
