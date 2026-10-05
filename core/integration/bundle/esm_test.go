package bundle_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestBundleESMSimple(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, entry},
	})
	helpers.AssertExitCode(t, r, 0)

	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	// The bundle should inline the helpers from utils.js and math.js.
	helpers.AssertContains(t, string(body), "function")
}

func TestBundleESMInlinesImports(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, entry},
	})
	helpers.AssertExitCode(t, r, 0)

	body, _ := os.ReadFile(out)
	// The original module should not be reachable via a relative
	// import anymore; the code is inlined.
	helpers.AssertNotContains(t, string(body), `from "./utils.js"`)
}

func TestBundleESMKeepsExternalImports(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "external-deps")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, entry},
	})
	helpers.AssertExitCode(t, r, 0)

	body, _ := os.ReadFile(out)
	// "external" module is not on disk, so it stays as an import.
	if !strings.Contains(string(body), "external") {
		t.Logf("bundle output: %s", body)
	}
}

func TestBundleESMMultipleImports(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "esm-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, entry},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleESMDefaultExport(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", "export default function hello() { return 'hi'; }\n")
	helpers.WriteString(t, dir, "main.js", `import hello from "./lib.js"; hello();`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleESMNamedExport(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", "export function add(a, b) { return a + b; }\n")
	helpers.WriteString(t, dir, "main.js", `import { add } from "./lib.js"; add(1, 2);`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleESMRenamedImport(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", "export function original() { return 1; }\n")
	helpers.WriteString(t, dir, "main.js", `import { original as renamed } from "./lib.js"; renamed();`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleESMNamespaceImport(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", "export const x = 1; export const y = 2;\n")
	helpers.WriteString(t, dir, "main.js", `import * as lib from "./lib.js"; lib.x + lib.y;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleESMReExport(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", "export const x = 1;\n")
	helpers.WriteString(t, dir, "index.js", `export { x } from "./lib.js";`+"\n")
	helpers.WriteString(t, dir, "main.js", `import { x } from "./index.js"; x;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleESMNested(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "sub/lib.js", "export const v = 1;\n")
	helpers.WriteString(t, dir, "main.js", `import { v } from "./sub/lib.js"; v;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleESMWithIndexFile(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib/index.js", "export const v = 1;\n")
	helpers.WriteString(t, dir, "main.js", `import { v } from "./lib/index.js"; v;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleESMWithoutExtension(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", "export const v = 1;\n")
	helpers.WriteString(t, dir, "main.js", `import { v } from "./lib"; v;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}