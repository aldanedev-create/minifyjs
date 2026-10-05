// Package diagnostics defines the structured error/warning type that
// every stage of the pipeline produces, plus formatters that turn
// them into terminal or JSON output.
//
// The CLI uses the terminal formatter; the Python bridge reads the
// JSON formatter's output. Keeping both in one package guarantees the
// two presentations stay in sync.
package diagnostics

// Severity classifies how serious a Diagnostic is.
type Severity int

const (
	// SeverityInfo is purely informational and never fails a build.
	SeverityInfo Severity = iota
	// SeverityWarning indicates something suspicious but not
	// incorrect (e.g. a no-op optimization the user asked for).
	SeverityWarning
	// SeverityError indicates the input could not be processed.
	SeverityError
)

// String returns the lowercase name of the severity, used in JSON
// output and in terminal prefixes.
func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeverityWarning:
		return "warning"
	case SeverityError:
		return "error"
	default:
		return "unknown"
	}
}