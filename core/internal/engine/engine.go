package engine

// engine.go holds the package-level documentation and the two public
// entry points are declared in transform.go and build.go. This file
// exists so godoc shows the package overview first.

// Version is the esbuild version this build of MinifyJS wraps. It is
// exposed so the CLI's --version output and the Python package can
// report which esbuild engine is in use, which matters when users
// file bug reports.
//
// The value is set via a Go build ldflag in the release pipeline
// (see build/build_core.py) so it always matches the pinned version
// in go.mod, even if someone forgets to update this file by hand.
var Version = "0.28.2"