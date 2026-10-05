package io

import "io"

// Stdout is the output stream the CLI writes minified code to when
// no -o flag is given. Overridable in tests, same as Stdin.
var Stdout io.Writer = defaultStdout()

// Stderr is where diagnostics and progress messages go.
var Stderr io.Writer = defaultStderr()