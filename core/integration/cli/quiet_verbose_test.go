package cli_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestQuietSuppressesProgress(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{Args: []string{"-q", in, "-o", out}})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoStderr(t, r)
}

func TestQuietLongFlag(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--quiet", in, "-o", out}})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoStderr(t, r)
}

func TestQuietStillPrintsErrors(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--quiet"},
		Stdin: []byte("function () { } }"),
	})
	if r.ExitCode == 0 {
		t.Fatalf("expected error")
	}
	if len(r.Stderr) == 0 {
		t.Fatal("errors should still print under --quiet")
	}
	helpers.AssertContains(t, strings.ToLower(string(r.Stderr)), "error")
}

func TestVerboseAddsOutput(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--verbose", in, "-o", out}})
	helpers.AssertExitCode(t, r, 0)
	// Verbose should emit something to stderr; the exact wording is
	// not part of the stable API.
	if len(r.Stderr) == 0 {
		t.Fatalf("--verbose should print extra output")
	}
}

func TestQuietAndVerboseConflict(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--quiet", "--verbose", in},
	})
	// Either the CLI errors or one wins. Both are acceptable, but
	// the CLI must not crash.
	if r.ExitCode == 4 {
		t.Fatalf("--quiet --verbose crashed")
	}
}

func TestNoQuietOverridesConfig(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{"quiet":true}`)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--no-quiet", in, "-o", out},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	if len(r.Stderr) == 0 {
		t.Fatalf("--no-quiet should re-enable progress output")
	}
}

func TestNoVerboseDisables(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{"verbose":true}`)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--no-verbose", in, "-o", out},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestQuietEnvironmentVariable(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := dir + "/out.js"

	env := []string{"MINIFYJS_QUIET=1", "PATH=" + pathEnv()}
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{in, "-o", out},
		Env:  env,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoStderr(t, r)
}