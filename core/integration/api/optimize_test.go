package api_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/api"
)

func TestAPIOptimizeAllPasses(t *testing.T) {
	r, err := api.Optimize("const x = 1 + 2 + 3;", api.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if r.HasErrors() {
		t.Fatalf("diagnostics: %v", r.Diagnostics)
	}
	if !strings.Contains(r.Code, "6") {
		t.Fatalf("folding not applied: %q", r.Code)
	}
}

func TestAPIOptimizeMangles(t *testing.T) {
	r, err := api.Optimize("function f(longName) { return longName; }", api.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(r.Code, "longName") >= 2 {
		t.Fatalf("mangle not applied: %q", r.Code)
	}
}

func TestAPIOptimizeEmpty(t *testing.T) {
	r, err := api.Optimize("", api.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if r.Code != "" {
		t.Fatalf("empty output %q", r.Code)
	}
}

func TestAPIOptimizeSyntaxError(t *testing.T) {
	r, err := api.Optimize("function () { } }", api.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !r.HasErrors() {
		t.Fatal("expected error diagnostic")
	}
}

func TestAPIOptimizeWithTarget(t *testing.T) {
	r, err := api.Optimize("const f = (a) => a + 1;", api.Options{Target: "es5"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(r.Code, "=>") {
		t.Fatalf("target not applied: %q", r.Code)
	}
}

func TestAPIOptimizeWithBanner(t *testing.T) {
	r, err := api.Optimize("const x = 1;", api.Options{Banner: "/*b*/"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(r.Code, "/*b*/") {
		t.Fatalf("banner missing: %q", r.Code)
	}
}

func TestAPIOptimizeWithSourcemap(t *testing.T) {
	r, err := api.Optimize("const x = 1;", api.Options{Sourcemap: "external"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Map == "" {
		t.Fatal("map not produced")
	}
}

func TestAPIOptimizeRatio(t *testing.T) {
	src := "function add(a, b) {\n    return a + b;\n}\n"
	r, err := api.Optimize(src, api.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if r.MinifiedBytes >= r.OriginalBytes {
		t.Fatalf("output not smaller: %d -> %d", r.OriginalBytes, r.MinifiedBytes)
	}
}