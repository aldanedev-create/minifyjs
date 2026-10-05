package minify_test

import (
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

// Idempotence: running minification on already-minified output must
// produce the same output. This is a strong invariant that catches
// many classes of bugs in the minifier itself.

func assertIdempotent(t *testing.T, args []string, src []byte) {
	t.Helper()
	r1 := helpers.Run(t, helpers.RunOptions{Args: args, Stdin: src})
	helpers.AssertExitCode(t, r1, 0)

	r2 := helpers.Run(t, helpers.RunOptions{Args: args, Stdin: r1.Stdout})
	helpers.AssertExitCode(t, r2, 0)

	if string(r1.Stdout) != string(r2.Stdout) {
		t.Fatalf("not idempotent:\nfirst:  %q\nsecond: %q", r1.Stdout, r2.Stdout)
	}
}

func TestIdempotentSimple(t *testing.T) {
	assertIdempotent(t, nil, []byte("function add(a, b) { return a + b; }"))
}

func TestIdempotentWithNewlines(t *testing.T) {
	assertIdempotent(t, nil, []byte("function f() {\n    return 1;\n}\n"))
}

func TestIdempotentWithComments(t *testing.T) {
	assertIdempotent(t, nil, []byte("// comment\nconst x = 1;\n"))
}

func TestIdempotentCompress(t *testing.T) {
	assertIdempotent(t, []string{"--compress"}, []byte("const x = 1 + 2 + 3;"))
}

func TestIdempotentMangle(t *testing.T) {
	assertIdempotent(t, []string{"--mangle"},
		[]byte("function f(longName) { return longName; }"))
}

func TestIdempotentCompressMangle(t *testing.T) {
	assertIdempotent(t, []string{"--compress", "--mangle"},
		[]byte("function f(x) { const y = 1 + 2; return x + y; }"))
}

func TestIdempotentModernSyntax(t *testing.T) {
	assertIdempotent(t, nil,
		[]byte("const f = (x) => x?.y ?? 0;"))
}

func TestIdempotentClasses(t *testing.T) {
	assertIdempotent(t, nil,
		[]byte("class A { constructor(x) { this.x = x; } }"))
}