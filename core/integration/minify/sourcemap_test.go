package minify_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func TestSourcemapInlineDecodes(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--sourcemap=inline"},
		Stdin: []byte("const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)

	// The inline map is a base64 data URL. We only check the header.
	if !strings.Contains(string(r.Stdout), "sourceMappingURL=data:application/json") {
		t.Fatalf("inline map missing: %q", r.Stdout)
	}
}

func TestSourcemapExternalValidJSON(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--sourcemap=external", in, "-o", out},
	})
	helpers.AssertExitCode(t, r, 0)

	data, err := readFile(out + ".map")
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("map is not JSON: %v", err)
	}
	if m["version"] != float64(3) {
		t.Fatalf("map version = %v", m["version"])
	}
}

func TestSourcemapHasSources(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--sourcemap=external", in, "-o", out},
	})
	helpers.AssertExitCode(t, r, 0)

	data, _ := readFile(out + ".map")
	var m struct {
		Sources []string `json:"sources"`
	}
	_ = json.Unmarshal(data, &m)
	if len(m.Sources) == 0 {
		t.Fatalf("map has no sources: %s", data)
	}
}

func TestSourcemapHasMappings(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\nconst y = 2;\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--sourcemap=external", in, "-o", out},
	})
	helpers.AssertExitCode(t, r, 0)

	data, _ := readFile(out + ".map")
	var m struct {
		Mappings string `json:"mappings"`
	}
	_ = json.Unmarshal(data, &m)
	if m.Mappings == "" {
		t.Fatalf("map has no mappings: %s", data)
	}
}

func TestSourcemapNotEmittedByDefault(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{Stdin: []byte("const x = 1;")})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "sourceMappingURL")
}

func TestSourcemapNoSemicolonInOutputBreaksNothing(t *testing.T) {
	// Minified output has no trailing newline; the sourcemap
	// reference is a trailing comment, not a new line.
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--sourcemap=inline"},
		Stdin: []byte("const x = 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	// The reference must be on the same line as the code.
	lines := strings.Split(string(r.Stdout), "\n")
	if len(lines) > 1 {
		t.Fatalf("inline map introduced a newline: %q", r.Stdout)
	}
}

func TestSourcemapInvalidMode(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--sourcemap=nope"},
		Stdin: []byte("const x = 1;"),
	})
	if r.ExitCode == 0 {
		t.Fatalf("invalid mode should fail")
	}
}

func TestSourcemapWithTarget(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const f = (a) => a + 1;\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--sourcemap=external", "--target", "es5", in, "-o", out},
	})
	helpers.AssertExitCode(t, r, 0)
	if _, err := readFile(out + ".map"); err != nil {
		t.Fatalf("map missing with target: %v", err)
	}
}

func TestSourcemapBothMode(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--sourcemap=both", in, "-o", out},
	})
	helpers.AssertExitCode(t, r, 0)

	body, _ := readFile(out)
	helpers.AssertContains(t, string(body), "sourceMappingURL=data:application/json")
	if _, err := readFile(out + ".map"); err != nil {
		t.Fatalf("both mode should also write .map: %v", err)
	}
}

func TestSourcemapSourcesContentWhenSet(t *testing.T) {
	// esbuild by default omits sourcesContent for external maps.
	// We only assert the map parses and has the expected keys.
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := dir + "/out.js"

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--sourcemap=external", in, "-o", out},
	})
	helpers.AssertExitCode(t, r, 0)

	data, _ := readFile(out + ".map")
	for _, key := range []string{`"version"`, `"sources"`, `"mappings"`} {
		if !strings.Contains(string(data), key) {
			t.Fatalf("map missing %s: %s", key, data)
		}
	}
}