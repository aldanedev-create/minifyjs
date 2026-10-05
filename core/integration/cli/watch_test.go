package cli_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

// Watch tests use a short timeout. The CLI runs forever; the test
// kills it after observing the behavior it wants to verify.

func TestWatchRequiresInput(t *testing.T) {
	dir := helpers.TmpDir(t)
	out := filepath.Join(dir, "out.js")
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--watch", "-o", out},
	})
	helpers.AssertExitCode(t, r, 2)
}

func TestWatchRequiresOutput(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	r := helpers.Run(t, helpers.RunOptions{
		Args: []string{"--watch", in},
	})
	helpers.AssertExitCode(t, r, 2)
}

func TestWatchInitialBuild(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const x = 1;\n")
	out := filepath.Join(dir, "out.js")

	// --watch runs forever; use a helper that kills it after a
	// short timeout and returns whatever output was produced.
	helpers.RunWithTimeout(t, 2*time.Second, helpers.RunOptions{
		Args: []string{"--watch", in, "-o", out},
	})

	if _, err := os.Stat(out); err != nil {
		t.Fatalf("watch mode did not produce output: %v", err)
	}
}

func TestWatchRebuildsOnChange(t *testing.T) {
	dir := helpers.TmpDir(t)
	in := helpers.WriteString(t, dir, "in.js", "const first = 1;\n")
	out := filepath.Join(dir, "out.js")

	done := make(chan struct{})
	go func() {
		helpers.RunWithTimeout(t, 3*time.Second, helpers.RunOptions{
			Args: []string{"--watch", in, "-o", out},
		})
		close(done)
	}()

	// Give the watcher time to do its initial build.
	time.Sleep(500 * time.Millisecond)

	// Change the input.
	if err := os.WriteFile(in, []byte("const second = 2;\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Give it time to pick up the change.
	time.Sleep(1500 * time.Millisecond)

	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(string(body), "second") {
		t.Fatalf("watch did not rebuild after change: %q", body)
	}

	<-done
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) &&
		(haystack == needle || containsRec(haystack, needle))
}

func containsRec(h, n string) bool {
	for i := 0; i+len(n) <= len(h); i++ {
		if h[i:i+len(n)] == n {
			return true
		}
	}
	return false
}