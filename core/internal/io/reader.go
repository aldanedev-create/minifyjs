// Package io provides filesystem and stdio helpers shared by the CLI,
// the runtime, and (later) any long-running watch/server mode.
//
// It deliberately avoids importing core/internal/* so that the CLI
// can use it without pulling the whole engine into the argument
// parser.
package io

import (
	"io"
	"os"
)

// ReadAll reads an io.Reader to completion and returns the bytes. It
// exists so callers do not each need to import io and os in the same
// file.
func ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// ReadFile reads an entire file. Returns os.ReadFile's error
// unchanged so callers can inspect it with os.IsNotExist.
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadStdin reads os.Stdin to completion.
func ReadStdin() ([]byte, error) {
	return io.ReadAll(os.Stdin)
}