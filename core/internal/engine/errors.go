package engine

import "errors"

// Sentinel errors returned by the engine. Callers can test for these
// with errors.Is so a CLI or Python wrapper can translate them into
// stable exit codes without string matching.
var (
	// ErrInvalidOptions is returned when Options contains a value
	// that esbuild cannot accept. The wrapped error message names
	// the specific field.
	ErrInvalidOptions = errors.New("engine: invalid options")

	// ErrInvalidBuildOptions is returned when a BuildOptions value
	// is internally inconsistent (missing entry points, both OutDir
	// and OutFile set, etc).
	ErrInvalidBuildOptions = errors.New("engine: invalid build options")

	// ErrTransformFailed is returned by Transform when esbuild
	// produced errors. The Result is still populated with the
	// diagnostics so callers can render them.
	ErrTransformFailed = errors.New("engine: transform failed")

	// ErrBuildFailed is returned by Build when esbuild produced
	// errors.
	ErrBuildFailed = errors.New("engine: build failed")
)