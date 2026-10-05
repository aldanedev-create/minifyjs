package minify_test

import (
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestSingleQuotesPreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const s = 'hello';\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestDoubleQuotesPreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte(`const s = "hello";` + "\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestEscapeSequencePreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte(`const s = "line\nbreak";` + "\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), `\n`)
}

func TestEscapeQuoteInsideString(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte(`const s = "she said \"hi\"";` + "\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), `\"`)
}

func TestBackslashInString(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte(`const s = "a\\b";` + "\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestEmptyString(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte(`const s = "";` + "\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), `""`)
}

func TestTemplateLiteralPreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const s = `hello`;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "`")
}

func TestTemplateWithInterpolation(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const s = `hello ${name}!`;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "${name}")
}

func TestTemplateNested(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const s = `a ${`b ${x}`} c`;"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestRegexLiteralPreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte(`const re = /ab+c/g;` + "\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "/ab+c/g")
}

func TestDivisionNotRegex(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = a / b / c;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestUnicodeEscapeInString(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte(`const s = "\u00e9";` + "\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), `\u00e9`)
}