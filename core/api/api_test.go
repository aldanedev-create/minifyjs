package api

import (
	"strings"
	"testing"
)

func TestMinifyBasic(t *testing.T) {
	r, err := Minify("function add(a, b) {\n    return a + b;\n}\n", Options{
		MinifyWhitespace: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.HasErrors() {
		t.Fatalf("diagnostics: %v", r.Diagnostics)
	}
	if strings.Contains(r.Code, "\n") {
		t.Fatalf("output still has newlines: %q", r.Code)
	}
	if r.BytesSaved() <= 0 {
		t.Fatalf("BytesSaved = %d", r.BytesSaved())
	}
}

func TestOptimizeEnablesAllPasses(t *testing.T) {
	r, err := Optimize("const x = 1 + 2 + 3;", Options{})
	if err != nil {
		t.Fatal(err)
	}
	if r.HasErrors() {
		t.Fatalf("diagnostics: %v", r.Diagnostics)
	}
	if !strings.Contains(r.Code, "6") {
		t.Fatalf("expected folding to 6: %q", r.Code)
	}
}

func TestMinifyReportsSyntaxError(t *testing.T) {
	r, err := Minify("function () { } }", Options{MinifyWhitespace: true})
	if err != nil {
		t.Fatal(err)
	}
	if !r.HasErrors() {
		t.Fatal("expected error diagnostic")
	}
}