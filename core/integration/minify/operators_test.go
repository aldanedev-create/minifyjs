package minify_test

import (
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestArithmeticOperators(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = a + b - c * d / e % f;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestComparisonOperators(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = a < b && c > d || e <= f && g >= h;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestEqualityOperators(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = a === b && c !== d && e == f && g != h;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestLogicalOperators(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = a && b || c;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBitwiseOperators(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = a & b | c ^ d;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestShiftOperators(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = a << 2 >> 1 >>> 3;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestUnaryOperators(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = -a + +b - !c + ~d;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestTypeofOperator(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const t = typeof x;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestVoidOperator(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const v = void 0;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestDeleteOperator(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("delete obj.prop;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestInOperator(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const has = 'key' in obj;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestInstanceofOperator(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const is = obj instanceof Date;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestTernaryOperator(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const x = a ? b : c;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestAssignmentOperators(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("a += 1; b -= 2; c *= 3; d /= 4; e %= 5; f **= 2;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBitwiseAssignment(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("a &= b; c |= d; e ^= f; g <<= 1; h >>= 2; i >>>= 3;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestIncrementDecrement(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("let i = 0; i++; ++i; i--; --i;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestExponentiationPrecedence(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const x = 2 ** 3 ** 4;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestCommaOperator(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const x = (a(), b(), c());\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}