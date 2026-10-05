package io

import (
	"io"
	"os"
)

// WriteStdout writes b to os.Stdout and returns any error.
//
// The CLI uses this instead of fmt.Print so the exact bytes are
// preserved (no added trailing newline), which matters when output
// is piped into another tool or compared byte-for-byte in tests.
func WriteStdout(b []byte) error {
	_, err := os.Stdout.Write(b)
	return err
}

// WriteAll writes b to w, returning the number of bytes written.
func WriteAll(w io.Writer, b []byte) (int, error) {
	return w.Write(b)
}