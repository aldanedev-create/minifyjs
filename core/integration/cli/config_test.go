package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func pathEnv() string {
	return os.Getenv("PATH")
}

// ---------------------------------------------------------------------------
// Basic discovery
// ---------------------------------------------------------------------------

// TestConfigFileDiscovered: a minifyjs.config.json in the working
// directory is picked up automatically.
func TestConfigFileDiscovered(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{"target":"es5"}`)
	helpers.WriteString(t, dir, "in.js", "const f = (a) => a;")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	// es5 target lowers arrow functions.
	helpers.AssertNotContains(t, string(r.Stdout), "=>")
}

// TestConfigFileDiscoveredUpTree: a config in a parent directory is
// found when running from a subdirectory.
func TestConfigFileDiscoveredUpTree(t *testing.T) {
	root := helpers.TmpDir(t)
	helpers.WriteString(t, root, "minifyjs.config.json", `{"target":"es5"}`)
	sub := filepath.Join(root, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	helpers.WriteString(t, sub, "in.js", "const f = (a) => a;")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  sub,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "=>")
}

// TestConfigRcJsonVariant: the .minifyjsrc.json filename is also
// recognized.
func TestConfigRcJsonVariant(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, ".minifyjsrc.json", `{"target":"es5"}`)
	helpers.WriteString(t, dir, "in.js", "const f = (a) => a;")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "=>")
}

// TestConfigRcVariant: the .minifyjsrc filename (no extension) is
// also recognized.
func TestConfigRcVariant(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, ".minifyjsrc", `{"target":"es5"}`)
	helpers.WriteString(t, dir, "in.js", "const f = (a) => a;")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "=>")
}

// TestConfigDiscoveryOrder: minifyjs.config.json wins over
// .minifyjsrc.json in the same directory.
func TestConfigDiscoveryOrder(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{"target":"es5"}`)
	helpers.WriteString(t, dir, ".minifyjsrc.json", `{"target":"esnext"}`)
	helpers.WriteString(t, dir, "in.js", "const f = (a) => a;")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	// minifyjs.config.json says es5, so arrows are lowered. If the
	// .minifyjsrc.json had won, arrows would still be present.
	helpers.AssertNotContains(t, string(r.Stdout), "=>")
}

// TestNoConfigFileIsFine: no config file anywhere is not an error.
func TestNoConfigFileIsFine(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "const x")
}

// ---------------------------------------------------------------------------
// --config and --no-config
// ---------------------------------------------------------------------------

// TestConfigFlagOverridesFile: --config points at a specific file,
// which takes precedence over the discovered one.
func TestConfigFlagOverridesFile(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{"target":"es5"}`)
	helpers.WriteString(t, dir, "custom.json", `{"target":"esnext"}`)
	helpers.WriteString(t, dir, "in.js", "const f = (a) => a;")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--config", "custom.json", "in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	// custom.json says esnext, so arrows survive.
	helpers.AssertContains(t, string(r.Stdout), "=>")
}

// TestConfigFlagWithMissingFile: --config pointing at a missing file
// is an error.
func TestConfigFlagWithMissingFile(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--config", "nope.json", "in.js"},
		Dir:  dir,
	})
	if r.ExitCode == 0 {
		t.Fatalf("missing config file should fail")
	}
	if !strings.Contains(strings.ToLower(string(r.Stderr)), "config") {
		t.Fatalf("stderr should mention config: %s", r.Stderr)
	}
}

// TestNoConfigIgnoresFile: --no-config disables discovery entirely.
func TestNoConfigIgnoresFile(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{"target":"es5"}`)
	helpers.WriteString(t, dir, "in.js", "const f = (a) => a;")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--no-config", "in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	// es5 was in the file but --no-config ignores it.
	helpers.AssertContains(t, string(r.Stdout), "=>")
}

// TestNoConfigIgnoresExplicitConfig: --no-config wins over --config.
func TestNoConfigIgnoresExplicitConfig(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "custom.json", `{"target":"es5"}`)
	helpers.WriteString(t, dir, "in.js", "const f = (a) => a;")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--no-config", "--config", "custom.json", "in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "=>")
}

// ---------------------------------------------------------------------------
// Precedence: flags > env > file > defaults
// ---------------------------------------------------------------------------

// TestFlagOverridesConfigTarget: --target on the command line wins.
func TestFlagOverridesConfigTarget(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{"target":"es5"}`)
	helpers.WriteString(t, dir, "in.js", "const f = (a) => a;")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--target", "esnext", "in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "=>")
}

// TestEnvOverridesConfigFile: MINIFYJS_TARGET wins over the file.
func TestEnvOverridesConfigFile(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{"target":"es5"}`)
	helpers.WriteString(t, dir, "in.js", "const f = (a) => a;")

	env := []string{
		"MINIFYJS_TARGET=esnext",
		"PATH=" + pathEnv(),
	}
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
		Env:  env,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "=>")
}

// TestFlagOverridesEnv: --target beats MINIFYJS_TARGET.
func TestFlagOverridesEnv(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "in.js", "const f = (a) => a;")

	env := []string{
		"MINIFYJS_TARGET=es5",
		"PATH=" + pathEnv(),
	}
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--target", "esnext", "in.js"},
		Dir:  dir,
		Env:  env,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "=>")
}

// TestEnvOverridesDefaults: MINIFYJS_* applies when no file and no
// flag is present.
func TestEnvOverridesDefaults(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "in.js", "const f = (a) => a;")

	env := []string{
		"MINIFYJS_TARGET=es5",
		"PATH=" + pathEnv(),
	}
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
		Env:  env,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "=>")
}

