//go:build !testoverride
// +build !testoverride

package io

import (
	"io"
	"os"
)

// defaultStdin / defaultStdout / defaultStderr return the real
// process streams. Kept in a build-tagged file so a test binary can
// substitute fakes if needed without affecting production builds.
func defaultStdin() io.Reader  { return os.Stdin }
func defaultStdout() io.Writer { return os.Stdout }
func defaultStderr() io.Writer { return os.Stderr }