package api_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/api"
)

func TestAPIOptionsZeroIsPassthrough(t *testing.T) {
	src := "const x = 1;\n"
	r, err := api.Minify(src, api.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if r.HasErrors() {
		t.Fatalf("diagnostics: %v", r.Diagnostics)
	}
	if r.Code != src {
		t.Fatalf("zero Options should pass through: got %q, want %q", r.Code, src)
	}
}

func TestAPIOptionsWhitespaceOnly(t *testing.T) {
	r, err := api.Minify("function f() {\n    return 1;\n}\n", api.Options{
		MinifyWhitespace: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(r.Code, "\n    ") {
		t.Fatalf("whitespace not removed: %q", r.Code)
	}
}

func TestAPIOptionsSyntaxOnly(t *testing.T) {
	r, err := api.Minify("const x = 1 + 2 + 3;", api.Options{
		MinifySyntax: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.Code, "6") {
		t.Fatalf("syntax folding not applied: %q", r.Code)
	}
}

func TestAPIOptionsIdentifiersOnly(t *testing.T) {
	r, err := api.Minify("function f(longName) { return longName; }", api.Options{
		MinifyIdentifiers: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(r.Code, "longName") >= 2 {
		t.Fatalf("mangle not applied: %q", r.Code)
	}
}

func TestAPIOptionsFormatIIFE(t *testing.T) {
	r, err := api.Minify("const x = 1;", api.Options{
		Format: "iife",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.Code, "(") {
		t.Fatalf("iife format not applied: %q", r.Code)
	}
}

func TestAPIOptionsFormatESM(t *testing.T) {
	r, err := api.Minify("export const x = 1;", api.Options{
		Format: "esm",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.Code, "export") {
		t.Fatalf("esm format dropped export: %q", r.Code)
	}
}

func TestAPIOptionsFormatCJS(t *testing.T) {
	r, err := api.Minify("const x = 1;", api.Options{
		Format: "cjs",
	})
	if err != nil {
		t.Fatal(err)
	}
	// cjs may add "module.exports" for top-level declarations.
}

func TestAPIOptionsSourcemapInline(t *testing.T) {
	r, err := api.Minify("const x = 1;", api.Options{
		Sourcemap: "inline",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.Code, "sourceMappingURL=data:application/json") {
		t.Fatalf("inline map not appended: %q", r.Code)
	}
}

func TestAPIOptionsSourcemapExternal(t *testing.T) {
	r, err := api.Minify("const x = 1;", api.Options{
		Sourcemap: "external",
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.Map == "" {
		t.Fatal("external map not populated in Result.Map")
	}
}

func TestAPIOptionsLegalCommentsNone(t *testing.T) {
	r, err := api.Minify("/*! (c) */ const x = 1;", api.Options{
		LegalComments: "none",
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(r.Code, "(c)") {
		t.Fatalf("legal comment not removed: %q", r.Code)
	}
}

func TestAPIOptionsBannerAndFooter(t *testing.T) {
	r, err := api.Minify("const x = 1;", api.Options{
		Banner: "/*b*/",
		Footer: "/*f*/",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(r.Code, "/*b*/") {
		t.Fatalf("banner missing: %q", r.Code)
	}
	if !strings.HasSuffix(r.Code, "/*f*/") {
		t.Fatalf("footer missing: %q", r.Code)
	}
}

func TestAPIOptionsSourceName(t *testing.T) {
	// SourceName only affects diagnostic locations, not output.
	r, err := api.Minify("const = 1;", api.Options{
		SourceName: "myfile.js",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !r.HasErrors() {
		t.Fatal("expected error")
	}
	if len(r.Diagnostics) == 0 {
		t.Fatal("no diagnostics")
	}
	if r.Diagnostics[0].File != "myfile.js" {
		t.Fatalf("SourceName not used: %q", r.Diagnostics[0].File)
	}
}