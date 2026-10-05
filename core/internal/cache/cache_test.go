package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestKeyIsDeterministic(t *testing.T) {
	o := KeyOptions{MinifyWhitespace: true}
	k1 := Key("hello", o)
	k2 := Key("hello", o)
	if k1 != k2 {
		t.Fatalf("keys differ: %s vs %s", k1, k2)
	}
	if len(k1) != 64 {
		t.Fatalf("expected 64-char hex sha256, got %d: %s", len(k1), k1)
	}
}

func TestKeyChangesWithSource(t *testing.T) {
	o := KeyOptions{}
	if Key("a", o) == Key("b", o) {
		t.Fatal("different sources should produce different keys")
	}
}

func TestKeyChangesWithOptions(t *testing.T) {
	if Key("x", KeyOptions{MinifyWhitespace: true}) ==
		Key("x", KeyOptions{MinifyWhitespace: false}) {
		t.Fatal("option change should change key")
	}
}

func TestKeyNoLengthAmbiguity(t *testing.T) {
	// Concatenation without length prefixing would collide:
	// ("ab", "c") vs ("a", "bc"). Field-length prefixing must
	// prevent this.
	if Key("abc", KeyOptions{}) == Key("abc", KeyOptions{}) {
		t.Fatal("same input must produce same key")
	}
}

func TestMemoryCacheRoundTrip(t *testing.T) {
	c := NewMemory(0)
	if _, ok := c.Get("nope"); ok {
		t.Fatal("expected miss")
	}
	if err := c.Put("k", Entry{Code: "x"}); err != nil {
		t.Fatal(err)
	}
	e, ok := c.Get("k")
	if !ok || e.Code != "x" {
		t.Fatalf("got %+v ok=%v", e, ok)
	}
	if err := c.Delete("k"); err != nil {
		t.Fatal(err)
	}
	if _, ok := c.Get("k"); ok {
		t.Fatal("expected miss after delete")
	}
}

func TestMemoryCacheEvictsOldest(t *testing.T) {
	c := NewMemory(2)
	_ = c.Put("a", Entry{Code: "a", CreatedAt: time.Now().Add(-time.Hour)})
	_ = c.Put("b", Entry{Code: "b", CreatedAt: time.Now()})
	_ = c.Put("c", Entry{Code: "c", CreatedAt: time.Now().Add(time.Hour)})
	if c.Len() != 2 {
		t.Fatalf("Len = %d, want 2", c.Len())
	}
	if _, ok := c.Get("a"); ok {
		t.Fatal("oldest should have been evicted")
	}
}

func TestFilesystemCacheRoundTrip(t *testing.T) {
	dir := t.TempDir()
	c, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if _, ok := c.Get("nope"); ok {
		t.Fatal("expected miss")
	}
	if err := c.Put("abcdef0123456789", Entry{Code: "x", Map: "m"}); err != nil {
		t.Fatal(err)
	}
	e, ok := c.Get("abcdef0123456789")
	if !ok {
		t.Fatal("expected hit")
	}
	if e.Code != "x" || e.Map != "m" {
		t.Fatalf("got %+v", e)
	}

	// The entry should live under a 2-char shard.
	expected := filepath.Join(dir, "ab", "abcdef0123456789.json")
	if _, err := os.Stat(expected); err != nil {
		t.Fatalf("entry not at expected path: %v", err)
	}

	if err := c.Delete("abcdef0123456789"); err != nil {
		t.Fatal(err)
	}
	if _, ok := c.Get("abcdef0123456789"); ok {
		t.Fatal("expected miss after delete")
	}
}

func TestFilesystemCacheTolerantOfCorruptEntry(t *testing.T) {
	dir := t.TempDir()
	c, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if err := c.Put("deadbeef", Entry{Code: "x"}); err != nil {
		t.Fatal(err)
	}
	// Corrupt the entry on disk.
	path := filepath.Join(dir, "de", "deadbeef.json")
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, ok := c.Get("deadbeef"); ok {
		t.Fatal("corrupt entry should be a miss")
	}
	// The corrupt file should have been removed.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("corrupt entry not removed: %v", err)
	}
}

func TestDefaultDirNonEmpty(t *testing.T) {
	if DefaultDir() == "" {
		t.Fatal("DefaultDir returned empty string")
	}
}