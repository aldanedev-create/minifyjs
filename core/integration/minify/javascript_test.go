package minify_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestFixtureSimple(t *testing.T) {
	src := helpers.InputFixture(t, "simple.js")
	r := helpers.Run(t, helpers.RunOptions{Stdin: src})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "function")
}

func TestFixtureWhitespaceHeavy(t *testing.T) {
	src := helpers.InputFixture(t, "whitespace-heavy.js")
	r := helpers.Run(t, helpers.RunOptions{Stdin: src})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoNewlines(t, string(r.Stdout))
}

func TestFixtureModernSyntax(t *testing.T) {
	src := helpers.InputFixture(t, "modern-syntax.js")
	r := helpers.Run(t, helpers.RunOptions{Stdin: src})
	helpers.AssertExitCode(t, r, 0)
}

func TestFixtureUnicode(t *testing.T) {
	src := helpers.InputFixture(t, "unicode.js")
	r := helpers.Run(t, helpers.RunOptions{Stdin: src})
	helpers.AssertExitCode(t, r, 0)
}

func TestFixtureEmpty(t *testing.T) {
	src := helpers.InputFixture(t, "empty.js")
	r := helpers.Run(t, helpers.RunOptions{Stdin: src})
	helpers.AssertExitCode(t, r, 0)
}

func TestFixtureSyntaxErrorFails(t *testing.T) {
	src := helpers.InputFixture(t, "syntax-error.js")
	r := helpers.Run(t, helpers.RunOptions{Stdin: src})
	if r.ExitCode == 0 {
		t.Fatalf("syntax-error fixture should fail")
	}
}

func TestFunctionDeclarationShrinks(t *testing.T) {
	src := "function add(a, b) {\n    return a + b;\n}\n"
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte(src)})
	helpers.AssertExitCode(t, r, 0)
	if len(r.Stdout) >= len(src) {
		t.Fatalf("output not smaller: %d >= %d", len(r.Stdout), len(src))
	}
}

func TestMultipleStatements(t *testing.T) {
	src := "let a = 1;\nlet b = 2;\nlet c = a + b;\nconsole.log(c);\n"
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte(src)})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoNewlines(t, string(r.Stdout))
}

func TestCodeWithOnlyCommentsProducesEmpty(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("// just a comment\n/* and another */\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoStdout(t, r)
}

func TestCodeWithBOMIsHandled(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("\xef\xbb\xbfconst x = 1;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "\xef\xbb\xbf")
}