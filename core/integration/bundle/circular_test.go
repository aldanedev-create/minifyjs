package bundle_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestCircularImportsHandled(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "circular")
	entry := filepath.Join(dir, "a.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, entry},
	})
	helpers.AssertExitCode(t, r, 0)

	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) == 0 {
		t.Fatal("empty output for circular imports")
	}
}

func TestCircularSelfImport(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "a.js", `import { b } from "./a.js"; export const a = 1;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "a.js")},
	})
	// esbuild resolves self-imports; either success or a clear error
	// is acceptable.
	if r.ExitCode == 4 {
		t.Fatalf("self import crashed")
	}
}

func TestCircularLongChain(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "a.js", `import "./b.js"; export const a = 1;`+"\n")
	helpers.WriteString(t, dir, "b.js", `import "./c.js"; export const b = 2;`+"\n")
	helpers.WriteString(t, dir, "c.js", `import "./a.js"; export const c = 3;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "a.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestCircularMinified(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "circular")
	entry := filepath.Join(dir, "a.js")
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

func TestCircularCommonJS(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "a.js", `const b = require("./b"); module.exports = { a: 1 };`+"\n")
	helpers.WriteString(t, dir, "b.js", `const a = require("./a"); module.exports = { b: 2 };`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "cjs", "--outfile", out, filepath.Join(dir, "a.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestCircularWithSplitting(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "a.js", `import "./b.js"; export const a = 1;`+"\n")
	helpers.WriteString(t, dir, "b.js", `import "./a.js"; export const b = 2;`+"\n")
	outdir := filepath.Join(dir, "dist")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--splitting",
			"--format", "esm",
			"--outdir", outdir,
			filepath.Join(dir, "a.js"),
		},
	})
	helpers.AssertExitCode(t, r, 0)
}