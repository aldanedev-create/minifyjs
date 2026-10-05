package cli_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestBannerPrepended(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--banner", "/* (c) 2026 */"},
		Stdin: []byte("const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	if !strings.HasPrefix(string(r.Stdout), "/* (c) 2026 */") {
		t.Fatalf("banner not at start: %q", r.Stdout)
	}
}

func TestFooterAppended(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--footer", "// end"},
		Stdin: []byte("const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	if !strings.HasSuffix(string(r.Stdout), "// end") {
		t.Fatalf("footer not at end: %q", r.Stdout)
	}
}

func TestBannerAndFooterTogether(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--banner", "/*b*/", "--footer", "/*f*/"},
		Stdin: []byte("const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "/*b*/")
	helpers.AssertContains(t, string(r.Stdout), "/*f*/")
}

func TestBannerEmptyStringAllowed(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--banner", ""},
		Stdin: []byte("const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestBannerFooterLegalCommentsNone(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--legal-comments", "none"},
		Stdin: []byte("/*! important */\nconst x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "important")
}

func TestLegalCommentsInline(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--legal-comments", "inline"},
		Stdin: []byte("/*! important */\nconst x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "important")
}

func TestBannerFooterLegalCommentsEOF(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--legal-comments", "eof"},
		Stdin: []byte("const x = 1;\n/*! tail comment */"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "tail comment")
}

func TestLegalCommentsInvalid(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--legal-comments", "bogus"},
		Stdin: []byte("const x = 1;"),
	})
	if r.ExitCode == 0 {
		t.Fatalf("invalid legal comments mode should fail")
	}
}

func TestRegularCommentsDroppedByDefault(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("// regular\nconst x = 1; // trailing\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "regular")
	helpers.AssertNotContains(t, string(r.Stdout), "trailing")
}

func TestJSDocCommentDroppedByDefault(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("/** JSDoc */\nfunction f() {}"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "JSDoc")
}