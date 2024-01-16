package store

import (
	"sync"
)

type KafkaStore struct {
	mu           sync.RWMutex
	Data         map[string][]int `json:"data"` // Data is stored in the form of  `"k1":[9,5,15]`
	CommitOffset map[string]int   `json:"commit_offset"`
}

func NewKafkaStore() *KafkaStore {
	return &KafkaStore{
		Data:         make(map[string][]int),
		CommitOffset: make(map[string]int),
	}
}

func (ks *KafkaStore) Set(key string, value any) int {

	ks.mu.Lock()
	defer ks.mu.Unlock()

	//Store the data in the Kafka store
	data := value.(int)
	ks.Data[key] = append(ks.Data[key], data)

	// Line based offset starting from offset=0
	return len(ks.Data[key]) - 1
}

func (ks *KafkaStore) Get(req map[string]any) map[string][][2]int {

	ks.mu.RLock()
	defer ks.mu.RUnlock()

	resMap := make(map[string][][2]int)
	for key, data := range req {
		startOffset := int(data.(float64))
		dataStore, exists := ks.Data[key]

		if !exists {
			continue
		}

		for i := startOffset; i < len(dataStore); i++ {
			resMap[key] = append(resMap[key], [2]int{i, dataStore[i]})
		}

	}

	return resMap
}

func (ks *KafkaStore) CommitOffsets(req map[string]any) {

	ks.mu.Lock()
	defer ks.mu.Unlock()

	for key, val := range req {
		ks.CommitOffset[key] = int(val.(float64))
	}
}

func (ks *KafkaStore) GetCommitedOffsets(keys []any) map[string]int {

	ks.mu.RLock()
	defer ks.mu.RUnlock()

	res := make(map[string]int)
	for _, keyI := range keys {
		key := keyI.(string)
		res[key] = ks.CommitOffset[key]
	}

	return res
}
