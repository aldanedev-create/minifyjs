package bundle_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestDynamicImportPreservedWithoutSplitting(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lazy.js", "export const v = 1;\n")
	helpers.WriteString(t, dir, "main.js", `import("./lazy.js").then(m => m.v);`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)

	body, _ := os.ReadFile(out)
	if !strings.Contains(string(body), "import(") {
		t.Fatalf("dynamic import was removed:\n%s", body)
	}
}

func TestDynamicImportMultiple(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "a.js", "export const a = 1;\n")
	helpers.WriteString(t, dir, "b.js", "export const b = 2;\n")
	helpers.WriteString(t, dir, "main.js",
		`Promise.all([import("./a.js"), import("./b.js")]);`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestDynamicImportVariable(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js",
		`const name = "./a.js"; import(name);`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestDynamicImportMinified(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lazy.js", "export const v = 1;\n")
	helpers.WriteString(t, dir, "main.js", `import("./lazy.js").then(m => m.v);`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--format", "esm",
			"--compress", "--mangle",
			"--outfile", out, filepath.Join(dir, "main.js"),
		},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestDynamicImportWithAwait(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lazy.js", "export const v = 1;\n")
	helpers.WriteString(t, dir, "main.js",
		`async function f() { const m = await import("./lazy.js"); return m.v; }`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestDynamicImportMissingFileFails(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js", `import("./nope.js");`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 1)
}

func TestDynamicImportNested(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "a.js", `import("./b.js"); export const a = 1;`+"\n")
	helpers.WriteString(t, dir, "b.js", "export const b = 2;\n")
	helpers.WriteString(t, dir, "main.js", `import("./a.js");`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestDynamicImportCommonJSLowered(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lazy.js", "export const v = 1;\n")
	helpers.WriteString(t, dir, "main.js", `import("./lazy.js");`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--format", "cjs",
			"--platform", "node",
			"--outfile", out, filepath.Join(dir, "main.js"),
		},
	})
	helpers.AssertExitCode(t, r, 0)
}