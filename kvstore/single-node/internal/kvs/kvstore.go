package kvstore

import "sync"

type value struct {
	Val int `json:"value"`
}

type kvStore struct {
	mu    sync.RWMutex
	store map[int]*value
}

func NewKVStore() *kvStore {
	return &kvStore{
		store: make(map[int]*value),
	}
}

func (kvs *kvStore) Set(key, val int) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	kvs.store[key] = &value{Val: val}
}

func (kvs *kvStore) Get(key int) *value {
	kvs.mu.RLock()
	defer kvs.mu.RUnlock()

	if value, exists := kvs.store[key]; exists {
		return value
	}

	return nil
}
