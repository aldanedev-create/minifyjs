package diagnostics

// Diagnostic is a single message produced by the engine.
//
// It carries enough context for both a terminal user and a machine
// consumer: a severity, a stable message, an optional source
// location, and an optional stable code (e.g. "lex.unterminated-string")
// that scripts and the Python API can switch on.
type Diagnostic struct {
	Severity Severity
	// Code is a stable identifier like "lex.unterminated-string".
	// Empty when the diagnostic has no stable code yet.
	Code string
	// Message is the human-readable description.
	Message string
	// File/Line/Col locate the diagnostic in the source, if known.
	// Zero values mean "no location".
	File string
	Line int
	Col  int
}

// IsError reports whether the diagnostic is an error (as opposed to
// info or warning).
func (d Diagnostic) IsError() bool { return d.Severity == SeverityError }

// HasLocation reports whether the diagnostic has a source location.
func (d Diagnostic) HasLocation() bool {
	return d.Line > 0 && d.Col > 0
}

// NewError builds an error diagnostic.
func NewError(code, msg string) Diagnostic {
	return Diagnostic{Severity: SeverityError, Code: code, Message: msg}
}

// NewWarning builds a warning diagnostic.
func NewWarning(code, msg string) Diagnostic {
	return Diagnostic{Severity: SeverityWarning, Code: code, Message: msg}
}

// NewInfo builds an informational diagnostic.
func NewInfo(code, msg string) Diagnostic {
	return Diagnostic{Severity: SeverityInfo, Code: code, Message: msg}
}

// WithLocation returns a copy of d with the location filled in.
func (d Diagnostic) WithLocation(file string, line, col int) Diagnostic {
	d.File = file
	d.Line = line
	d.Col = col
	return d
}