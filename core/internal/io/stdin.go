package io

import "io"

// Stdin is the input stream the CLI reads from when no file argument
// is given. It is a package-level variable so tests can replace it
// with a strings.Reader without touching the process's real stdin.
var Stdin io.Reader = defaultStdin()