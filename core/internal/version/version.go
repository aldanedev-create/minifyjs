// Package version is the single source of truth for the MinifyJS
// engine's version string.
//
// Both the native CLI (core/cmd/minifyjs) and the public Go API
// (core/api) read the version from here so there is exactly one
// place to change it during a release. The Python package
// (python/minifyjs/_version.py) keeps its own copy because it ships
// independently of the Go source, but the release pipeline keeps the
// two in sync (see build/build_release.py).
package version

// Version is the semantic version of the MinifyJS engine.
//
// The value here is the fallback used when the binary is built
// without -ldflags (e.g. `go run ./cmd/minifyjs` from a source
// checkout). Release builds override it at link time:
//
//	go build -ldflags "-X github.com/minifyjs/minifyjs/core/internal/version.Version=v1.2.3"
//
// The variable is therefore declared as a var, not a const: a const
// cannot be overridden by -ldflags. Everything else in the codebase
// reads version.Version and gets whichever value was baked in at
// build time.
//
// It follows semver: MAJOR.MINOR.PATCH, optionally followed by a
// pre-release suffix.
var Version = "0.1.2"

// String returns the version string. It exists so callers can write
// version.String() instead of version.Version, which gives us a
// stable API if we later decide to derive the value from build
// metadata rather than an ldflag.
func String() string {
	return Version
}