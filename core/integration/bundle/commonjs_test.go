package bundle_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestBundleCJSSimple(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "cjs-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "cjs", "--outfile", out, entry},
	})
	helpers.AssertExitCode(t, r, 0)

	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) == 0 {
		t.Fatal("empty bundle output")
	}
}

func TestBundleCJSRequire(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", "module.exports = { v: 1 };\n")
	helpers.WriteString(t, dir, "main.js", `const lib = require("./lib.js"); lib.v;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "cjs", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleCJSExportsFunction(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", "module.exports = function() { return 1; };\n")
	helpers.WriteString(t, dir, "main.js", `const f = require("./lib.js"); f();`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "cjs", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleCJSExportsNamed(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", "exports.a = 1; exports.b = 2;\n")
	helpers.WriteString(t, dir, "main.js", `const lib = require("./lib.js"); lib.a + lib.b;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "cjs", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleCJSModuleDir(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "pkg/package.json", `{"main":"./index.js"}`)
	helpers.WriteString(t, dir, "pkg/index.js", "module.exports = 1;\n")
	helpers.WriteString(t, dir, "main.js", `const p = require("./pkg"); p;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "cjs", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleCJSMissingRequireFails(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js", `const x = require("./nope");`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "cjs", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	if r.ExitCode == 0 {
		t.Fatalf("missing require should fail")
	}
}

func TestBundleCJSWithExportsAndRequire(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "a.js", "module.exports = { x: 1 };\n")
	helpers.WriteString(t, dir, "b.js", `const a = require("./a"); module.exports = { y: a.x + 1 };`+"\n")
	helpers.WriteString(t, dir, "main.js", `const b = require("./b"); b.y;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "cjs", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleCJSWithTarget(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "cjs-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--format", "cjs",
			"--target", "es2015",
			"--outfile", out, entry,
		},
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleCJSNodePlatform(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "lib.js", "module.exports = 1;\n")
	helpers.WriteString(t, dir, "main.js", `require("./lib");`+"\n")
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

func TestBundleCJSMinified(t *testing.T) {
	dir := helpers.CopyFixtureDir(t, "cjs-simple")
	entry := filepath.Join(dir, "main.js")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle", "--format", "cjs",
			"--compress", "--mangle",
			"--outfile", out, entry,
		},
	})
	helpers.AssertExitCode(t, r, 0)
}