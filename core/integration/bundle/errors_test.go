package bundle_test

import (
	"path/filepath"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestBundleMissingEntryFails(t *testing.T) {
	dir := helpers.TmpDir(t)
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--outfile", out, filepath.Join(dir, "nope.js")},
	})
	if r.ExitCode == 0 {
		t.Fatalf("missing entry should fail")
	}
}

func TestBundleMissingImportFails(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js", `import x from "./nope.js"; x;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	if r.ExitCode == 0 {
		t.Fatalf("missing import should fail")
	}
}

func TestBundleSyntaxErrorFails(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js", "function () { } }")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	if r.ExitCode == 0 {
		t.Fatalf("syntax error should fail")
	}
}

func TestBundleWithoutOutputFails(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", filepath.Join(dir, "main.js")},
	})
	helpers.AssertExitCode(t, r, 2)
}

func TestBundleWithoutEntryFails(t *testing.T) {
	dir := helpers.TmpDir(t)
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--outfile", out},
	})
	helpers.AssertExitCode(t, r, 2)
}

func TestBundleOutfileAndOutdirConflict(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle",
			"--outfile", filepath.Join(dir, "out.js"),
			"--outdir", filepath.Join(dir, "dist"),
			filepath.Join(dir, "main.js"),
		},
	})
	// Using both is ambiguous.
	if r.ExitCode == 0 {
		t.Fatalf("both outfile and outdir should fail")
	}
}

func TestBundleImportOutsideRoot(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js", `import x from "/absolute/path.js"; x;`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	if r.ExitCode == 0 {
		t.Fatalf("missing absolute path should fail")
	}
}

func TestBundleImportCircularBrokenPath(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "a.js", `import "./b.js"; export const a = 1;`+"\n")
	// b.js missing.
	helpers.WriteString(t, dir, "main.js", `import "./a.js";`+"\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--format", "esm", "--outfile", out, filepath.Join(dir, "main.js")},
	})
	if r.ExitCode == 0 {
		t.Fatalf("missing file in chain should fail")
	}
}

func TestBundleNoEntryButInputFile(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js", "const x = 1;\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--bundle", "--outfile", out, filepath.Join(dir, "main.js")},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBundleInvalidFormat(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "main.js", "const x = 1;\n")
	out := filepath.Join(dir, "bundle.js")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--bundle",
			"--format", "bogus",
			"--outfile", out,
			filepath.Join(dir, "main.js"),
		},
	})
	if r.ExitCode == 0 {
		t.Fatalf("invalid format should fail")
	}
}