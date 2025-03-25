package app

import (
	"context"
	"fmt"

	storetypes "cosmossdk.io/store/types"

	apptypes "github.com/0xPellNetwork/dvs-template/types"
)

type appDataManager struct {
	provider apptypes.StoreProvider
}

func (m *appDataManager) Get(ctx context.Context, storeKey storetypes.StoreKey, key []byte) ([]byte, error) {
	store := m.provider.QueryMultiStore().GetKVStore(storeKey)
	if store == nil {
		return nil, fmt.Errorf("store %s not found", storeKey)
	}
	return store.Get(key), nil
}

func (m *appDataManager) Set(ctx context.Context, storeKey storetypes.StoreKey, key, value []byte) error {
	store := m.provider.CommitMultiStore().GetKVStore(storeKey)
	if store == nil {
		return fmt.Errorf("store %s not found", storeKey)
	}
	store.Set(key, value)
	return nil
}

func (m *appDataManager) Delete(ctx context.Context, storeKey storetypes.StoreKey, key []byte) error {
	store := m.provider.CommitMultiStore().GetKVStore(storeKey)
	if store == nil {
		return fmt.Errorf("store %s not found", storeKey)
	}
	store.Delete(key)
	return nil
}

type AppTxManager struct {
	*appDataManager
}

func NewAppTxManager(provider apptypes.StoreProvider) *AppTxManager {
	return &AppTxManager{
		appDataManager: &appDataManager{
			provider: provider,
		},
	}
}

func NewAppQueryManager(provider apptypes.StoreProvider) *AppQueryManager {
	return &AppQueryManager{
		appDataManager: &appDataManager{
			provider: provider,
		},
	}
}

func (m *AppTxManager) Commit() (storetypes.CommitID, error) {
	return m.provider.CommitMultiStore().Commit(), nil
}

type AppQueryManager struct {
	*appDataManager
}
