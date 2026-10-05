package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestFileToFile(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js",
		"function add(a, b) {\n    return a + b;\n}\n")
	out := filepath.Join(dir, "out.js")

	r := helpers.Run(t, helpers.RunOptions{Args: []string{in, "-o", out}})
	helpers.AssertExitCode(t, r, 0)
	got, _ := os.ReadFile(out)
	helpers.AssertNoNewlines(t, string(got))
	helpers.AssertContains(t, string(got), "function")
}

func TestFileToStdout(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	r := helpers.Run(t, helpers.RunOptions{Args: []string{in}})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "const x")
	helpers.AssertNoStderr(t, r)
}

func TestMissingFileExitsThree(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Args: []string{"/does/not/exist.js"}})
	helpers.AssertExitCode(t, r, 3)
	helpers.AssertContains(t, string(r.Stderr), "no such file")
}

func TestDirectoryAsInputExitsOne(t *testing.T) {
	dir := helpers.TmpDir(t)
	r := helpers.Run(t, helpers.RunOptions{Args: []string{dir}})
	if r.ExitCode == 0 {
		t.Fatalf("directory as input should fail")
	}
}

func TestOutputToUnwritablePathExitsOne(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	blocker := helpers.WriteString(t, dir, "blocker", "x")
	out := filepath.Join(blocker, "out.js")

	r := helpers.Run(t, helpers.RunOptions{Args: []string{in, "-o", out}})
	if r.ExitCode == 0 {
		t.Fatalf("writing under a file should fail")
	}
}

func TestOutputOverwritesExistingFile(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := helpers.WriteString(t, dir, "out.js", "// OLD\n")

	r := helpers.Run(t, helpers.RunOptions{Args: []string{in, "-o", out}})
	helpers.AssertExitCode(t, r, 0)
	got, _ := os.ReadFile(out)
	helpers.AssertNotContains(t, string(got), "OLD")
}

func TestMultipleInputFilesExitsTwo(t *testing.T) {
	dir := helpers.TmpDir(t)
	a := helpers.WriteString(t, dir, "a.js", "const a = 1;\n")
	b := helpers.WriteString(t, dir, "b.js", "const b = 2;\n")
	r := helpers.Run(t, helpers.RunOptions{Args: []string{a, b}})
	helpers.AssertExitCode(t, r, 2)
}

func TestInputFileWithBOM(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "\xef\xbb\xbfconst x = 1;\n")
	r := helpers.Run(t, helpers.RunOptions{Args: []string{in}})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "const x")
}

func TestInputFileWithCRLF(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\r\nconst y = 2;\r\n")
	r := helpers.Run(t, helpers.RunOptions{Args: []string{in}})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "\r")
}

func TestOutputCreatesFile(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := filepath.Join(dir, "never-existed.js")
	r := helpers.Run(t, helpers.RunOptions{Args: []string{in, "-o", out}})
	helpers.AssertExitCode(t, r, 0)
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("output not created: %v", err)
	}
}

func TestInputAndOutputCanBeSameFile(t *testing.T) {
	// minifyjs in-place: read, minify, write back.
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js",
		"function add(a, b) {\n    return a + b;\n}\n")
	r := helpers.Run(t, helpers.RunOptions{Args: []string{in, "-o", in}})
	helpers.AssertExitCode(t, r, 0)
	got, _ := os.ReadFile(in)
	helpers.AssertNoNewlines(t, string(got))
}

func TestRelativeInputPath(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestInputWithSpacesInPath(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "with spaces.js", "const x = 1;\n")
	r := helpers.Run(t, helpers.RunOptions{Args: []string{in}})
	helpers.AssertExitCode(t, r, 0)
}

func TestInputWithUnicodePath(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "café.js", "const x = 1;\n")
	r := helpers.Run(t, helpers.RunOptions{Args: []string{in}})
	helpers.AssertExitCode(t, r, 0)
}

func TestInputWithVeryLongLine(t *testing.T) {
	dir := helpers.TmpDir(t)
	var b strings.Builder
	b.WriteString("const x = ")
	for i := 0; i < 50000; i++ {
		b.WriteString("1+")
	}
	b.WriteString("1;\n")
	in := helpers.WriteString(t, dir, "in.js", b.String())
	r := helpers.Run(t, helpers.RunOptions{Args: []string{in}})
	helpers.AssertExitCode(t, r, 0)
}

func TestEmptyFile(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "empty.js", "")
	r := helpers.Run(t, helpers.RunOptions{Args: []string{in}})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoStdout(t, r)
}