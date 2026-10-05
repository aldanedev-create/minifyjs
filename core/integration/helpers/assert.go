// Package helpers: assertion helpers. Go's testing package has no
// built-in equality assertions; these fill the gap with good
// failure messages for the common cases.
package helpers

import (
	"os"
	"bytes"
	"strings"
	"testing"
)

// AssertEqual fails the test if got != want, printing both.
func AssertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("got  %q\nwant %q", got, want)
	}
}


// in core/integration/helpers/assert.go, append:

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}


func pathEnv() string {
	return os.Getenv("PATH")
}

// AssertEqualBytes is AssertEqual for byte slices. It compares by
// value, not by identity.
func AssertEqualBytes(t *testing.T, got, want []byte) {
	t.Helper()
	if !bytes.Equal(got, want) {
		t.Fatalf("got  %q\nwant %q", got, want)
	}
}

// AssertContains fails the test if haystack does not contain needle.
func AssertContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected %q to contain %q", haystack, needle)
	}
}

// AssertNotContains fails the test if haystack contains needle.
func AssertNotContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if strings.Contains(haystack, needle) {
		t.Fatalf("expected %q not to contain %q", haystack, needle)
	}
}

// AssertExitCode fails the test if the CLI result's exit code differs
// from want, printing stderr for context.
func AssertExitCode(t *testing.T, r Result, want int) {
	t.Helper()
	if r.ExitCode != want {
		t.Fatalf("exit = %d, want %d\nstdout: %s\nstderr: %s",
			r.ExitCode, want, r.Stdout, r.Stderr)
	}
}

// AssertNoNewlines fails if s contains \n or \r.
func AssertNoNewlines(t *testing.T, s string) {
	t.Helper()
	if strings.ContainsAny(s, "\r\n") {
		t.Fatalf("expected no newlines, got %q", s)
	}
}

// AssertNoStderr fails if stderr is non-empty.
func AssertNoStderr(t *testing.T, r Result) {
	t.Helper()
	if len(r.Stderr) > 0 {
		t.Fatalf("expected empty stderr, got %q", r.Stderr)
	}
}

// AssertNoStdout fails if stdout is non-empty.
func AssertNoStdout(t *testing.T, r Result) {
	t.Helper()
	if len(r.Stdout) > 0 {
		t.Fatalf("expected empty stdout, got %q", r.Stdout)
	}
}