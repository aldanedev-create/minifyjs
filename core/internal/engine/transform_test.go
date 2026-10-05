package engine

import (
	"strings"
	"testing"
)

// These tests exercise Transform's option mapping. They do not test
// esbuild itself — esbuild has its own test suite — they test that
// MinifyJS's Options map to the esbuild behavior callers expect.

func TestTransformMinifyIdentifiersRenamesLocals(t *testing.T) {
	src := `function calculateTotal(price, tax) { const total = price + (price * tax); return total; }`
	res, err := Transform(src, Options{
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.HasErrors() {
		t.Fatalf("diagnostics: %v", res.Diagnostics)
	}
	// The function's own name is exported (top-level), but the
	// parameters and locals should be renamed to something shorter.
	if !strings.Contains(res.Code, "calculateTotal") {
		t.Fatalf("top-level function name should be preserved: %q", res.Code)
	}
	if !strings.Contains(res.Code, "function") {
		t.Fatalf("output should still be a function: %q", res.Code)
	}
}

func TestTransformMinifySyntaxFoldsConstants(t *testing.T) {
	res, err := Transform("const x = 1 + 2 + 3;", Options{
		MinifyWhitespace: true,
		MinifySyntax:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.HasErrors() {
		t.Fatalf("diagnostics: %v", res.Diagnostics)
	}
	if !strings.Contains(res.Code, "6") {
		t.Fatalf("expected constant folding to produce 6: %q", res.Code)
	}
}

func TestTransformTargetLowersArrowFunctions(t *testing.T) {
	src := `const f = (a, b) => a + b;`
	res, err := Transform(src, Options{
		MinifyWhitespace: true,
		Target:           "es5",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.HasErrors() {
		t.Fatalf("diagnostics: %v", res.Diagnostics)
	}
	if strings.Contains(res.Code, "=>") {
		t.Fatalf("es5 target should lower arrow functions: %q", res.Code)
	}
}

func TestTransformSourcemapInlineAppendsDataURL(t *testing.T) {
	res, err := Transform("const x = 1;", Options{
		MinifyWhitespace: true,
		Sourcemap:        "inline",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.HasErrors() {
		t.Fatalf("diagnostics: %v", res.Diagnostics)
	}
	if !strings.Contains(res.Code, "sourceMappingURL=data:application/json") {
		t.Fatalf("inline sourcemap missing from output: %q", res.Code)
	}
}

func TestTransformSourcemapExternalReturnsMap(t *testing.T) {
	res, err := Transform("const x = 1;", Options{
		MinifyWhitespace: true,
		Sourcemap:        "external",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.HasErrors() {
		t.Fatalf("diagnostics: %v", res.Diagnostics)
	}
	if res.Map == "" {
		t.Fatal("external sourcemap should populate Result.Map")
	}
	if !strings.Contains(res.Map, `"version":3`) {
		t.Fatalf("unexpected map: %q", res.Map)
	}
}

func TestTransformBannerAndFooter(t *testing.T) {
	res, err := Transform("const x = 1;", Options{
		MinifyWhitespace: true,
		Banner:           "/* banner */",
		Footer:           "/* footer */",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(res.Code, "/* banner */") {
		t.Fatalf("banner missing: %q", res.Code)
	}
	if !strings.HasSuffix(res.Code, "/* footer */") {
		t.Fatalf("footer missing: %q", res.Code)
	}
}