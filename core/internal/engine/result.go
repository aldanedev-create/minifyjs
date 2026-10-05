package engine

import (
	"github.com/minifyjs/minifyjs/core/internal/diagnostics"
)

// Result is what Transform and Build return. It is deliberately
// independent of esbuild's own result type so the rest of MinifyJS
// does not need to import esbuild.
type Result struct {
	// Code is the transformed JavaScript source.
	Code string

	// Map is the source map, when one was requested. Its format is
	// determined by Options.Sourcemap: for "inline" it is appended
	// to Code as a data URL and Map is empty; for "external" and
	// "both" it is the raw JSON of the map.
	Map string

	// Diagnostics is the ordered list of errors and warnings
	// produced by esbuild, converted into MinifyJS's own diagnostic
	// type. Errors here mean Code is not usable; callers should
	// check HasErrors before using Code.
	Diagnostics []diagnostics.Diagnostic

	// OriginalBytes is the byte length of the input.
	OriginalBytes int

	// MinifiedBytes is the byte length of Code.
	MinifiedBytes int
}

// HasErrors reports whether any diagnostic is an error.
func (r Result) HasErrors() bool {
	return diagnostics.HasErrors(r.Diagnostics)
}

// HasWarnings reports whether any diagnostic is a warning (and none
// are errors).
func (r Result) HasWarnings() bool {
	for _, d := range r.Diagnostics {
		if d.Severity == diagnostics.SeverityWarning {
			return true
		}
	}
	return false
}

// Ratio returns the compression ratio as minified/original (0..1).
// Returns 0 when the input was empty so callers do not divide by
// zero.
func (r Result) Ratio() float64 {
	if r.OriginalBytes == 0 {
		return 0
	}
	return float64(r.MinifiedBytes) / float64(r.OriginalBytes)
}

// BytesSaved returns original - minified. Negative when the output
// grew (which happens with tiny inputs and banner/footer options).
func (r Result) BytesSaved() int {
	return r.OriginalBytes - r.MinifiedBytes
}