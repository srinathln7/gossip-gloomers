package cache

import (
	"errors"
	"sync"
)

type KVStore struct {
	store map[int]struct{}
	mu    sync.Mutex
}

// newCache : initializes and returns a new key-value store
func NewCache() *KVStore {
	return &KVStore{
		store: make(map[int]struct{}),
	}
}

// Set: Sets the given value in the key-value store
func (kv *KVStore) Set(val int) error {

	kv.mu.Lock()
	defer kv.mu.Unlock()

	if _, exists := kv.store[val]; exists {
		return errors.New("specified key already exists in the cache")
	}

	kv.store[val] = struct{}{}

	return nil
}

// SetAll: Sets all the values specified in the slice in the key-value store
func (kv *KVStore) SetAll(vals []int) error {

	for _, val := range vals {
		kv.Set(val)
	}

	return nil
}

// Get: Returns all the values specified in the store
func (kv *KVStore) Get() []int {

	kv.mu.Lock()
	defer kv.mu.Unlock()

	var keys []int
	for key := range kv.store {
		keys = append(keys, key)
	}

	return keys
}
