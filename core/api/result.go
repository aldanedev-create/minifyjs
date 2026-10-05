package api

import "github.com/minifyjs/minifyjs/core/internal/diagnostics"

// Result mirrors engine.Result.
type Result struct {
	Code          string
	Map           string
	Diagnostics   []diagnostics.Diagnostic
	OriginalBytes int
	MinifiedBytes int
}

// HasErrors reports whether any diagnostic is an error.
func (r Result) HasErrors() bool {
	return diagnostics.HasErrors(r.Diagnostics)
}

// BytesSaved returns original - minified.
func (r Result) BytesSaved() int {
	return r.OriginalBytes - r.MinifiedBytes
}