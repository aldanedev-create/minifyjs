package main

import (
	"io"

	miniio "github.com/minifyjs/minifyjs/core/internal/io"
)

// readInput reads from stdin or from a file path. The empty string
// and "-" both mean stdin.
//
// readInput is used by the watch loop and by tests; the main
// runTransform path inlines the same logic so it can produce more
// specific error messages.
func readInput(path string, stdin io.Reader) ([]byte, string, error) {
	if path == "" || path == "-" {
		b, err := miniio.ReadAll(stdin)
		return b, "<stdin>", err
	}
	b, err := miniio.ReadFile(path)
	return b, path, err
}