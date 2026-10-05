package bundle_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestTreeShakingRemovesUnusedExport(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "tree-shaking")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, entry},
	})
	helpers.AssertExitCode(t, r, 0)

	body, _ := os.ReadFile(out)
	if strings.Contains(string(body), "unusedFunction") {
		t.Fatalf("unused export survived tree shaking:\n%s", body)
	}
}

func TestTreeShakingKeepsUsedExport(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "tree-shaking")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, entry},
	})
	helpers.AssertExitCode(t, r, 0)

	body, _ := os.ReadFile(out)
	if !strings.Contains(string(body), `"used"`) {
		t.Fatalf("used export was removed:\n%s", body)
	}
}

func TestTreeShakingRemovesUnusedLocal(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js", `const unused = 1; const used = 2; console.log(used);`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--format", "esm",
			"--compress",
			"--outfile", out, filepath.Join(dir, "main.js"),
		},
	})
	helpers.AssertExitCode(t, r, 0)

	body, _ := os.ReadFile(out)
	if strings.Contains(string(body), "unused") {
		t.Fatalf("unused local survived:\n%s", body)
	}
}

func TestTreeShakingWithSideEffects(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js", `import "./side-effects.js"; console.log(1);`+"\n")
	helpers.WriteString(t, dir, "side-effects.js", `console.log("side effect");`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)

	body, _ := os.ReadFile(out)
	// Side-effect import must be preserved.
	if !strings.Contains(string(body), "side effect") {
		t.Fatalf("side-effect import was removed:\n%s", body)
	}
}

func TestTreeShakingPureAnnotation(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", `/* @__PURE__ */ export function unused() {}`+"\n")
	helpers.WriteString(t, dir, "main.js", `const x = 1;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestTreeShakingMinified(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "tree-shaking")
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

	body, _ := os.ReadFile(out)
	if strings.Contains(string(body), "unusedFunction") {
		t.Fatalf("unused export survived compress+mangle:\n%s", body)
	}
}

func TestTreeShakingExportStar(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "a.js", "export const used = 1; export const unused = 2;\n")
	helpers.WriteString(t, dir, "index.js", `export * from "./a.js";`+"\n")
	helpers.WriteString(t, dir, "main.js", `import { used } from "./index.js"; used;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestTreeShakingReExportChain(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "base.js", "export const v = 1;\n")
	helpers.WriteString(t, dir, "mid.js", `export { v } from "./base.js";`+"\n")
	helpers.WriteString(t, dir, "main.js", `import { v } from "./mid.js"; console.log(v);`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)

	body, _ := os.ReadFile(out)
	if !strings.Contains(string(body), "1") {
		t.Fatalf("re-export chain dropped the value:\n%s", body)
	}
}

func TestTreeShakingSideEffectFreePackage(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "package.json",
		`{"name":"test","version":"1.0.0","sideEffects":false}`)
	helpers.WriteString(t, dir, "lib.js", "export const unused = 1;\n")
	helpers.WriteString(t, dir, "main.js", `import { unused } from "./lib.js";`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestTreeShakingEmptyBundle(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", "export const x = 1;\n")
	helpers.WriteString(t, dir, "main.js", `import { x } from "./lib.js";`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}
