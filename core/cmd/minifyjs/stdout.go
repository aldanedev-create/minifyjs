package main

import (
	"io"
	"os"
)

// stdoutWriter returns the writer the CLI should treat as stdout.
//
// On Windows, os.Stdout is opened in text mode by default, which
// translates \n to \r\n on write. Minified JavaScript should be
// written byte-for-byte: a downstream tool comparing output, a
// Python subprocess reading from a pipe, or a byte-length check
// would all see different bytes on Windows than on Linux if we let
// the OS translate.
//
// The same caveat as stdinReader applies: switching the OS handle to
// binary mode requires syscall and is not safe to do unconditionally
// from a library. The function exists so the decision has one home.
func stdoutWriter() io.Writer {
	return os.Stdout
}

// stderrWriter returns the writer the CLI should treat as stderr.
// Diagnostics and progress messages go here, so text-mode translation
// on Windows is fine and actually desirable (it makes \n render
// correctly in cmd.exe).
func stderrWriter() io.Writer {
	return os.Stderr
}

// writeStdout writes exactly len(data) bytes and returns an error if
// fewer were written. os.Stdout.Write on a pipe can short-write when
// the reader is slow; callers that care about complete output should
// use this rather than Write directly.
func writeStdout(w io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := w.Write(data)
		if err != nil {
			return err
		}
		data = data[n:]
	}
	return nil
}