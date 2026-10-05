package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultsEnablesMinify(t *testing.T) {
	d := Defaults()
	if !d.IsMinifying() {
		t.Fatal("defaults should enable minification")
	}
	if d.Minify.Whitespace != true || d.Minify.Identifiers != true || d.Minify.Syntax != true {
		t.Fatalf("defaults should enable all minify passes: %+v", d.Minify)
	}
	if d.Target != "" || d.Format != "" || d.Sourcemap != "" {
		t.Fatalf("defaults should not set target/format/sourcemap: %+v", d)
	}
}

func TestParseEmptyConfig(t *testing.T) {
	c, err := Parse([]byte(`{}`), ".")
	if err != nil {
		t.Fatal(err)
	}
	// Zero value: nothing set, caller is expected to start from
	// Defaults and overlay.
	if c.IsMinifying() {
		t.Fatal("empty config should not enable minify")
	}
}

func TestParseMinifyBlock(t *testing.T) {
	c, err := Parse([]byte(`{"minify": {"whitespace": true}}`), ".")
	if err != nil {
		t.Fatal(err)
	}
	if !c.Minify.Whitespace {
		t.Fatal("whitespace should be true")
	}
	if c.Minify.Identifiers || c.Minify.Syntax {
		t.Fatal("identifiers and syntax should be false")
	}
}

func TestParseResolvesPaths(t *testing.T) {
	c, err := Parse([]byte(`{"inputs": ["a.js"], "output": "dist/b.js"}`), "/repo")
	if err != nil {
		t.Fatal(err)
	}
	if c.Inputs[0] != filepath.Join("/repo", "a.js") {
		t.Fatalf("input not resolved: %q", c.Inputs[0])
	}
	if c.Output != filepath.Join("/repo", "dist", "b.js") {
		t.Fatalf("output not resolved: %q", c.Output)
	}
}

func TestParseAbsolutePathUnchanged(t *testing.T) {
	c, err := Parse([]byte(`{"output": "/tmp/x.js"}`), "/repo")
	if err != nil {
		t.Fatal(err)
	}
	if c.Output != "/tmp/x.js" {
		t.Fatalf("absolute path changed: %q", c.Output)
	}
}

func TestParseRejectsUnknownJSON(t *testing.T) {
	_, err := Parse([]byte(`not json`), ".")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestValidateAcceptsKnownValues(t *testing.T) {
	cases := []Config{
		{},
		{Target: "es2020"},
		{Target: "es2020,chrome90"},
		{Format: "esm"},
		{Sourcemap: "inline"},
		{LegalComments: "none"},
		{Quiet: true},
		{Verbose: true},
	}
	for i, c := range cases {
		if err := c.Validate(); err != nil {
			t.Errorf("case %d: %v", i, err)
		}
	}
}

func TestValidateRejectsBadValues(t *testing.T) {
	cases := []Config{
		{Format: "not-a-format"},
		{Sourcemap: "not-a-mode"},
		{LegalComments: "not-a-mode"},
		{Target: "ES2015"},
		{Quiet: true, Verbose: true},
		{Bundle: true}, // no entry points
		{Bundle: true, BundleOptions: BundleOptions{Splitting: true}}, // splitting without esm
	}
	for i, c := range cases {
		if err := c.Validate(); err == nil {
			t.Errorf("case %d should have failed: %+v", i, c)
		}
	}
}

func TestFindWalksUp(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(root, "minifyjs.config.json")
	if err := os.WriteFile(cfgPath, []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	found, ok, err := Find(sub)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected to find config")
	}
	if found != cfgPath {
		t.Fatalf("found %q, want %q", found, cfgPath)
	}
}

func TestFindReturnsFalseWhenNoConfig(t *testing.T) {
	root := t.TempDir()
	_, ok, err := Find(root)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected not to find config in empty dir")
	}
}

func TestLoadReadsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "minifyjs.config.json")
	if err := os.WriteFile(path, []byte(`{"target": "es2015"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if c.Target != "es2015" {
		t.Fatalf("target = %q", c.Target)
	}
}

func TestEnvOverlay(t *testing.T) {
	env := map[string]string{
		"MINIFYJS_TARGET":     "es2018",
		"MINIFYJS_FORMAT":     "esm",
		"MINIFYJS_SOURCEMAP":  "inline",
		"MINIFYJS_CACHE":      "1",
		"MINIFYJS_CACHE_DIR":  "/tmp/mjc",
		"MINIFYJS_QUIET":      "yes",
		"MINIFYJS_NO_MINIFY":  "true",
	}
	getenv := func(k string) string { return env[k] }

	c := Defaults()
	c.EnvOverlay(getenv)

	if c.Target != "es2018" {
		t.Fatalf("Target = %q", c.Target)
	}
	if c.Format != "esm" {
		t.Fatalf("Format = %q", c.Format)
	}
	if c.Sourcemap != "inline" {
		t.Fatalf("Sourcemap = %q", c.Sourcemap)
	}
	if !c.Cache {
		t.Fatal("Cache should be true")
	}
	if c.CacheDir != "/tmp/mjc" {
		t.Fatalf("CacheDir = %q", c.CacheDir)
	}
	if !c.Quiet {
		t.Fatal("Quiet should be true")
	}
	if c.IsMinifying() {
		t.Fatal("MINIFYJS_NO_MINIFY should turn off all minify passes")
	}
}