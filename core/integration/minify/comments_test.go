package minify_test

import (
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestLineCommentRemoved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("// comment\nconst x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "comment")
}

func TestBlockCommentRemoved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("/* comment */ const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "comment")
}

func TestTrailingLineCommentRemoved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const x = 1; // trailing\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "trailing")
}

func TestMultipleCommentsRemoved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("// a\n// b\n/* c */\nconst x = 1; /* d */ // e\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "a")
	helpers.AssertNotContains(t, string(r.Stdout), "b")
	helpers.AssertNotContains(t, string(r.Stdout), "c")
	helpers.AssertNotContains(t, string(r.Stdout), "d")
	helpers.AssertNotContains(t, string(r.Stdout), "e")
}

func TestJSDocRemovedByDefault(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("/**\n * @param x\n */\nfunction f(x) { return x; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "@param")
}

func TestLegalCommentAtEof(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--legal-comments", "eof"},
		Stdin: []byte("/*! (c) 2026 */\nconst x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "(c) 2026")
}

func TestLegalCommentAtInline(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--legal-comments", "inline"},
		Stdin: []byte("/*! (c) 2026 */\nconst x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "(c) 2026")
}

func TestLegalCommentNone(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--legal-comments", "none"},
		Stdin: []byte("/*! (c) 2026 */\nconst x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "(c) 2026")
}

func TestCommentInsideStringPreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte(`const s = "// not a comment";` + "\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "// not a comment")
}

func TestCommentInsideTemplatePreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const s = `// not a comment`;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "// not a comment")
}

func TestCommentInsideRegexPreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const re = /\\/\\//;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestCommentBetweenCodeAndComment(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const a = 1; // one\nconst b = 2; // two\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "one")
	helpers.AssertNotContains(t, string(r.Stdout), "two")
}