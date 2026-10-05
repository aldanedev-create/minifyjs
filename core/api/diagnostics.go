package api

// This file re-exports the diagnostics types from
// internal/diagnostics so external callers who import only the
// public API do not need to reach into internal/ (which the Go
// compiler forbids for code outside this module).
//
// The re-exported names are stable. Their underlying definitions may
// evolve, but the names in this package will not be removed without a
// major version bump.

import (
	"github.com/minifyjs/minifyjs/core/internal/diagnostics"
)

// Severity classifies how serious a Diagnostic is.
type Severity = diagnostics.Severity

// Severity values. Aliased rather than redeclared so they remain
// comparable and switch-compatible with values produced by the
// engine.
const (
	SeverityInfo    = diagnostics.SeverityInfo
	SeverityWarning = diagnostics.SeverityWarning
	SeverityError   = diagnostics.SeverityError
)

// Diagnostic is a single message produced by the engine.
//
// It carries enough context for both a terminal user and a machine
// consumer: a severity, a stable message, an optional source
// location, and an optional stable code that scripts can switch on.
type Diagnostic = diagnostics.Diagnostic

// Formatter renders a slice of Diagnostics to a writer.
type Formatter = diagnostics.Formatter

// TerminalFormatter prints diagnostics in the compact "file:line:col:
// severity: message" format used by the CLI.
type TerminalFormatter = diagnostics.TerminalFormatter

// JSONFormatter emits diagnostics as a JSON array. This is the
// format the Python bridge consumes when structured output is
// requested.
type JSONFormatter = diagnostics.JSONFormatter

// HasErrors reports whether any diagnostic in the slice is an error.
// It is a package-level function rather than a method so callers can
// use it on a nil slice without constructing a Result first.
func HasErrors(diags []Diagnostic) bool {
	return diagnostics.HasErrors(diags)
}

// SortStable orders diagnostics deterministically: by file, then
// line, then column, then severity. It is exported so external
// callers can produce the same ordering the CLI does.
func SortStable(diags []Diagnostic) {
	diagnostics.SortStable(diags)
}