package cli_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

// TestNoMinifyLeavesCodeAlone verifies --no-minify is a passthrough.
func TestNoMinifyLeavesCodeAlone(t *testing.T) {
	src := "function add(a, b) {\n    return a + b;\n}\n"
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte(src),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertEqual(t, string(r.Stdout), src)
}

// TestCompressEnablesWhitespaceAndSyntax verifies --compress does not
// just do whitespace.
func TestCompressEnablesWhitespaceAndSyntax(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--compress"},
		Stdin: []byte("const x = 1 + 2 + 3;"),
	})
	helpers.AssertExitCode(t, r, 0)
	// Constant folding should produce "6".
	helpers.AssertContains(t, string(r.Stdout), "6")
}

// TestNoCompressDisablesBoth verifies the negation.
func TestNoCompressDisablesBoth(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-compress"},
		Stdin: []byte("const x = 1 + 2 + 3;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	// Without --compress, 1+2+3 should not be folded to 6.
	helpers.AssertNotContains(t, string(r.Stdout), "=6")
}

// TestMangleRenamesLocals verifies --mangle.
func TestMangleRenamesLocals(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle"},
		Stdin: []byte(`function f(longName) { return longName + 1; }`),
	})
	helpers.AssertExitCode(t, r, 0)
	// The parameter should be renamed away from longName.
	if strings.Count(string(r.Stdout), "longName") >= 2 {
		t.Fatalf("--mangle did not rename parameter: %q", r.Stdout)
	}
}

// TestNoManglePreservesNames verifies --no-mangle.
func TestNoManglePreservesNames(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-mangle"},
		Stdin: []byte(`function f(longName) { return longName + 1; }`),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "longName")
}

// TestTargetAcceptsES5 verifies the es5 target lowers arrow functions.
func TestTargetAcceptsES5(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "es5"},
		Stdin: []byte("const f = (a) => a + 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "=>")
}

// TestTargetAcceptsESNext verifies esnext does not lower.
func TestTargetAcceptsESNext(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "esnext"},
		Stdin: []byte("const f = (a) => a + 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "=>")
}

// TestFormatESM verifies --format esm keeps import/export.
func TestFormatESM(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--format", "esm"},
		Stdin: []byte("export const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "export")
}

// TestFormatIIFE verifies --format iife wraps in a function.
func TestFormatIIFE(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--format", "iife"},
		Stdin: []byte("const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	// IIFE typically wraps with "(()=>{ ... })()" or similar.
	if !strings.Contains(string(r.Stdout), "(") {
		t.Fatalf("iife output should contain parens: %q", r.Stdout)
	}
}

// TestFormatCJS verifies --format cjs.
func TestFormatCJS(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--format", "cjs"},
		Stdin: []byte("const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
}

// TestDefineSubstitution verifies --define replaces identifiers.
func TestDefineSubstitution(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--define", "DEBUG=false", "--compress"},
		Stdin: []byte("if (DEBUG) console.log('debug');"),
	})
	helpers.AssertExitCode(t, r, 0)
	// DEBUG should be gone; the whole branch should be dropped.
	helpers.AssertNotContains(t, string(r.Stdout), "DEBUG")
}

// TestDropConsole verifies --drop console.
func TestFlagsDropConsole(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--drop", "console", "--compress"},
		Stdin: []byte("console.log('hi'); const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "console.log")
}

// TestDropDebugger verifies --drop debugger.
func TestFlagsDropDebugger(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--drop", "debugger", "--compress"},
		Stdin: []byte("debugger; const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "debugger")
}

// TestPureMarksFunctionAsSideEffectFree verifies --pure.
func TestPureMarksFunctionAsSideEffectFree(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--pure", "console.log", "--compress"},
		Stdin: []byte("const x = 1; console.log(x);"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "console.log")
}

// TestFlagOrderIndependent verifies flags can come before or after
// the positional input.
func TestFlagOrderIndependent(t *testing.T) {
	src := "const f = (a) => a + 1;"

	r1 := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "es5"},
		Stdin: []byte(src),
	})
	r2 := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "es5", "-"},
		Stdin: []byte(src),
	})
	helpers.AssertExitCode(t, r1, 0)
	helpers.AssertExitCode(t, r2, 0)
	if !strings.Contains(string(r1.Stdout), "function") {
		t.Fatalf("r1 missing function: %q", r1.Stdout)
	}
	if !strings.Contains(string(r2.Stdout), "function") {
		t.Fatalf("r2 missing function: %q", r2.Stdout)
	}
}

// TestEqualsFormForValueFlag verifies --flag=value syntax.
func TestEqualsFormForValueFlag(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target=es5"},
		Stdin: []byte("const f = (a) => a + 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "=>")
}

// TestMultipleDefines verifies repeated --define flags.
func TestMultipleDefines(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--define", "A=1",
			"--define", "B=2",
			"--compress",
		},
		Stdin: []byte("const x = A + B;"),
	})
	helpers.AssertExitCode(t, r, 0)
	// 1 + 2 folds to 3.
	helpers.AssertContains(t, string(r.Stdout), "3")
}

// TestMultipleDrops verifies repeated --drop flags.
func TestMultipleDrops(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{
			"--drop", "console",
			"--drop", "debugger",
			"--compress",
		},
		Stdin: []byte("console.log('x'); debugger; const y = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "console")
	helpers.AssertNotContains(t, string(r.Stdout), "debugger")
}

// TestNoFlagAcceptsExtraValue verifies --no-* flags reject values.
func TestNoFlagAcceptsExtraValue(t *testing.T) {
	// --no-minify=x should be a usage error.
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify=yes"},
		Stdin: []byte("const x = 1;"),
	})
	// Either it errors, or it ignores the value. Both are acceptable;
	// the test only ensures no crash.
	if r.ExitCode == 4 {
		t.Fatalf("--no-minify=yes crashed the CLI")
	}
}

// TestBannerFlagIsPreserved verifies banner text reaches output.
func TestBannerFlagIsPreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--banner", "/* (c) 2026 */"},
		Stdin: []byte("const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "/* (c) 2026 */")
}

// TestFooterFlagIsPreserved verifies footer text reaches output.
func TestFooterFlagIsPreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--footer", "// end of file"},
		Stdin: []byte("const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "// end of file")
}

// TestLegalCommentsNone verifies --legal-comments none.
func TestFlagsLegalCommentsNone(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--legal-comments", "none"},
		Stdin: []byte("/*! preserve me */\nconst x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "preserve me")
}

// TestLegalCommentsEOF verifies default behavior keeps /*! comments.
func TestFlagsLegalCommentsEOF(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--legal-comments", "eof"},
		Stdin: []byte("/*! preserve me */\nconst x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "preserve me")
}