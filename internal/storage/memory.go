package storage

import (
	"bytes"
	"context"
	"io"
	"sync"
)

// MemoryStore is an in-memory ObjectStore for tests and local development.
// Put reads the full reader into memory; size and contentType are ignored.
type MemoryStore struct {
	mu      sync.RWMutex
	objects map[string][]byte
}

// NewMemoryStore returns an empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		objects: make(map[string][]byte),
	}
}

// Put implements ObjectStore.
func (m *MemoryStore) Put(ctx context.Context, key string, r io.Reader, _ int64, _ string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if key == "" {
		return ErrInvalidKey
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.objects == nil {
		m.objects = make(map[string][]byte)
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	m.objects[key] = cp
	return nil
}

// Get implements ObjectStore.
func (m *MemoryStore) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if key == "" {
		return nil, ErrInvalidKey
	}

	m.mu.RLock()
	data, ok := m.objects[key]
	m.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

// Delete implements ObjectStore.
func (m *MemoryStore) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if key == "" {
		return ErrInvalidKey
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.objects[key]; !ok {
		return ErrNotFound
	}
	delete(m.objects, key)
	return nil
}

var _ ObjectStore = (*MemoryStore)(nil)
