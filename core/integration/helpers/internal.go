package helpers

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

func runtimeCaller() (pc uintptr, file string, line int, ok bool) {
	return runtime.Caller(1)
}

func errWithOutput(err error, output []byte) error {
	return fmt.Errorf("%w\n%s", err, output)
}

func strconvItoa(n int) string {
	return strconv.Itoa(n)
}

// MustGetenv is a test helper that fails the test if the variable is
// not set. It exists so tests that depend on CI-provided values
// (e.g. a fixture path) fail with a clear message rather than
// silently using the empty string.
func MustGetenv(t interface{ Fatal(...any) }, key string) string {
	v := os.Getenv(key)
	if v == "" {
		t.Fatal("environment variable not set:", key)
	}
	return v
}