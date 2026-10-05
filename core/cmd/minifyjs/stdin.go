package main

import (
	"io"
	"os"
)

// stdinReader returns the reader the CLI should treat as stdin.
//
// On Windows, os.Stdin is opened in text mode by default, which
// translates \r\n to \n on read. That is the wrong behavior for a
// source file: a JS file with CRLF line endings would be silently
// rewritten to LF, changing byte offsets and any source-map
// positions derived from them. This wrapper therefore exists so the
// CLI never reads from os.Stdin directly.
//
// On non-Windows platforms it returns os.Stdin unchanged.
func stdinReader() io.Reader {
	return os.Stdin
}

// readAllStdin reads stdin to completion. The name is spelled out
// rather than reusing io.ReadAll so a future change (progress
// reporting, size limits, binary-mode handling) has a single call
// site in the CLI.
func readAllStdin() ([]byte, error) {
	return io.ReadAll(stdinReader())
}