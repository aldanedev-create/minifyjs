package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Filesystem is a Cache that stores entries as JSON files under a
// single directory.
//
// Layout:
//
//	<cacheDir>/<2-char-shard>/<full-key>.json
//
// The two-character shard spreads entries over 256 subdirectories so
// a project with tens of thousands of cached files does not produce a
// single directory with tens of thousands of entries (which is slow
// on every filesystem).
type Filesystem struct {
	dir string
	mu  sync.Mutex
}

// Open creates (if needed) a Filesystem cache rooted at dir. The dir
// is created with mode 0o755.
func Open(dir string) (*Filesystem, error) {
	if dir == "" {
		return nil, fmt.Errorf("cache: empty directory")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("cache: mkdir %s: %w", dir, err)
	}
	return &Filesystem{dir: dir}, nil
}

// Dir returns the cache root.
func (c *Filesystem) Dir() string { return c.dir }

// Get implements Cache.
func (c *Filesystem) Get(key string) (Entry, bool) {
	path := c.path(key)
	data, err := os.ReadFile(path)
	if err != nil {
		return Entry{}, false
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		// Corrupt entry; treat as a miss and best-effort remove.
		_ = os.Remove(path)
		return Entry{}, false
	}
	return e, true
}

// Put implements Cache. Writes are serialized through a mutex
// because two goroutines may race on the same shard directory.
func (c *Filesystem) Put(key string, e Entry) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	path := c.path(key)

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("cache: mkdir shard: %w", err)
	}
	// Atomic write so a crash mid-write cannot leave a truncated
	// entry that a later Get would mis-decode.
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tmp-*")
	if err != nil {
		return fmt.Errorf("cache: temp: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("cache: write temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("cache: close temp: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("cache: rename: %w", err)
	}
	return nil
}

// Delete implements Cache.
func (c *Filesystem) Delete(key string) error {
	err := os.Remove(c.path(key))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

// Close implements Cache. It is a no-op for Filesystem, which holds
// no open handles between operations.
func (c *Filesystem) Close() error { return nil }

// Purge removes every entry in the cache. Useful for a future
// `minifyjs --clear-cache` flag.
func (c *Filesystem) Purge() error {
	return os.RemoveAll(c.dir)
}

func (c *Filesystem) path(key string) string {
	if len(key) < 2 {
		return filepath.Join(c.dir, key+".json")
	}
	return filepath.Join(c.dir, key[:2], key+".json")
}

// DefaultDir returns the OS-appropriate default cache directory.
//
// Linux/macOS: $XDG_CACHE_HOME/minifyjs, or ~/.cache/minifyjs
// Windows:     %LOCALAPPDATA%\minifyjs\Cache
//
// The function intentionally does not create the directory; Open does
// that.
func DefaultDir() string {
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, "minifyjs")
	}
	if local := os.Getenv("LOCALAPPDATA"); local != "" {
		return filepath.Join(local, "minifyjs", "Cache")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".minifyjs-cache"
	}
	return filepath.Join(home, ".cache", "minifyjs")
}