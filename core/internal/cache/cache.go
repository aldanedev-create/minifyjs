// Package cache stores MinifyJS transform results on disk so repeated
// builds of unchanged input skip the esbuild call entirely.
//
// This matters for two concrete scenarios:
//
//   - Watch mode: a source file is edited, saved, saved again with no
//     change, and saved a third time with a real change. Only the
//     third save should do work.
//
//   - CI: a project has 500 JS files and a build step that runs
//     `minifyjs` on each. If only one file changed since the last
//     green build, only that file needs re-minifying.
//
// The cache is opt-in (config.Cache is false by default) because it
// writes to disk. When enabled, it is a content-addressed store: the
// key is a hash of (input bytes, options), so a cache hit guarantees
// the cached output is exactly what a fresh run would produce.
package cache

import (
	"time"
)

// Entry is one cached transform result.
type Entry struct {
	// Code is the minified output.
	Code string

	// Map is the source map, if one was produced. Empty otherwise.
	Map string

	// CreatedAt is when the entry was written. Used only for cache
	// eviction; the CLI does not read it.
	CreatedAt time.Time
}

// Cache is the interface satisfied by the memory and filesystem
// implementations. The CLI holds one Cache; tests can substitute a
// memory cache.
type Cache interface {
	// Get looks up an entry by key. Returns the entry and true on
	// hit, or the zero Entry and false on miss.
	Get(key string) (Entry, bool)

	// Put stores an entry under key. It overwrites any existing
	// entry with the same key.
	Put(key string, e Entry) error

	// Delete removes the entry under key. A no-op if the key is not
	// present.
	Delete(key string) error

	// Close releases any resources held by the cache. Calling Close
	// on a memory cache is a no-op.
	Close() error
}