package diagnostics

import (
	"io"
	"sort"
)

// Formatter renders a slice of Diagnostics to a writer in a specific
// format. Implementations must be safe to reuse across calls.
type Formatter interface {
	Format(w io.Writer, diags []Diagnostic) error
}

// SortStable orders diagnostics deterministically: by file, then
// line, then column, then severity (errors before warnings before
// info). This makes CLI output reproducible regardless of the order
// the pipeline happened to produce the diagnostics in.
func SortStable(diags []Diagnostic) {
	sort.SliceStable(diags, func(i, j int) bool {
		a, b := diags[i], diags[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if a.Col != b.Col {
			return a.Col < b.Col
		}
		return a.Severity > b.Severity
	})
}

// HasErrors reports whether any diagnostic is an error.
func HasErrors(diags []Diagnostic) bool {
	for _, d := range diags {
		if d.IsError() {
			return true
		}
	}
	return false
}