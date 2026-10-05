package cache

import (
	"sync"
	"time"
)

// Memory is an in-process Cache. It is used by tests and by the CLI
// when the user asks for caching within a single run (watch mode
// that does not want a persistent cache).
//
// Memory is safe for concurrent use.
type Memory struct {
	mu      sync.RWMutex
	entries map[string]Entry
	maxSize int
}

// NewMemory returns a memory cache bounded to maxEntries. A
// maxEntries of 0 means unbounded. When the limit is exceeded, the
// oldest entry (by CreatedAt) is evicted on the next Put.
func NewMemory(maxEntries int) *Memory {
	return &Memory{
		entries: make(map[string]Entry),
		maxSize: maxEntries,
	}
}

// Get implements Cache.
func (m *Memory) Get(key string) (Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[key]
	return e, ok
}

// Put implements Cache.
func (m *Memory) Put(key string, e Entry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	m.entries[key] = e
	if m.maxSize > 0 && len(m.entries) > m.maxSize {
		m.evictOldestLocked()
	}
	return nil
}

// Delete implements Cache.
func (m *Memory) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, key)
	return nil
}

// Close implements Cache.
func (m *Memory) Close() error { return nil }

// Len returns the number of entries. Used by tests.
func (m *Memory) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.entries)
}

func (m *Memory) evictOldestLocked() {
	var oldestKey string
	var oldestTime time.Time
	first := true
	for k, e := range m.entries {
		if first || e.CreatedAt.Before(oldestTime) {
			oldestKey = k
			oldestTime = e.CreatedAt
			first = false
		}
	}
	if oldestKey != "" {
		delete(m.entries, oldestKey)
	}
}