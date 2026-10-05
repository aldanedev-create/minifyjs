package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

// TmpDir creates a temp directory that is NOT auto-cleaned on test
// completion. It is used by tests that want to inspect the directory
// after a failure, or that build artifacts the test suite will reuse.
//
// When the test passes, the directory is removed. When it fails, the
// directory path is logged and left on disk.
func TmpDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "minifyjs-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if t.Failed() {
			t.Logf("test failed; leaving temp dir at %s", dir)
			return
		}
		_ = os.RemoveAll(dir)
	})
	return dir
}

// WriteFile writes a file under dir with the given name and content,
// creating intermediate directories as needed. It returns the full
// path.
func WriteFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// WriteString is WriteFile for string content.
func WriteString(t *testing.T, dir, name, content string) string {
	t.Helper()
	return WriteFile(t, dir, name, []byte(content))
}