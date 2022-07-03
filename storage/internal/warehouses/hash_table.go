package warehouses

import (
	"github.com/dimaglushkov/dkvs/storage/internal"
	"sync"
)

type HashTable struct {
	data map[string]string
	lock sync.RWMutex
}

func NewHashTable() *HashTable {
	return &HashTable{data: make(map[string]string)}
}

func (m *HashTable) Get(k string) (string, error) {
	m.lock.RLock()

	res, ok := m.data[k]

	m.lock.RUnlock()

	if !ok {
		return "", internal.NewUnknownKeyError(k)
	}
	return res, nil
}

func (m *HashTable) Put(k, v string) error {
	m.lock.Lock()

	m.data[k] = v

	m.lock.Unlock()
	return nil
}

func (m *HashTable) Delete(k string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.data[k]; !ok {
		return internal.NewUnknownKeyError(k)
	}

	delete(m.data, k)
	return nil
}
