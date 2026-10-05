package minify_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestLeadingWhitespaceRemoved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("    const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	if strings.HasPrefix(string(r.Stdout), " ") {
		t.Fatalf("leading whitespace not removed: %q", r.Stdout)
	}
}

func TestTrailingWhitespaceRemoved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const x = 1;    "),
	})
	helpers.AssertExitCode(t, r, 0)
	if strings.HasSuffix(string(r.Stdout), " ") {
		t.Fatalf("trailing whitespace not removed: %q", r.Stdout)
	}
}

func TestNewlinesRemoved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const a = 1;\nconst b = 2;\nconst c = 3;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoNewlines(t, string(r.Stdout))
}

func TestTabsRemoved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("function f() {\n\treturn 1;\n}"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "\t")
}

func TestCRLFNormalized(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const a = 1;\r\nconst b = 2;\r\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "\r")
}

func TestMultipleBlankLinesCollapsed(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const a = 1;\n\n\n\nconst b = 2;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoNewlines(t, string(r.Stdout))
}

func TestSpaceAroundOperatorsRemoved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const x = 1 + 2 * 3 - 4 / 5;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), " + ")
	helpers.AssertNotContains(t, string(r.Stdout), " * ")
}

func TestSpaceAfterCommaRemoved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("f(a, b, c);"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), ", ")
}

func TestSpaceAroundBracesRemoved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("if (x) { y(); }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "{ ")
	helpers.AssertNotContains(t, string(r.Stdout), " }")
}

func TestSpaceNeededBetweenIdentifiers(t *testing.T) {
	// `return x` must keep the space: `returnx` would be one token.
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("function f() { return x; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "return")
}

func TestSpaceNeededBetweenKeywordsAndIdentifiers(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("var x = 1; let y = 2; const z = 3;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "var x")
	helpers.AssertContains(t, string(r.Stdout), "let y")
	helpers.AssertContains(t, string(r.Stdout), "const z")
}

func TestIdempotentWhitespace(t *testing.T) {
	// Running the CLI on its own output produces the same output.
	src := "function f() {\n    return 1;\n}\n"
	r1 := helpers.Run(t, helpers.RunOptions{Stdin: []byte(src)})
	helpers.AssertExitCode(t, r1, 0)

	r2 := helpers.Run(t, helpers.RunOptions{Stdin: r1.Stdout})
	helpers.AssertExitCode(t, r2, 0)

	if string(r1.Stdout) != string(r2.Stdout) {
		t.Fatalf("not idempotent:\nfirst:  %q\nsecond: %q", r1.Stdout, r2.Stdout)
	}
}

func TestShebangPreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("#!/usr/bin/env node\nconst x = 1;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "#!/usr/bin/env node")
}