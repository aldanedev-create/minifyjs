package main

import (
	"fmt"

	"github.com/minifyjs/minifyjs/core/internal/engine"
	"github.com/minifyjs/minifyjs/core/internal/version"
)

// printVersion writes the version banner. It reports both the
// MinifyJS version and the esbuild version it wraps, because bug
// reports about minification behavior usually need both.
func printVersion(out writer) {
	fmt.Fprintf(out, "minifyjs %s (esbuild %s)\n", version.Version, engine.Version)
}