// ---------------------------------------------------------------------------
// Config file shape
// ---------------------------------------------------------------------------

// TestConfigMinifyBlock: nested minify object works.
func TestConfigMinifyBlock(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json",
		`{"minify":{"whitespace":false,"identifiers":false,"syntax":false}}`)
	helpers.WriteString(t, dir, "in.js", "const x = 1 + 2;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	// No minify passes means no constant folding: 1+2 stays.
	helpers.AssertContains(t, string(r.Stdout), "1 + 2")
}

// TestConfigPartialMinify: only syntax on, whitespace off.
func TestConfigPartialMinify(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json",
		`{"minify":{"whitespace":false,"identifiers":false,"syntax":true}}`)
	helpers.WriteString(t, dir, "in.js", "const x = 1 + 2 + 3;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	// Syntax folding turns 1+2+3 into 6.
	helpers.AssertContains(t, string(r.Stdout), "6")
}

// TestConfigBanner: banner from config appears in output.
func TestConfigBanner(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json",
		`{"banner":"/* config banner */"}`)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "config banner")
}

// TestConfigOutput: output path from config is used.
func TestConfigOutput(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{"output":"out.js"}`)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	if _, err := os.Stat(filepath.Join(dir, "out.js")); err != nil {
		t.Fatalf("config output not written: %v", err)
	}
}

// TestConfigSourcemap: sourcemap mode from config is honored.
func TestConfigSourcemap(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{"sourcemap":"inline"}`)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "sourceMappingURL=data:application/json")
}

// TestConfigRelativePathsResolvedAgainstConfigFile: an input path in
// the config is resolved relative to the config file, not the cwd.
func TestConfigRelativePathsResolvedAgainstConfigFile(t *testing.T) {
	root := helpers.TmpDir(t)
	sub := filepath.Join(root, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	helpers.WriteString(t, root, "minifyjs.config.json",
		`{"output":"dist/out.js"}`)
	helpers.WriteString(t, sub, "in.js", "const x = 1;\n")

	// Run from the subdirectory. The config is discovered at root,
	// so "dist/out.js" resolves to root/dist/out.js, not
	// sub/dist/out.js.
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  sub,
	})
	helpers.AssertExitCode(t, r, 0)

	expected := filepath.Join(root, "dist", "out.js")
	if _, err := os.Stat(expected); err != nil {
		t.Fatalf("config output not at %s: %v", expected, err)
	}
}

// ---------------------------------------------------------------------------
// Malformed config
// ---------------------------------------------------------------------------

// TestConfigInvalidJSON: a config that is not valid JSON produces an
// error, not a silent fallback.
func TestConfigInvalidJSON(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{not json}`)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	if r.ExitCode == 0 {
		t.Fatalf("invalid config should fail")
	}
	if !strings.Contains(strings.ToLower(string(r.Stderr)), "config") {
		t.Fatalf("stderr should mention config: %s", r.Stderr)
	}
}

// TestConfigInvalidTarget: a syntactically valid JSON file with an
// unknown target fails during validation.
func TestConfigInvalidTarget(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{"target":"ES2015"}`)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	if r.ExitCode == 0 {
		t.Fatalf("invalid target should fail")
	}
}

// TestConfigQuietAndVerboseConflict: config validation rejects
// quiet+verbose together.
func TestConfigQuietAndVerboseConflict(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json",
		`{"quiet":true,"verbose":true}`)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	if r.ExitCode == 0 {
		t.Fatalf("quiet+verbose should fail validation")
	}
}

// TestConfigUnknownKeysIgnored: unknown keys are tolerated so a
// config written for a newer version does not fail on an older one.
func TestConfigUnknownKeysIgnored(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json",
		`{"futureOption":"whatever","target":"esnext"}`)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
}

// TestConfigEmptyObject: an empty config object is valid.
func TestConfigEmptyObject(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "minifyjs.config.json", `{}`)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
	})
	helpers.AssertExitCode(t, r, 0)
}

// ---------------------------------------------------------------------------
// Env vars beyond MINIFYJS_TARGET
// ---------------------------------------------------------------------------

// TestEnvFormat: MINIFYJS_FORMAT is honored.
func TestEnvFormat(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "in.js", "export const x = 1;\n")

	env := []string{"MINIFYJS_FORMAT=esm", "PATH=" + pathEnv()}
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
		Env:  env,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "export")
}

// TestEnvSourcemap: MINIFYJS_SOURCEMAP is honored.
func TestEnvSourcemap(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")

	env := []string{"MINIFYJS_SOURCEMAP=inline", "PATH=" + pathEnv()}
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
		Env:  env,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "sourceMappingURL=data:application/json")
}

// TestEnvNoMinify: MINIFYJS_NO_MINIFY=1 turns off all minify passes.
func TestEnvNoMinify(t *testing.T) {
	dir := helpers.TmpDir(t)
	src := "function f() {\n    return 1;\n}\n"
	helpers.WriteString(t, dir, "in.js", src)

	env := []string{"MINIFYJS_NO_MINIFY=1", "PATH=" + pathEnv()}
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js"},
		Dir:  dir,
		Env:  env,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertEqual(t, string(r.Stdout), src)
}

// TestEnvQuiet: MINIFYJS_QUIET=1 suppresses progress output.
func TestEnvQuiet(t *testing.T) {
	dir := helpers.TmpDir(t)
	helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := filepath.Join(dir, "out.js")

	env := []string{"MINIFYJS_QUIET=1", "PATH=" + pathEnv()}
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"in.js", "-o", out},
		Dir:  dir,
		Env:  env,
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNoStderr(t, r)
}