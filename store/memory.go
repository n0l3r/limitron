package store

import (
	"sync"
	"time"
)

type memoryItem struct {
	value  int64
	expiry int64
}

type MemoryStore struct {
	sync.RWMutex
	data map[string]memoryItem
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]memoryItem),
	}
}

func (ms *MemoryStore) Get(key string) (int64, error) {
	ms.RLock()
	defer ms.RUnlock()
	item, exists := ms.data[key]
	if !exists || item.expiry < time.Now().Unix() {
		return 0, nil
	}
	return item.value, nil
}

func (ms *MemoryStore) Set(key string, value int64, expiry int64) error {
	ms.Lock()
	defer ms.Unlock()
	ms.data[key] = memoryItem{value, time.Now().Unix() + expiry}
	return nil
}

func (ms *MemoryStore) Incr(key string, n int64, expiry int64) (int64, error) {
	ms.Lock()
	defer ms.Unlock()
	item, exists := ms.data[key]
	if !exists || item.expiry < time.Now().Unix() {
		ms.data[key] = memoryItem{n, time.Now().Unix() + expiry}
		return n, nil
	}
	item.value += n
	ms.data[key] = item
	return item.value, nil
}
