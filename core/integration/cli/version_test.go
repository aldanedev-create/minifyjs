package cli_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestVersionFormat(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--version"}})
	helpers.AssertExitCode(t, r, 0)
	out := strings.TrimSpace(string(r.Stdout))
	if !strings.HasPrefix(out, "minifyjs ") {
		t.Fatalf("version should start with 'minifyjs ': %q", out)
	}
	if !strings.Contains(out, "esbuild") {
		t.Fatalf("version should mention esbuild: %q", out)
	}
}

func TestVersionExitsBeforeReadingStdin(t *testing.T) {
	// Passing --version with stdin available should not read it.
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--version"},
		Stdin: []byte("this would be a syntax error"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoStderr(t, r)
}

func TestHelpExitsBeforeReadingStdin(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--help"},
		Stdin: []byte("this would be a syntax error"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestVersionNotAffectedByOtherFlags(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--version", "--target", "es5", "--mangle"},
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "minifyjs")
}

func TestVersionGoesToStdout(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--version"}})
	helpers.AssertNoStderr(t, r)
	if len(r.Stdout) == 0 {
		t.Fatal("version output empty")
	}
}

func TestVersionOneLine(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--version"}})
	lines := strings.Split(strings.TrimRight(string(r.Stdout), "\n"), "\n")
	if len(lines) != 1 {
		t.Fatalf("version should be one line, got %d: %q", len(lines), r.Stdout)
	}
}