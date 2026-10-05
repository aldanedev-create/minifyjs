package diagnostics

import (
	"fmt"
	"io"
)

// TerminalFormatter prints diagnostics in the same compact format the
// CLI has always used:
//
//	file:line:col: severity: message
//
// Diagnostics without a location omit the "file:line:col:" prefix.
type TerminalFormatter struct {
	// Prefix, if non-empty, is printed before every line. The CLI
	// sets this to "minifyjs: " so its output is greppable.
	Prefix string
}

// Format implements Formatter.
func (f TerminalFormatter) Format(w io.Writer, diags []Diagnostic) error {
	SortStable(diags)
	for _, d := range diags {
		var loc string
		if d.HasLocation() {
			if d.File != "" {
				loc = fmt.Sprintf("%s:%d:%d: ", d.File, d.Line, d.Col)
			} else {
				loc = fmt.Sprintf("%d:%d: ", d.Line, d.Col)
			}
		}
		line := fmt.Sprintf("%s%s%s: %s\n", f.Prefix, loc, d.Severity, d.Message)
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}