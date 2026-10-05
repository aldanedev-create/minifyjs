package helpers

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestdataDir returns the absolute path to core/integration/testdata.
func TestdataDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "..", "testdata")
}

// InputFixture returns the bytes of a file under testdata/inputs.
func InputFixture(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join(TestdataDir(t), "inputs", name)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return b
}

// InputFixturePath returns the absolute path of a file under
// testdata/inputs. Use this when a test needs to pass the path to the
// CLI rather than the content.
func InputFixturePath(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join(TestdataDir(t), "inputs", name)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("fixture %s: %v", name, err)
	}
	return path
}

// ConfigFixturePath returns the absolute path of a file under
// testdata/configs.
func ConfigFixturePath(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join(TestdataDir(t), "configs", name)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("config fixture %s: %v", name, err)
	}
	return path
}

// CopyFixtureDir copies a bundle project under
// testdata/bundle-projects into a fresh temp directory and returns
// the temp directory path. Bundle tests need to write output next to
// the inputs, and testdata is read-only.
func CopyFixtureDir(t *testing.T, name string) string {
	t.Helper()
	src := filepath.Join(TestdataDir(t), "bundle-projects", name)
	if _, err := os.Stat(src); err != nil {
		t.Fatalf("bundle fixture %s: %v", name, err)
	}
	dst := TmpDir(t)
	if err := copyTree(src, dst); err != nil {
		t.Fatalf("copy fixture: %v", err)
	}
	return dst
}

func copyTree(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}