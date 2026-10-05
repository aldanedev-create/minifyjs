package engine

import (
	"os"
	"path/filepath"
	"testing"
)

// Build writes to the filesystem, so these tests use t.TempDir and
// never touch the real working tree.

func TestBuildSimpleBundle(t *testing.T) {
	dir := t.TempDir()

	// Two-module project: main.js imports utils.js.
	writeFile(t, filepath.Join(dir, "utils.js"), `
export function add(a, b) { return a + b; }
`)
	writeFile(t, filepath.Join(dir, "main.js"), `
import { add } from "./utils.js";
console.log(add(1, 2));
`)

	outDir := filepath.Join(dir, "dist")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}

	res, err := Build(BuildOptions{
		EntryPoints: []string{filepath.Join(dir, "main.js")},
		OutDir:      outDir,
		Bundle:      true,
		Format:      "esm",
		AbsWorkingDir: dir,
		Options: Options{
			MinifyWhitespace:  true,
			MinifyIdentifiers: true,
			MinifySyntax:      true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.HasErrors() {
		t.Fatalf("diagnostics: %v", res.Diagnostics)
	}

	// The bundled output should exist on disk.
	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one output file")
	}
}

func TestBuildMissingImportReportsDiagnostic(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "main.js"), `
import { nope } from "./does-not-exist.js";
`)
	outDir := filepath.Join(dir, "dist")
	_ = os.MkdirAll(outDir, 0o755)

	res, err := Build(BuildOptions{
		EntryPoints:   []string{filepath.Join(dir, "main.js")},
		OutDir:        outDir,
		Bundle:        true,
		Format:        "esm",
		AbsWorkingDir: dir,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.HasErrors() {
		t.Fatal("expected a diagnostic for the missing import")
	}
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}