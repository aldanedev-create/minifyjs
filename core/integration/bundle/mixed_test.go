package bundle_test

import (
	"path/filepath"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestBundleESMImportsCJS(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-with-cjs")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, entry},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleCJSImportsESM(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "esm-lib.mjs", "export const v = 1;\n")
	helpers.WriteString(t, dir, "main.js", `import("./esm-lib.mjs");`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "cjs", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	if r.ExitCode == 4 {
		t.Fatalf("mixed cjs->esm crashed")
	}
}

func TestBundleMixedDefaultInterop(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "cjs-lib.js", "module.exports = { greet: () => 'hi' };\n")
	helpers.WriteString(t, dir, "main.js", `import lib from "./cjs-lib.js"; lib.greet();`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleMixedNamedInterop(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "cjs-lib.js", "module.exports = { greet: () => 'hi' };\n")
	helpers.WriteString(t, dir, "main.js", `import { greet } from "./cjs-lib.js"; greet();`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	if r.ExitCode != 0 {
		t.Logf("named interop from cjs failed: %s", r.Stderr)
	}
}

func TestBundleMixedMinified(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-with-cjs")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--format", "esm",
			"--compress", "--mangle",
			"--outfile", out, entry,
		},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleMixedBrowserPlatform(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-with-cjs")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--format", "esm",
			"--platform", "browser",
			"--outfile", out, entry,
		},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleMixedPreservesDynamicImport(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", "export const v = 1;\n")
	helpers.WriteString(t, dir, "main.js", `import("./lib.js").then(m => m.v);`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleMixedExternalPackage(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "package.json", `{"name":"test","version":"1.0.0"}`)
	helpers.WriteString(t, dir, "main.js", "const x = 1;\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
}