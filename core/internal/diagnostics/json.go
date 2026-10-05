package diagnostics

import (
	"encoding/json"
	"io"
)

// JSONFormatter emits diagnostics as a JSON array. This is what the
// Python bridge consumes when it wants structured output (the CLI
// gains a --diagnostics-json flag for this in a later phase).
type JSONFormatter struct {
	// Pretty enables indented output; default (false) is compact.
	Pretty bool
}

type jsonDiagnostic struct {
	Severity string `json:"severity"`
	Code     string `json:"code,omitempty"`
	Message  string `json:"message"`
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Col      int    `json:"col,omitempty"`
}

// Format implements Formatter.
func (f JSONFormatter) Format(w io.Writer, diags []Diagnostic) error {
	SortStable(diags)
	out := make([]jsonDiagnostic, 0, len(diags))
	for _, d := range diags {
		out = append(out, jsonDiagnostic{
			Severity: d.Severity.String(),
			Code:     d.Code,
			Message:  d.Message,
			File:     d.File,
			Line:     d.Line,
			Col:      d.Col,
		})
	}
	enc := json.NewEncoder(w)
	if f.Pretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(out)
}