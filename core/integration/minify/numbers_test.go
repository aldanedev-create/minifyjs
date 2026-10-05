package minify_test

import (
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestIntegerPreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = 42;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "42")
}

func TestFloatPreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = 3.14;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "3.14")
}

func TestLeadingDecimalPoint(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = .5;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestExponentNotation(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = 1e10;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestNegativeExponent(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = 1e-10;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestHexLiteral(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = 0xFF;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestOctalLiteral(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = 0o17;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBinaryLiteral(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = 0b101;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBigIntLiteral(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const x = 123n;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestNumericSeparators(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const x = 1_000_000;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestNumberMethodCall(t *testing.T) {
	// `1 .toString()` must not become `1.toString()`.
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const s = 1 .toString();\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestNumberWithDecimals(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = 0.5;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestInfinityAndNaN(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const a = Infinity; const b = NaN;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestConstantFoldingOfAddition(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--compress"},
		Stdin: []byte("const x = 1 + 2 + 3;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "6")
}