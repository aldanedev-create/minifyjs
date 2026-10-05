package cli_test

import (
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestDefineBooleanFalse(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--define", "DEBUG=false", "--compress"},
		Stdin: []byte("if (DEBUG) { console.log('debug'); }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "DEBUG")
	helpers.AssertNotContains(t, string(r.Stdout), "console")
}

func TestDefineBooleanTrue(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--define", "DEBUG=true", "--compress"},
		Stdin: []byte("if (DEBUG) { const y = 1; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "DEBUG")
}

func TestDefineString(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--define", `VERSION="1.0.0"`},
		Stdin: []byte("console.log(VERSION);"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "1.0.0")
}

func TestDefineNumber(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--define", "PORT=8080"},
		Stdin: []byte("const p = PORT;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "8080")
}

func TestDefineObjectLiteral(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--define", `CONFIG={"debug":true}`},
		Stdin: []byte("const c = CONFIG.debug;"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestDefineWithoutEqualsExitsTwo(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--define", "noequalsign"},
		Stdin: []byte("const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 2)
}

func TestDropConsole(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--drop", "console", "--compress"},
		Stdin: []byte("console.log('a'); console.error('b');"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "console")
}

func TestDropDebugger(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--drop", "debugger", "--compress"},
		Stdin: []byte("debugger; const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "debugger")
}

func TestDropInvalidKindExitsTwo(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--drop", "banana"},
		Stdin: []byte("const x = 1;"),
	})
	if r.ExitCode == 0 {
		t.Fatalf("invalid --drop kind should be rejected")
	}
}

func TestPureRemovesCall(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--pure", "console.log", "--compress"},
		Stdin: []byte("const x = 1; console.log(x);"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "console.log")
}

func TestPureMultipleCalls(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{
			"--pure", "console.log",
			"--pure", "console.warn",
			"--compress",
		},
		Stdin: []byte("const x = 1; console.log(x); console.warn(x);"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "console.log")
	helpers.AssertNotContains(t, string(r.Stdout), "console.warn")
}

func TestDefineAndCompressTogether(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--define", "A=1", "--define", "B=2", "--compress"},
		Stdin: []byte("const x = A + B;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "3")
}