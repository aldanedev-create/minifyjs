package minify_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestMangleShortensParameter(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle"},
		Stdin: []byte("function f(longParameterName) { return longParameterName; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	occurrences := strings.Count(string(r.Stdout), "longParameterName")
	if occurrences >= 2 {
		t.Fatalf("parameter not mangled: %d occurrences in %q",
			occurrences, r.Stdout)
	}
}

func TestMangleShortensLocal(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle"},
		Stdin: []byte("function f() { const longLocalName = 1; return longLocalName; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	if strings.Count(string(r.Stdout), "longLocalName") >= 2 {
		t.Fatalf("local not mangled: %q", r.Stdout)
	}
}

func TestManglePreservesTopLevel(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle"},
		Stdin: []byte("function exportedName() { return 1; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "exportedName")
}

func TestManglePreservesPropertyNames(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle"},
		Stdin: []byte("const obj = { propertyName: 1 }; console.log(obj.propertyName);"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "propertyName")
}

func TestManglePreservesGlobals(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle"},
		Stdin: []byte("console.log(document.body);"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "console")
	helpers.AssertContains(t, string(r.Stdout), "document")
}

func TestNoMangleKeepsNames(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--no-mangle"},
		Stdin: []byte("function f(longName) { return longName; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "longName")
}

func TestMangleNestedFunctions(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle"},
		Stdin: []byte("function outer(outerArg) { return function inner(innerArg) { return outerArg + innerArg; }; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	if strings.Count(string(r.Stdout), "outerArg") >= 2 {
		t.Fatalf("outerArg not mangled: %q", r.Stdout)
	}
	if strings.Count(string(r.Stdout), "innerArg") >= 2 {
		t.Fatalf("innerArg not mangled: %q", r.Stdout)
	}
}

func TestMangleClosure(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle"},
		Stdin: []byte("function makeCounter(startValue) { let countValue = startValue; return () => countValue++; }"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestMangleClassMethods(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle"},
		Stdin: []byte("class A { methodName(x) { return x; } } new A().methodName(1);"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "methodName")
}

func TestMangleCatchBinding(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle"},
		Stdin: []byte("try {} catch (longErrorName) { console.log(longErrorName); }"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestMangleIdempotentOnMangledOutput(t *testing.T) {
	src := "function f(longName) { return longName; }"
	r1 := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle"},
		Stdin: []byte(src),
	})
	helpers.AssertExitCode(t, r1, 0)

	r2 := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle"},
		Stdin: r1.Stdout,
	})
	helpers.AssertExitCode(t, r2, 0)
	helpers.AssertEqual(t, string(r1.Stdout), string(r2.Stdout))
}

func TestMangleWithCompress(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--mangle", "--compress"},
		Stdin: []byte("function f(x) { const y = 1 + 2; return x + y; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "3")
}