package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestCacheFlagCreatesCacheDir(t *testing.T) {
	dir := helpers.TmpDir(t)
	cachedir := filepath.Join(dir, "cache")
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := filepath.Join(dir, "out.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--cache", "--cache-dir", cachedir, in, "-o", out},
	})
	helpers.AssertExitCode(t, r, 0)
	if _, err := os.Stat(cachedir); err != nil {
		t.Fatalf("cache dir not created: %v", err)
	}
}

func TestCacheSecondRunUsesCache(t *testing.T) {
	dir := helpers.TmpDir(t)
	cachedir := filepath.Join(dir, "cache")
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := filepath.Join(dir, "out.js")

	// First run writes the cache.
	r1 := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--cache", "--cache-dir", cachedir, in, "-o", out},
	})
	helpers.AssertExitCode(t, r1, 0)

	if _, err := os.ReadDir(cachedir); err != nil {
		t.Fatal(err)
	}
	count1 := countFiles(cachedir)

	// Second run reads it and writes the same output.
	r2 := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--cache", "--cache-dir", cachedir, in, "-o", out},
	})
	helpers.AssertExitCode(t, r2, 0)

	count2 := countFiles(cachedir)

	if count1 != count2 {
		t.Fatalf("cache entry count changed: %d -> %d", count1, count2)
	}
}

func TestNoCacheDoesNotWriteCache(t *testing.T) {
	dir := helpers.TmpDir(t)
	cachedir := filepath.Join(dir, "cache")
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := filepath.Join(dir, "out.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--no-cache", "--cache-dir", cachedir, in, "-o", out},
	})
	helpers.AssertExitCode(t, r, 0)
	if _, err := os.Stat(cachedir); err == nil {
		entries, _ := os.ReadDir(cachedir)
		if len(entries) > 0 {
			t.Fatalf("--no-cache wrote entries: %v", entries)
		}
	}
}

func TestCacheDifferentContentDiffKey(t *testing.T) {
	dir := helpers.TmpDir(t)
	cachedir := filepath.Join(dir, "cache")
	in1 := helpers.WriteString(t, dir, "a.js", "const a = 1;\n")
	in2 := helpers.WriteString(t, dir, "b.js", "const b = 2;\n")
	out := filepath.Join(dir, "out.js")

	r1 := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--cache", "--cache-dir", cachedir, in1, "-o", out},
	})
	helpers.AssertExitCode(t, r1, 0)
	r2 := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--cache", "--cache-dir", cachedir, in2, "-o", out},
	})
	helpers.AssertExitCode(t, r2, 0)

	if countFiles(cachedir) < 2 {
		t.Fatalf("expected 2 cache entries, got %d", countFiles(cachedir))
	}
}

func TestCacheDirFlag(t *testing.T) {
	dir := helpers.TmpDir(t)
	custom := filepath.Join(dir, "custom-cache")
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := filepath.Join(dir, "out.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--cache", "--cache-dir", custom, in, "-o", out},
	})
	helpers.AssertExitCode(t, r, 0)
	if _, err := os.Stat(custom); err != nil {
		t.Fatalf("custom cache dir missing: %v", err)
	}
}

func TestCacheDisabledByDefault(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := filepath.Join(dir, "out.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{in, "-o", out},
		Dir:  dir, // run from a dir without a cache
	})
	helpers.AssertExitCode(t, r, 0)
	// The default cache dir is under the user's home; the test does
	// not touch it. We only check the run succeeded without --cache.
}

func TestCacheCorruptEntryRecovers(t *testing.T) {
	dir := helpers.TmpDir(t)
	cachedir := filepath.Join(dir, "cache")
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := filepath.Join(dir, "out.js")

	// First run creates a cache entry.
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--cache", "--cache-dir", cachedir, in, "-o", out},
	})
	helpers.AssertExitCode(t, r, 0)

	// Corrupt every entry.
	_ = filepath.Walk(cachedir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			_ = os.WriteFile(path, []byte("not json"), 0o644)
		}
		return nil
	})

	// Second run should still succeed (corrupt entry -> miss -> rebuild).
	r2 := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--cache", "--cache-dir", cachedir, in, "-o", out},
	})
	helpers.AssertExitCode(t, r2, 0)
}

func countFiles(root string) int {
	entries, err := os.ReadDir(root)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() {
			n++
		} else {
			n += countFiles(filepath.Join(root, e.Name()))
		}
	}
	return n
}