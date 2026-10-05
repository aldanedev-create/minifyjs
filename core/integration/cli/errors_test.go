package cli_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestSyntaxErrorGoesToStderr(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("function () { } }")})
	helpers.AssertExitCode(t, r, 1)
	helpers.AssertNoStdout(t, r)
	helpers.AssertContains(t, string(r.Stderr), "minifyjs:")
}

func TestUnknownFlagExitCodeIsTwo(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"--bogus-flag"}})
	helpers.AssertExitCode(t, r, 2)
}

func TestMissingValueExitCodeIsTwo(t *testing.T) {
	flags := []string{"-o", "--output", "--target", "--format", "--banner", "--footer", "--legal-comments", "--define", "--drop", "--pure", "--outdir", "--platform", "--cache-dir", "--config"}
	for _, flag := range flags {
		t.Run(flag, func(t *testing.T) {
			r := helpers.Run(t, helpers.RunOptions{Args: []string{flag}})
			helpers.AssertExitCode(t, r, 2)
		})
	}
}

func TestErrorIncludesFilename(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "bad.js", "const x = 1;\nconst = 2;\n")
	r := helpers.Run(t, helpers.RunOptions{Args: []string{in}})
	if r.ExitCode == 0 {
		t.Fatalf("expected error")
	}
	helpers.AssertContains(t, string(r.Stderr), "bad.js")
}

func TestErrorIncludesLineAndColumn(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const x = 1;\nconst = 2;\n")})
	if r.ExitCode == 0 {
		t.Fatalf("expected error")
	}
	// esbuild's message includes the location as "stdin:2:6" or
	// similar. We check for the line number.
	helpers.AssertContains(t, string(r.Stderr), "2")
}

func TestUnterminatedStringError(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte(`const s = "unterminated`)})
	if r.ExitCode == 0 {
		t.Fatalf("expected error for unterminated string")
	}
}

func TestUnterminatedTemplateError(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const s = `unterminated")})
	if r.ExitCode == 0 {
		t.Fatalf("expected error for unterminated template")
	}
}

func TestUnterminatedCommentError(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const x = 1;\n/* unterminated")})
	if r.ExitCode == 0 {
		t.Fatalf("expected error for unterminated comment")
	}
}

func TestInvalidRegexError(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte(`const re = /[unclosed/;`)})
	if r.ExitCode == 0 {
		t.Fatalf("expected error for invalid regex")
	}
}

func TestReservedWordAsIdentifierError(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const const = 1;")})
	if r.ExitCode == 0 {
		t.Fatalf("expected error for reserved word")
	}
}

func TestDuplicateDeclarationError(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("let x = 1; let x = 2;")})
	if r.ExitCode == 0 {
		t.Fatalf("expected error for duplicate declaration")
	}
}

func TestInvalidImportSyntax(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("import { from \"x\";")})
	if r.ExitCode == 0 {
		t.Fatalf("expected error for invalid import")
	}
}

func TestInvalidJSONInConfig(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", "{not valid json}")
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	if r.ExitCode == 0 {
		t.Fatalf("invalid config should fail")
	}
	if !strings.Contains(strings.ToLower(string(r.Stderr)), "config") {
		t.Fatalf("stderr should mention config: %s", r.Stderr)
	}
}