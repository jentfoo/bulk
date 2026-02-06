package store

import (
	"slices"
	"sync"
)

type memStore struct {
	mu   sync.Mutex
	data map[string][]byte
}

// NewMemStore returns an in-memory Storage implementation.
// All operations are safe for concurrent use.
func NewMemStore() Storage {
	return &memStore{data: make(map[string][]byte)}
}

func (m *memStore) Set(key string, blob []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = slices.Clone(blob)
	return nil
}

func (m *memStore) Get(key string) ([]byte, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	blob, ok := m.data[key]
	if !ok {
		return nil, false, nil
	}
	return slices.Clone(blob), true, nil
}

func (m *memStore) ContainsKey(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.data[key]
	return ok
}

func (m *memStore) KeySet() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

func (m *memStore) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
}

func (m *memStore) DeleteAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	clear(m.data)
}

func (m *memStore) Size() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.data)
}

func (m *memStore) Close() error {
	return nil
}
