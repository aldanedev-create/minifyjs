package bundle_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestSplittingCreatesMultipleFiles(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "multi-entry")
	outdir := filepath.Join(dir, "dist")

	entryA := filepath.Join(dir, "a.js")
	entryB := filepath.Join(dir, "b.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--splitting",
			"--format", "esm",
			"--outdir", outdir,
			entryA, entryB,
		},
	})
	helpers.AssertExitCode(t, r, 0)

	entries, err := os.ReadDir(outdir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 2 {
		t.Fatalf("expected multiple output files, got %d", len(entries))
	}
}

func TestSplittingRequiresESMFormat(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "multi-entry")
	outdir := filepath.Join(dir, "dist")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--splitting",
			"--format", "cjs",
			"--outdir", outdir,
			filepath.Join(dir, "a.js"),
		},
	})
	if r.ExitCode == 0 {
		t.Fatalf("splitting without esm should fail")
	}
}

func TestSplittingRequiresBundle(t *testing.T) {
	dir := helpers.TmpDir(t)
	outdir := filepath.Join(dir, "dist")
	helpers.WriteString(t, dir, "a.js", "const a = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--splitting",
			"--format", "esm",
			"--outdir", outdir,
			filepath.Join(dir, "a.js"),
		},
	})
	if r.ExitCode == 0 {
		t.Fatalf("splitting without --bundle should fail")
	}
}

func TestSplittingSharedChunk(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "shared.js", "export const shared = 1;\n")
	helpers.WriteString(t, dir, "a.js", `import { shared } from "./shared.js"; console.log("a", shared);`+"\n")
	helpers.WriteString(t, dir, "b.js", `import { shared } from "./shared.js"; console.log("b", shared);`+"\n")
	outdir := filepath.Join(dir, "dist")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--splitting",
			"--format", "esm",
			"--outdir", outdir,
			filepath.Join(dir, "a.js"),
			filepath.Join(dir, "b.js"),
		},
	})
	helpers.AssertExitCode(t, r, 0)

	entries, _ := os.ReadDir(outdir)
	if len(entries) < 3 {
		t.Fatalf("expected shared chunk, got %d files", len(entries))
	}
}

func TestSplittingDynamicImport(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lazy.js", "export const v = 1;\n")
	helpers.WriteString(t, dir, "main.js", `import("./lazy.js").then(m => m.v);`+"\n")
	outdir := filepath.Join(dir, "dist")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--splitting",
			"--format", "esm",
			"--outdir", outdir,
			filepath.Join(dir, "main.js"),
		},
	})
	helpers.AssertExitCode(t, r, 0)

	entries, _ := os.ReadDir(outdir)
	if len(entries) < 2 {
		t.Fatalf("expected dynamic chunk, got %d files", len(entries))
	}
}

func TestSplittingMinified(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "multi-entry")
	outdir := filepath.Join(dir, "dist")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--splitting",
			"--format", "esm",
			"--compress", "--mangle",
			"--outdir", outdir,
			filepath.Join(dir, "a.js"),
			filepath.Join(dir, "b.js"),
		},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestSplittingWithSourcemap(t *testing.T) {
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

func TestSplittingSingleEntry(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js", `import("./lazy.js");`+"\n")
	helpers.WriteString(t, dir, "lazy.js", "export const v = 1;\n")
	outdir := filepath.Join(dir, "dist")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--splitting",
			"--format", "esm",
			"--outdir", outdir,
			filepath.Join(dir, "main.js"),
		},
	})
	helpers.AssertExitCode(t, r, 0)
}