// Package helpers contains shared utilities for MinifyJS's
// integration tests. Integration tests exercise the compiled CLI
// binary and the full API surface, as opposed to unit tests that
// exercise individual internal packages.
//
// The CLI is compiled exactly once per `go test` invocation and
// cached for the duration. Compiling it fresh for every test would
// add tens of seconds to a suite that already runs an external
// process per test.
package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

// buildOnce ensures the CLI is built at most once per test binary
// run. The compiled binary is placed in a temp dir that lives as
// long as the test process.
var (
	buildOnce     sync.Once
	builtCLIPath  string
	buildFailErr  error
)

// CLIPath returns the path to the compiled minifyjs binary, building
// it on first call. It fails the test if the build fails.
//
// The path is cached for the lifetime of the test binary, so calling
// it from many tests is cheap.
func CLIPath(t *testing.T) string {
	t.Helper()
	buildOnce.Do(func() {
		builtCLIPath, buildFailErr = buildCLI()
	})
	if buildFailErr != nil {
		t.Fatalf("build CLI: %v", buildFailErr)
	}
	return builtCLIPath
}

func buildCLI() (string, error) {
	// The test file lives at core/integration/helpers/build.go.
	// Two levels up lands at core/, where cmd/minifyjs is.
	_, thisFile, _, _ := runtime.Caller(0)
	coreRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")

	out := filepath.Join(os.TempDir(), "minifyjs-test-"+pidStr())
	if runtime.GOOS == "windows" {
		out += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", out, "./cmd/minifyjs")
	cmd.Dir = coreRoot
	cmd.Env = os.Environ()
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", errWithOutput(err, output)
	}
	return out, nil
}

func pidStr() string {
	return strconvItoa(os.Getpid())
}