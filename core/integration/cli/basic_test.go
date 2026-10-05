package cli_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestHelpExitsZero(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--help"}})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "Usage:")
}

func TestHelpShortFlag(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"-h"}})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "Usage:")
}

func TestHelpListsAllFlags(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--help"}})
	helpers.AssertExitCode(t, r, 0)
	for _, flag := range []string{
		"-o, --output",
		"--compress",
		"--mangle",
		"--target",
		"--format",
		"--sourcemap",
		"--banner",
		"--footer",
		"--bundle",
		"--cache",
		"--config",
		"--watch",
		"-q, --quiet",
		"-v, --verbose",
		"-V, --version",
	} {
		if !strings.Contains(string(r.Stdout), flag) {
			t.Errorf("--help does not mention %q", flag)
		}
	}
}

func TestHelpGoesToStdoutNotStderr(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--help"}})
	if len(r.Stderr) != 0 {
		t.Fatalf("stderr should be empty: %q", r.Stderr)
	}
}

func TestVersionExitsZero(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--version"}})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "minifyjs")
}

func TestVersionReportsEsbuildVersion(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--version"}})
	helpers.AssertContains(t, string(r.Stdout), "esbuild")
}

func TestVersionShortFlag(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"-V"}})
	helpers.AssertExitCode(t, r, 0)
}

func TestVersionGoesToStdoutNotStderr(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--version"}})
	if len(r.Stderr) != 0 {
		t.Fatalf("stderr should be empty: %q", r.Stderr)
	}
}

func TestUnknownFlagExitsTwo(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--not-a-flag"}})
	helpers.AssertExitCode(t, r, 2)
	helpers.AssertContains(t, string(r.Stderr), "unknown flag")
}

func TestMissingValueExitsTwo(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--target"}})
	helpers.AssertExitCode(t, r, 2)
}

func TestNoArgsEmptyStdinProducesEmptyStdout(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("")})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoStdout(t, r)
}

func TestDoubleDashEndsFlagParsing(t *testing.T) {
	// After "--", everything is a positional argument. A file
	// literally named "--weird.js" should be read as a file.
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "--weird.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--", "--weird.js"},
		Dir:  dir,
	})
	if r.ExitCode == 2 {
		t.Fatalf("-- should have ended flag parsing")
	}
	helpers.AssertExitCode(t, r, 0)
}