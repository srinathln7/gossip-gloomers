package kvstore

import (
	"context"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type value struct {
	Val     int `json:"value"`
	Version int `json:"version"`
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

func (kvs *kvStore) GetSnapShot() map[int]*value {
	return kvs.store
}

func (kvs *kvStore) Set(key, val int) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	kvs.set(key, val)
}

func (kvs *kvStore) SyncLocal(key, val, version int) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	var localVersion int
	if value, exists := kvs.store[key]; exists {
		localVersion = value.Version
	} else {
		localVersion = 0
	}

	if localVersion < version {
		kvs.set(key, val)
	}
}

func (kvs *kvStore) set(key, val int) int {
	var version int
	if _, exists := kvs.store[key]; !exists {
		version = 1
	} else {
		version = kvs.store[key].Version
		version++
	}

	kvs.store[key] = &value{Val: val, Version: version}
	return version
}

func (kvs *kvStore) Get(key int) *value {
	kvs.mu.RLock()
	defer kvs.mu.RUnlock()

	if value, exists := kvs.store[key]; exists {
		return value
	}

	return nil
}

func Broadcast(node *maelstrom.Node, key, value, version int) {
	var body map[string]any = make(map[string]any)
	body["type"] = "gossip"
	body["key"] = key
	body["value"] = value
	body["version"] = version

	nodeList := node.NodeIDs()
	for _, peer := range nodeList {
		if peer == node.ID() {
			continue
		}

		go func(peer string) {
			for {
				_, err := node.SyncRPC(context.Background(), peer, body)
				if err == nil {
					break
				}
			}

		}(peer)
	}
}
