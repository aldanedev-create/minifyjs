package io

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadAllFromString(t *testing.T) {
	b, err := ReadAll(strings.NewReader("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hello" {
		t.Fatalf("got %q", b)
	}
}

func TestExistsIsFileIsDir(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !Exists(file) || !Exists(dir) {
		t.Fatal("Exists returned false for existing paths")
	}
	if !IsFile(file) {
		t.Fatal("IsFile(file) = false")
	}
	if IsDir(file) {
		t.Fatal("IsDir(file) = true")
	}
	if !IsDir(dir) {
		t.Fatal("IsDir(dir) = false")
	}
	if IsFile(dir) {
		t.Fatal("IsFile(dir) = true")
	}
	if Exists(filepath.Join(dir, "nope")) {
		t.Fatal("Exists(nope) = true")
	}
}