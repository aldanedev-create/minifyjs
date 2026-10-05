package bundle_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestBundleSourcemapExternal(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--format", "esm",
			"--sourcemap=external",
			"--outfile", out, entry,
		},
	})
	helpers.AssertExitCode(t, r, 0)

	if _, err := os.Stat(out + ".map"); err != nil {
		t.Fatalf("map file missing: %v", err)
	}
}

func TestBundleSourcemapListsAllSources(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--format", "esm",
			"--sourcemap=external",
			"--outfile", out, entry,
		},
	})
	helpers.AssertExitCode(t, r, 0)

	data, _ := os.ReadFile(out + ".map")
	body := string(data)
	// The map's sources array should reference the original files.
	if !strings.Contains(body, "main.js") {
		t.Fatalf("map does not reference main.js:\n%s", body)
	}
}

func TestBundleSourcemapInline(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--format", "esm",
			"--sourcemap=inline",
			"--outfile", out, entry,
		},
	})
	helpers.AssertExitCode(t, r, 0)

	body, _ := os.ReadFile(out)
	helpers.AssertContains(t, string(body), "sourceMappingURL=data:application/json")
}

func TestBundleSourcemapWithSplitting(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "multi-entry")
	outdir := filepath.Join(dir, "dist")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--splitting",
			"--format", "esm",
			"--sourcemap=external",
			"--outdir", outdir,
			filepath.Join(dir, "a.js"),
			filepath.Join(dir, "b.js"),
		},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleSourcemapWithMinify(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--format", "esm",
			"--compress", "--mangle",
			"--sourcemap=external",
			"--outfile", out, entry,
		},
	})
	helpers.AssertExitCode(t, r, 0)

	if _, err := os.Stat(out + ".map"); err != nil {
		t.Fatalf("map missing: %v", err)
	}
}

func TestBundleSourcemapNotEmittedByDefault(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, entry},
	})
	helpers.AssertExitCode(t, r, 0)

	if _, err := os.Stat(out + ".map"); err == nil {
		t.Fatalf("map should not be written by default")
	}
}