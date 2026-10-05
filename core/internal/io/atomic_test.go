package io

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.js")
	if err := WriteFileAtomic(path, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello" {
		t.Fatalf("got %q", got)
	}
	// Overwrite.
	if err := WriteFileAtomic(path, []byte("world"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, _ = os.ReadFile(path)
	if string(got) != "world" {
		t.Fatalf("after overwrite got %q", got)
	}
	// No temp files left behind.
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d: %v", len(entries), entries)
	}
}