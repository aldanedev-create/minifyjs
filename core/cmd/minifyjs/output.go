package main

import (
	"io"

	miniio "github.com/minifyjs/minifyjs/core/internal/io"
)

// writeFile is a small wrapper used by the watch loop. The main
// runTransform path calls miniio.WriteFileAtomic directly.
func writeFile(path string, data []byte) error {
	return miniio.WriteFileAtomic(path, data, 0o644)
}

// writeStdout writes to stdout. Kept as a named function so tests
// and future variants (e.g. --append) have a single place to hook.
func writeStdout(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}