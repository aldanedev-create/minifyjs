package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runCLI is a test helper that runs the CLI with fake streams and an
// empty environment. It returns exit code, stdout, stderr.
func runCLI(t *testing.T, args []string, stdin string) (int, string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	code := runWith(args, strings.NewReader(stdin), &out, &errb, func(string) string { return "" })
	return code, out.String(), errb.String()
}

func TestHelp(t *testing.T) {
	code, out, _ := runCLI(t, []string{"--help"}, "")
	if code != exitOK {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(out, "Usage:") {
		t.Fatalf("help missing Usage: %q", out)
	}
}

func TestVersion(t *testing.T) {
	code, out, _ := runCLI(t, []string{"--version"}, "")
	if code != exitOK {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(out, "minifyjs") || !strings.Contains(out, "esbuild") {
		t.Fatalf("version output = %q", out)
	}
}

func TestStdinToStdout(t *testing.T) {
	code, out, errS := runCLI(t, nil, "function add(a, b) {\n    return a + b;\n}\n")
	if code != exitOK {
		t.Fatalf("exit %d, stderr=%q", code, errS)
	}
	if strings.Contains(out, "\n") {
		t.Fatalf("output still contains newlines: %q", out)
	}
	if !strings.Contains(out, "function") {
		t.Fatalf("output = %q", out)
	}
}

func TestNoMinifyPassesThrough(t *testing.T) {
	src := "const x = 1;\n"
	code, out, errS := runCLI(t, []string{"--no-minify"}, src)
	if code != exitOK {
		t.Fatalf("exit %d, stderr=%q", code, errS)
	}
	if !strings.Contains(out, "const x") {
		t.Fatalf("output = %q", out)
	}
}

func TestUnknownFlag(t *testing.T) {
	code, _, errS := runCLI(t, []string{"--not-a-flag"}, "")
	if code != exitUsage {
		t.Fatalf("exit %d, want %d", code, exitUsage)
	}
	if !strings.Contains(errS, "unknown flag") {
		t.Fatalf("stderr = %q", errS)
	}
}

func TestMissingValue(t *testing.T) {
	code, _, errS := runCLI(t, []string{"--target"}, "")
	if code != exitUsage {
		t.Fatalf("exit %d, want %d", code, exitUsage)
	}
	if !strings.Contains(errS, "requires a value") {
		t.Fatalf("stderr = %q", errS)
	}
}

func TestFileToFile(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.js")
	out := filepath.Join(dir, "out.js")
	if err := os.WriteFile(in, []byte("function add(a, b) {\n    return a + b;\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, _, errS := runCLI(t, []string{in, "-o", out}, "")
	if code != exitOK {
		t.Fatalf("exit %d, stderr=%q", code, errS)
	}
	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(got), "\n") {
		t.Fatalf("output not minified: %q", got)
	}
}

func TestMissingInputFile(t *testing.T) {
	code, _, errS := runCLI(t, []string{"/does/not/exist.js"}, "")
	if code != exitNotFound {
		t.Fatalf("exit %d, want %d (stderr=%q)", code, exitNotFound, errS)
	}
}

func TestSyntaxErrorExitsNonZero(t *testing.T) {
	code, _, errS := runCLI(t, nil, "function () { } }")
	if code != exitError {
		t.Fatalf("exit %d, want %d (stderr=%q)", code, exitError, errS)
	}
	if !strings.Contains(errS, "error") {
		t.Fatalf("stderr = %q", errS)
	}
}

func TestTargetFlag(t *testing.T) {
	code, out, errS := runCLI(t, []string{"--target", "es5"}, "const f = (a, b) => a + b;")
	if code != exitOK {
		t.Fatalf("exit %d, stderr=%q", code, errS)
	}
	if strings.Contains(out, "=>") {
		t.Fatalf("es5 target should lower arrows: %q", out)
	}
}

func TestSourcemapInline(t *testing.T) {
	code, out, errS := runCLI(t, []string{"--sourcemap=inline"}, "const x = 1;")
	if code != exitOK {
		t.Fatalf("exit %d, stderr=%q", code, errS)
	}
	if !strings.Contains(out, "sourceMappingURL=data:application/json") {
		t.Fatalf("inline map missing: %q", out)
	}
}

func TestBannerAndFooter(t *testing.T) {
	code, out, errS := runCLI(t,
		[]string{"--banner", "/* B */", "--footer", "/* F */"},
		"const x = 1;")
	if code != exitOK {
		t.Fatalf("exit %d, stderr=%q", code, errS)
	}
	if !strings.HasPrefix(out, "/* B */") {
		t.Fatalf("banner missing: %q", out)
	}
	if !strings.HasSuffix(out, "/* F */") {
		t.Fatalf("footer missing: %q", out)
	}
}

func TestFlagAfterPositional(t *testing.T) {
	// `minifyjs app.js -o out.js` and `minifyjs -o out.js app.js`
	// must both work — the hand-written parser exists for this.
	dir := t.TempDir()
	in := filepath.Join(dir, "in.js")
	out := filepath.Join(dir, "out.js")
	_ = os.WriteFile(in, []byte("let x = 1;\n"), 0o644)

	code, _, errS := runCLI(t, []string{in, "-o", out}, "")
	if code != exitOK {
		t.Fatalf("file-then-flag failed: exit %d, stderr=%q", code, errS)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("output not written: %v", err)
	}
}

func TestDefinesParsed(t *testing.T) {
	f, err := parseFlags([]string{"--define", "DEBUG=false", "--define", "ENV=\"prod\""})
	if err != nil {
		t.Fatal(err)
	}
	if f.defines["DEBUG"] != "false" {
		t.Fatalf("DEBUG = %q", f.defines["DEBUG"])
	}
	if f.defines["ENV"] != `"prod"` {
		t.Fatalf("ENV = %q", f.defines["ENV"])
	}
}

func TestDropsParsed(t *testing.T) {
	f, err := parseFlags([]string{"--drop", "console", "--drop", "debugger"})
	if err != nil {
		t.Fatal(err)
	}
	if len(f.drops) != 2 || f.drops[0] != "console" || f.drops[1] != "debugger" {
		t.Fatalf("drops = %v", f.drops)
	}
}