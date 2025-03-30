package app

import (
	"context"
	"fmt"
	"sync"

	storetypes "cosmossdk.io/store/types"

	apptypes "github.com/0xPellNetwork/dvs-template/types"
)

type appDataManager struct {
	provider apptypes.StoreProvider
	mtx      sync.RWMutex
}

// Get retrieves a value from the store using the provided store key and key.
func (m *appDataManager) Get(ctx context.Context, storeKey storetypes.StoreKey, key []byte) ([]byte, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	store := m.provider.QueryMultiStore().GetKVStore(storeKey)
	if store == nil {
		return nil, fmt.Errorf("store %s not found", storeKey)
	}
	return store.Get(key), nil
}

// Set stores a value in the store using the provided store key and key, returning the commit ID and error.
func (m *appDataManager) Set(ctx context.Context,
	storeKey storetypes.StoreKey,
	key, value []byte,
) (storetypes.CommitID, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	store := m.provider.CommitMultiStore().GetKVStore(storeKey)
	if store == nil {
		return storetypes.CommitID{}, fmt.Errorf("store %s not found", storeKey)
	}
	store.Set(key, value)
	return m.provider.CommitMultiStore().Commit(), nil
}

// Delete removes a value from the store using the provided store key and key, returning the commit ID and error.
func (m *appDataManager) Delete(ctx context.Context, storeKey storetypes.StoreKey, key []byte) (storetypes.CommitID, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	store := m.provider.CommitMultiStore().GetKVStore(storeKey)
	if store == nil {
		return storetypes.CommitID{}, fmt.Errorf("store %s not found", storeKey)
	}
	store.Delete(key)

	return m.provider.CommitMultiStore().Commit(), nil
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

type AppQueryManager struct {
	*appDataManager
}
