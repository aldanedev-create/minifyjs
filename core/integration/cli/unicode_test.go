package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestUnicodeIdentifiers(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte("const café = 1;\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "café")
}

func TestUnicodeStrings(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte(`const s = "日本語";` + "\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "日本語")
}

func TestUnicodeEscapes(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte(`const s = "\u00e9";`),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestEmojiInString(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte(`const s = "🎉";` + "\n"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "🎉")
}

func TestEmojiInComment(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("// 🎉 comment\nconst x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "🎉")
}

func TestUnicodeEscapePreserved(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-minify"},
		Stdin: []byte(`const s = "\u{1F389}";` + "\n"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestUnicodeInFilename(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "café-日本語.js", "const x = 1;\n")
	r := helpers.Run(t, helpers.RunOptions{Args: []string{in}})
	helpers.AssertExitCode(t, r, 0)
}

func TestUnicodeInOutputPath(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := filepath.Join(dir, "café.js")

	r := helpers.Run(t, helpers.RunOptions{Args: []string{in, "-o", out}})
	helpers.AssertExitCode(t, r, 0)
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("unicode output path: %v", err)
	}
}