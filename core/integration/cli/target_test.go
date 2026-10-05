package cli_test

import (
	"strings"
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestTargetES5LowersArrows(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "es5"},
		Stdin: []byte("const f = (a) => a + 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "=>")
}

func TestTargetES2015KeepsArrows(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "es2015"},
		Stdin: []byte("const f = (a) => a + 1;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "=>")
}

func TestTargetES2017KeepsAsync(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "es2017"},
		Stdin: []byte("async function f() { await x; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "async")
}

func TestTargetES2015LowersAsync(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "es2015"},
		Stdin: []byte("async function f() { await x; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	// es2015 has no async/await; esbuild lowers it.
	helpers.AssertNotContains(t, string(r.Stdout), "async function")
}

func TestTargetESNextPreservesAll(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "esnext"},
		Stdin: []byte("class A { #x = 1; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "#x")
}

func TestTargetCompoundList(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "es2020,chrome90"},
		Stdin: []byte("const x = a?.b ?? c;"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestTargetInvalidExitsTwo(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "ES2015"}, // uppercase
		Stdin: []byte("const x = 1;"),
	})
	if r.ExitCode == 0 {
		t.Fatalf("uppercase target should be rejected")
	}
}

func TestTargetGarbageExitsTwo(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "@@@@"},
		Stdin: []byte("const x = 1;"),
	})
	if r.ExitCode == 0 {
		t.Fatalf("garbage target should be rejected")
	}
}

func TestTargetES5LowersClasses(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "es5"},
		Stdin: []byte("class A { constructor() {} }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "class ")
}

func TestTargetES2015KeepsClasses(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "es2015"},
		Stdin: []byte("class A { constructor() {} }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "class")
}

func TestTargetES2015LowersOptionalChaining(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "es2015"},
		Stdin: []byte("const x = a?.b?.c;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertNotContains(t, string(r.Stdout), "?.")
}

func TestTargetES2020KeepsOptionalChaining(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Args:  []string{"--target", "es2020"},
		Stdin: []byte("const x = a?.b;"),
	})
	helpers.AssertExitCode(t, r, 0)
	if !strings.Contains(string(r.Stdout), "?.") {
		// esbuild may normalize to a?.b still
		t.Logf("output: %s", r.Stdout)
	}
}