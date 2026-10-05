package engine

import (
	"strings"
	"testing"
)

func TestOptionsIsNoOp(t *testing.T) {
	if !(Options{}).IsNoOp() {
		t.Fatal("zero Options should be a no-op")
	}
	if (Options{MinifyWhitespace: true}).IsNoOp() {
		t.Fatal("MinifyWhitespace should not be a no-op")
	}
}

func TestOptionsValidateAcceptsKnownValues(t *testing.T) {
	cases := []Options{
		{},
		{Target: "es2015"},
		{Target: "esnext"},
		{Format: "esm"},
		{Format: "iife"},
		{Sourcemap: "inline"},
		{Sourcemap: "external"},
		{LegalComments: "none"},
		{Target: "es2020,chrome90"},
	}
	for i, c := range cases {
		if err := c.validate(); err != nil {
			t.Errorf("case %d: %v", i, err)
		}
	}
}

func TestOptionsValidateRejectsBadValues(t *testing.T) {
	cases := []Options{
		{Format: "not-a-format"},
		{Sourcemap: "not-a-mode"},
		{LegalComments: "not-a-mode"},
		{Target: "ES2015"}, // uppercase is not accepted
	}
	for i, c := range cases {
		if err := c.validate(); err == nil {
			t.Errorf("case %d should have failed: %+v", i, c)
		}
	}
}

func TestMapFormat(t *testing.T) {
	if mapFormat("esm").String() != "esm" {
		t.Fatal("esm mapping wrong")
	}
	if mapFormat("").String() != "default" {
		t.Fatal("default mapping wrong")
	}
}

func TestTransformEmptyIsNoError(t *testing.T) {
	res, err := Transform("", Options{})
	if err != nil {
		t.Fatal(err)
	}
	if res.HasErrors() {
		t.Fatalf("unexpected diagnostics: %v", res.Diagnostics)
	}
}

func TestTransformMinifyWhitespace(t *testing.T) {
	src := "function add(a, b) {\n    return a + b;\n}\n"
	res, err := Transform(src, Options{MinifyWhitespace: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.HasErrors() {
		t.Fatalf("unexpected diagnostics: %v", res.Diagnostics)
	}
	if strings.Contains(res.Code, "\n") {
		t.Fatalf("newlines survived whitespace minification: %q", res.Code)
	}
	if res.MinifiedBytes >= res.OriginalBytes {
		t.Fatalf("output not smaller: %d -> %d", res.OriginalBytes, res.MinifiedBytes)
	}
	if res.BytesSaved() <= 0 {
		t.Fatalf("BytesSaved = %d", res.BytesSaved())
	}
	if res.Ratio() <= 0 || res.Ratio() >= 1 {
		t.Fatalf("Ratio = %v, want 0<r<1", res.Ratio())
	}
}

func TestTransformInvalidTargetIsError(t *testing.T) {
	_, err := Transform("let x = 1;", Options{Target: "ES2015"})
	if err == nil {
		t.Fatal("expected error for uppercase target")
	}
}

func TestTransformReportsSyntaxError(t *testing.T) {
	// A lone closing brace is a syntax error in every JS grammar.
	res, err := Transform("function () { } }", Options{MinifyWhitespace: true})
	if err != nil {
		t.Fatal(err)
	}
	if !res.HasErrors() {
		t.Fatalf("expected errors, got %v", res.Diagnostics)
	}
}

func TestBuildRequiresEntryPoint(t *testing.T) {
	_, err := Build(BuildOptions{OutDir: "out"})
	if err == nil {
		t.Fatal("expected error for missing entry points")
	}
}

func TestBuildRequiresOutDirOrOutFile(t *testing.T) {
	_, err := Build(BuildOptions{EntryPoints: []string{"a.js"}})
	if err == nil {
		t.Fatal("expected error for missing out dir/file")
	}
}

func TestBuildRejectsOutDirAndOutFile(t *testing.T) {
	_, err := Build(BuildOptions{
		EntryPoints: []string{"a.js"},
		OutDir:      "out",
		OutFile:     "out.js",
	})
	if err == nil {
		t.Fatal("expected error for both out dir and file")
	}
}

func TestBuildSplittingRequiresESMAndBundle(t *testing.T) {
	_, err := Build(BuildOptions{
		EntryPoints: []string{"a.js"},
		OutDir:      "out",
		Splitting:   true,
		// Bundle false, Format empty
	})
	if err == nil {
		t.Fatal("expected error for splitting without esm+bundle")
	}
}