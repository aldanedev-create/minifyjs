package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestBundleSimpleESM(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "dist", "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--outfile", out, entry},
	})
	helpers.AssertExitCode(t, r, 0)
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("bundled output missing: %v", err)
	}
}

func TestBundleWithOutdir(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	outdir := filepath.Join(dir, "dist")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--outdir", outdir, entry},
	})
	helpers.AssertExitCode(t, r, 0)
	entries, err := os.ReadDir(outdir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("no output files written")
	}
}

func TestBundleRequiresOutput(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", entry},
	})
	helpers.AssertExitCode(t, r, 2)
}

func TestBundleRequiresEntry(t *testing.T) {
	dir := helpers.TmpDir(t)
	outdir := filepath.Join(dir, "dist")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--outdir", outdir},
	})
	helpers.AssertExitCode(t, r, 2)
}

func TestBundleTreeShaking(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "tree-shaking")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--outfile", out, "--format", "esm", entry},
	})
	helpers.AssertExitCode(t, r, 0)

	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	// The tree-shaking fixture has an unused export. It should be gone.
	if strings.Contains(string(body), "unusedFunction") {
		t.Fatalf("tree shaking did not remove unused export:\n%s", body)
	}
}

func TestBundleMinificationApplied(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--outfile", out,
			"--compress", "--mangle",
			"--format", "esm",
			entry,
		},
	})
	helpers.AssertExitCode(t, r, 0)

	body, _ := os.ReadFile(out)
	if strings.Contains(string(body), "\n    ") {
		t.Fatalf("bundle output does not look minified:\n%s", body)
	}
}

func TestBundlePlatformBrowser(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--platform", "browser", "--outfile", out, "--format", "esm", entry},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundlePlatformNode(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--platform", "node", "--outfile", out, "--format", "esm", entry},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleSplittingRequiresESM(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "multi-entry")
	entry := filepath.Join(dir, "a.js")
	outdir := filepath.Join(dir, "dist")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--splitting",
			"--outdir", outdir,
			"--format", "cjs", // wrong on purpose
			entry,
		},
	})
	if r.ExitCode == 0 {
		t.Fatalf("splitting with non-esm format should fail")
	}
}

func TestBundleMissingEntry(t *testing.T) {
	dir := helpers.TmpDir(t)
	entry := filepath.Join(dir, "does-not-exist.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--outfile", out, entry},
	})
	if r.ExitCode == 0 {
		t.Fatalf("missing entry should fail")
	}
}

func TestBundleSyntaxError(t *testing.T) {
	dir := helpers.TmpDir(t)
	entry := helpers.WriteString(t, dir, "bad.js", "function () { } }")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--outfile", out, entry},
	})
	if r.ExitCode == 0 {
		t.Fatalf("syntax error should fail")
	}
}

func TestBundleWithTarget(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--outfile", out,
			"--format", "esm",
			"--target", "es2015",
			entry,
		},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleSourcemap(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--outfile", out,
			"--sourcemap=external",
			"--format", "esm",
			entry,
		},
	})
	helpers.AssertExitCode(t, r, 0)
	if _, err := os.Stat(out + ".map"); err != nil {
		t.Fatalf("bundle map missing: %v", err)
	}
}