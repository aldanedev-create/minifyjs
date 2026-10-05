// Command minifyjs is the native CLI for MinifyJS.
//
// Usage:
//
//	minifyjs [input.js] [flags]
//	cat input.js | minifyjs [flags]
//	minifyjs --bundle src/main.js --outdir dist
//
// The CLI is a thin wrapper around the public Go API in core/api. All
// option resolution (defaults → config file → env vars → flags)
// happens here; the engine does not read config or environment.
package main

import (
	"os"
)

func main() {
	os.Exit(run(os.Args[1:]))
}