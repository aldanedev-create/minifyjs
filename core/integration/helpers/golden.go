package helpers

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// AssertGolden compares got against the file at
// testdata/golden/<name>.golden. If the UPDATE_GOLDEN environment
// variable is set to a non-empty value, the golden file is written
// (or overwritten) with got instead.
//
// Golden files are useful for outputs that are correct but tedious to
// spell out inline: minified JavaScript, source maps, JSON
// diagnostics. Tests stay readable because the expected value is a
// file next to the test, not a 300-character string literal.
func AssertGolden(t *testing.T, name string, got []byte) {
	t.Helper()

	_, thisFile, _, _ := runtimeCaller()
	dir := filepath.Join(filepath.Dir(thisFile), "testdata", "golden")
	path := filepath.Join(dir, name+".golden")

	if os.Getenv("UPDATE_GOLDEN") != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("golden file missing: %v (run with UPDATE_GOLDEN=1 to create)", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("golden mismatch for %s:\n--- want ---\n%s\n--- got ---\n%s",
			name, want, got)
	}
}