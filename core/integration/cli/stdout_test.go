package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestStdoutByteExactNoTrailingNewline(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const x = 1;")})
	helpers.AssertExitCode(t, r, 0)
	if bytes.HasSuffix(r.Stdout, []byte("\n")) {
		t.Fatalf("stdout has trailing newline: %q", r.Stdout)
	}
}

func TestStdoutEndsWithSemicolon(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const x = 1;")})
	helpers.AssertExitCode(t, r, 0)
	if r.Stdout[len(r.Stdout)-1] != ';' {
		t.Fatalf("stdout should end with ';': %q", r.Stdout)
	}
}

func TestStdoutOnSyntaxErrorIsEmpty(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("function () { } }")})
	helpers.AssertNoStdout(t, r)
}

func TestStdoutOnMissingFileIsEmpty(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"/does/not/exist.js"}})
	helpers.AssertNoStdout(t, r)
}

func TestQuietSuppressesOutputFilePath(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{Args: []string{"-q", in, "-o", out}})
	helpers.AssertExitCode(t, r, 0)
	if len(r.Stderr) != 0 {
		t.Fatalf("--quiet should suppress stderr: %q", r.Stderr)
	}
}

func TestNonQuietReportsSavings(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "function add(a, b) {\n    return a + b;\n}\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{Args: []string{in, "-o", out}})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stderr), "bytes")
}

func TestStdoutExactBytes(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = 1;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertEqual(t, string(r.Stdout), "const x = 1;\n")
}

func TestStdoutLargeOutput(t *testing.T) {
	var b strings.Builder
	for i := 0; i < 2000; i++ {
		b.WriteString("const a")
		b.WriteString(itoa(i))
		b.WriteString(" = ")
		b.WriteString(itoa(i))
		b.WriteString(";\n")
	}
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte(b.String())})
	helpers.AssertExitCode(t, r, 0)
	if len(r.Stdout) == 0 {
		t.Fatal("expected non-empty stdout")
	}
}

func TestStdoutNoBOMAdded(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const x = 1;")})
	helpers.AssertExitCode(t, r, 0)
	if bytes.HasPrefix(r.Stdout, []byte("\xef\xbb\xbf")) {
		t.Fatalf("stdout should not start with BOM: %q", r.Stdout)
	}
}

func TestStdoutPreservesUnicode(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = \"café\";\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "café")
}

func TestStdoutQuietIsQuiet(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := dir + "/out.js"
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--quiet", in, "-o", out}})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoStderr(t, r)
}