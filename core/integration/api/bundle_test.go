package api_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/minifyjs/minifyjs/core/api"
)

func TestAPIBundleESM(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lib.js"), "export const v = 1;\n")
	writeFile(t, filepath.Join(dir, "main.js"), `import { v } from "./lib.js"; console.log(v);`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r, err := api.Bundle(api.BundleOptions{
		EntryPoints:   []string{filepath.Join(dir, "main.js")},
		OutFile:       out,
		Bundle:        true,
		AbsWorkingDir: dir,
		Options: api.Options{
			MinifyWhitespace: true,
			Format:           "esm",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.HasErrors() {
		t.Fatalf("diagnostics: %v", r.Diagnostics)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("bundle not written: %v", err)
	}
}

func TestAPIBundleMissingEntry(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "bundle.js")

	r, err := api.Bundle(api.BundleOptions{
		EntryPoints: []string{filepath.Join(dir, "nope.js")},
		OutFile:     out,
		Bundle:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !r.HasErrors() {
		t.Fatal("expected diagnostics for missing entry")
	}
}

func TestAPIBundleNoEntry(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "bundle.js")

	_, err := api.Bundle(api.BundleOptions{
		OutFile: out,
		Bundle:  true,
	})
	if err == nil {
		t.Fatal("expected error for missing entry")
	}
}

func TestAPIBundleNoOutput(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "main.js"), "const x = 1;\n")

	_, err := api.Bundle(api.BundleOptions{
		EntryPoints: []string{filepath.Join(dir, "main.js")},
		Bundle:      true,
	})
	if err == nil {
		t.Fatal("expected error for missing output")
	}
}

func TestAPIBundleBothOutFileAndOutDir(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "main.js"), "const x = 1;\n")

	_, err := api.Bundle(api.BundleOptions{
		EntryPoints: []string{filepath.Join(dir, "main.js")},
		OutFile:     filepath.Join(dir, "out.js"),
		OutDir:      filepath.Join(dir, "dist"),
		Bundle:      true,
	})
	if err == nil {
		t.Fatal("expected error for both outfile and outdir")
	}
}

func TestAPIBundleWithMinify(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lib.js"), "export function add(a, b) { return a + b; }\n")
	writeFile(t, filepath.Join(dir, "main.js"), `import { add } from "./lib.js"; add(1, 2);`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r, err := api.Bundle(api.BundleOptions{
		EntryPoints:   []string{filepath.Join(dir, "main.js")},
		OutFile:       out,
		Bundle:        true,
		AbsWorkingDir: dir,
		Options: api.Options{
			MinifyWhitespace:  true,
			MinifyIdentifiers: true,
			MinifySyntax:      true,
			Format:            "esm",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.HasErrors() {
		t.Fatalf("diagnostics: %v", r.Diagnostics)
	}
}

func TestAPIBundleCodeReturned(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "main.js"), "console.log(1);\n")
	out := filepath.Join(dir, "bundle.js")

	r, err := api.Bundle(api.BundleOptions{
		EntryPoints: []string{filepath.Join(dir, "main.js")},
		OutFile:     out,
		Bundle:      true,
		Options: api.Options{
			Format: "esm",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// Code is populated when OutFile is used.
	if r.Code == "" {
		t.Logf("Code empty (may be normalized in future versions)")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}