package cli_test

import (
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

// These tests exist to lock in the exit codes documented in
// docs/reference/exit-codes.md. If you change a code, change the
// docs and this file together.

func TestExitZeroOnSuccess(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const x = 1;")})
	helpers.AssertExitCode(t, r, 0)
}

func TestExitOneOnSyntaxError(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("function () { } }")})
	helpers.AssertExitCode(t, r, 1)
}

func TestExitOneOnInvalidConfig(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{invalid}`)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{Args: []string{"in.js"}, Dir: dir})
	helpers.AssertExitCode(t, r, 1)
}

func TestExitTwoOnUnknownFlag(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--nope"}})
	helpers.AssertExitCode(t, r, 2)
}

func TestExitTwoOnMissingValue(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--target"}})
	helpers.AssertExitCode(t, r, 2)
}

func TestExitTwoOnTooManyInputs(t *testing.T) {
	dir := helpers.TmpDir(t)
	a := helpers.WriteString(t, dir, "a.js", "const a = 1;")
	b := helpers.WriteString(t, dir, "b.js", "const b = 2;")

	r := helpers.Run(t, helpers.RunOptions{Args: []string{a, b}})
	helpers.AssertExitCode(t, r, 2)
}

func TestExitThreeOnMissingFile(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"/no/such/file.js"}})
	helpers.AssertExitCode(t, r, 3)
}

func TestExitZeroOnHelp(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--help"}})
	helpers.AssertExitCode(t, r, 0)
}

func TestExitZeroOnVersion(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--version"}})
	helpers.AssertExitCode(t, r, 0)
}

func TestExitOneOnUnwritableOutput(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;")
	blocker := helpers.WriteString(t, dir, "blocker", "x")
	out := blocker + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{Args: []string{in, "-o", out}})
	if r.ExitCode != 1 {
		t.Fatalf("exit = %d, want 1", r.ExitCode)
	}
}