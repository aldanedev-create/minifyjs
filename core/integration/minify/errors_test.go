package minify_test

import (
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestSyntaxErrorUnbalancedBrace(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("function f() { }")})
	// This is actually valid; only the closing brace of the
	// function body is present. Extra closing braces are the error.
	helpers.AssertExitCode(t, r, 0)
}

func TestSyntaxErrorExtraBrace(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("function f() { } }")})
	helpers.AssertExitCode(t, r, 1)
}

func TestSyntaxErrorUnclosedParen(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("f(1, 2")})
	helpers.AssertExitCode(t, r, 1)
}

func TestSyntaxErrorUnterminatedString(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte(`const s = "abc`)})
	helpers.AssertExitCode(t, r, 1)
}

func TestSyntaxErrorUnterminatedTemplate(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const s = `abc")})
	helpers.AssertExitCode(t, r, 1)
}

func TestSyntaxErrorUnterminatedComment(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("/* abc")})
	helpers.AssertExitCode(t, r, 1)
}

func TestSyntaxErrorInvalidRegex(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte(`const re = /[abc/;`)})
	helpers.AssertExitCode(t, r, 1)
}

func TestSyntaxErrorReservedWordAsVariable(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const if = 1;")})
	helpers.AssertExitCode(t, r, 1)
}

func TestSyntaxErrorDuplicateLet(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("let x = 1; let x = 2;")})
	helpers.AssertExitCode(t, r, 1)
}

func TestSyntaxErrorImportAtTop(t *testing.T) {
	// import outside a module is still a syntax error in a plain
	// script context; esbuild accepts it as an ES module.
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte(`import from "x";`)})
	helpers.AssertExitCode(t, r, 1)
}

func TestSyntaxErrorAwaitOutsideAsync(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("function f() { await x; }")})
	// Depending on esbuild's version, this may or may not be an
	// error. We only require that the CLI does not crash.
	if r.ExitCode == 4 {
		t.Fatalf("await outside async crashed: %s", r.Stderr)
	}
}

func TestSyntaxErrorReturnOutsideFunction(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("return 1;")})
	// esbuild tolerates top-level return in some contexts.
	if r.ExitCode == 4 {
		t.Fatalf("top-level return crashed: %s", r.Stderr)
	}
}

func TestSyntaxErrorEmptyObjectProperty(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const o = { : 1 };")})
	helpers.AssertExitCode(t, r, 1)
}

func TestSyntaxErrorMissingInitializer(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const x;")})
	helpers.AssertExitCode(t, r, 1)
}

func TestErrorOutputIncludesSourceSnippet(t *testing.T) {
	// esbuild reports errors with a snippet of the source. The CLI
	// should preserve it.
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const x = 1;\nconst = 2;\n")})
	if r.ExitCode == 0 {
		t.Fatalf("expected error")
	}
	if len(r.Stderr) < 10 {
		t.Fatalf("stderr too short to include a snippet: %q", r.Stderr)
	}
}