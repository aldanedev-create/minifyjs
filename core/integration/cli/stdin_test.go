package cli_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestStdinPipe(t *testing.T) {
	src := []byte("function add(a, b) {\n    return a + b;\n}\n")
	r := helpers.Run(t, helpers.RunOptions{Stdin: src})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoNewlines(t, string(r.Stdout))
	helpers.AssertContains(t, string(r.Stdout), "function")
}

func TestStdinDashReadsFromStdin(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"-"},
		Stdin: []byte("const x = 1;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "const x")
}

func TestStdinEmpty(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("")})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoStdout(t, r)
}

func TestStdinWhitespaceOnly(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("   \n\t\n  ")})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoStdout(t, r)
}

func TestStdinCommentsOnly(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("// only a comment\n/* and a block */\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoStdout(t, r)
}

func TestStdinWithBOM(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("\xef\xbb\xbfconst x = 1;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "const x")
}

func TestStdinCRLF(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("function add(a, b) {\r\n    return a + b;\r\n}\r\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "\r")
	helpers.AssertNotContains(t, string(r.Stdout), "\n")
}

func TestStdinSyntaxError(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("function () { } }")})
	if r.ExitCode == 0 {
		t.Fatalf("syntax error should exit non-zero")
	}
	helpers.AssertContains(t, strings.ToLower(string(r.Stderr)), "error")
}

func TestStdinNoMinifyFlag(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = 1;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "const x")
}

func TestStdinFromFixture(t *testing.T) {
	src := helpers.InputFixture(t, "simple.js")
	r := helpers.Run(t, helpers.RunOptions{Stdin: src})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "function")
}

func TestStdinUnicode(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const café = 1;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestStdinLongInput(t *testing.T) {
	// A larger input than the CLI's internal buffers.
	var b strings.Builder
	for i := 0; i < 5000; i++ {
		b.WriteString("const v")
		b.WriteString(itoa(i))
		b.WriteString(" = ")
		b.WriteString(itoa(i))
		b.WriteString(";\n")
	}
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte(b.String())})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "const")
}

func TestStdinInvalidUTF8(t *testing.T) {
	// Invalid UTF-8 bytes should produce an error, not a crash.
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte{0xff, 0xfe, 0xfd, 0x00},
	})
	if r.ExitCode == 0 {
		t.Fatalf("invalid UTF-8 should not silently succeed")
	}
	if r.ExitCode == 4 {
		t.Fatalf("invalid UTF-8 should be a processing error, not internal")
	}
}

func TestStdinBinaryContent(t *testing.T) {
	// A binary blob is not JavaScript. It should be reported as a
	// syntax error, not crash the binary.
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte{0x00, 0x01, 0x02, 0x03, 0xff},
	})
	if r.ExitCode == 0 {
		t.Fatalf("binary content should fail")
	}
}

func itoa(n int) string {
	// Local helper to avoid importing strconv in a test file that
	// does not otherwise need it.
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}