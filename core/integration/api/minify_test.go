// Package api_test exercises core/api end-to-end. These are still
// integration tests because api.Minify calls into engine, which calls
// into esbuild. They do not compile or spawn the CLI; they import the
// API directly and call it in-process.
package api_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/api"
)

func TestAPIMinifySimple(t *testing.T) {
	r, err := api.Minify("function add(a, b) { return a + b; }", api.Options{
		MinifyWhitespace: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.HasErrors() {
		t.Fatalf("diagnostics: %v", r.Diagnostics)
	}
	if strings.Contains(r.Code, "\n") {
		t.Fatalf("output has newline: %q", r.Code)
	}
}

func TestAPIMinifyEmpty(t *testing.T) {
	r, err := api.Minify("", api.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if r.Code != "" {
		t.Fatalf("empty input produced %q", r.Code)
	}
}

func TestAPIMinifySyntaxError(t *testing.T) {
	r, err := api.Minify("function () { } }", api.Options{MinifyWhitespace: true})
	if err != nil {
		t.Fatal(err)
	}
	if !r.HasErrors() {
		t.Fatal("expected error diagnostic")
	}
}

func TestAPIMinifyBytesSaved(t *testing.T) {
	r, err := api.Minify("function add(a, b) {\n    return a + b;\n}\n", api.Options{
		MinifyWhitespace: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.BytesSaved() <= 0 {
		t.Fatalf("BytesSaved = %d", r.BytesSaved())
	}
}

func TestAPIMinifyCompress(t *testing.T) {
	r, err := api.Minify("const x = 1 + 2 + 3;", api.Options{
		MinifyWhitespace: true,
		MinifySyntax:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.Code, "6") {
		t.Fatalf("constant folding not applied: %q", r.Code)
	}
}

func TestAPIMinifyMangle(t *testing.T) {
	r, err := api.Minify("function f(longParameter) { return longParameter; }", api.Options{
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(r.Code, "longParameter") >= 2 {
		t.Fatalf("mangling not applied: %q", r.Code)
	}
}

func TestAPIMinifyTarget(t *testing.T) {
	r, err := api.Minify("const f = (a) => a + 1;", api.Options{
		MinifyWhitespace: true,
		Target:           "es5",
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(r.Code, "=>") {
		t.Fatalf("es5 target not applied: %q", r.Code)
	}
}

func TestAPIMinifySourcemap(t *testing.T) {
	r, err := api.Minify("const x = 1;", api.Options{
		MinifyWhitespace: true,
		Sourcemap:        "external",
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.Map == "" {
		t.Fatal("external sourcemap should populate Map")
	}
}

func TestAPIMinifyBanner(t *testing.T) {
	r, err := api.Minify("const x = 1;", api.Options{
		MinifyWhitespace: true,
		Banner:           "/* (c) */",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(r.Code, "/* (c) */") {
		t.Fatalf("banner missing: %q", r.Code)
	}
}

func TestAPIMinifyInvalidTarget(t *testing.T) {
	_, err := api.Minify("const x = 1;", api.Options{Target: "ES2015"})
	if err == nil {
		t.Fatal("invalid target should error")
	}
}