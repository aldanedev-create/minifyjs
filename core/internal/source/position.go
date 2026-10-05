// Package source owns source-file bookkeeping: byte offsets, line and
// column numbers, spans, and a line map used to translate between
// byte offsets and human-readable positions.
//
// Every later phase (diagnostics, source maps, error reporting) needs
// this information, so it lives in its own package with no
// dependencies on the rest of the engine.
package source

import "strconv"

// Position is a 1-based line/column pair plus a 0-based byte offset
// into the source text.
//
// Line and Col are what humans see in error messages and what
// source maps encode into VLQ mappings. Offset is what the engine
// uses internally because it is cheap to compare and slice.
type Position struct {
	Offset int // 0-based byte offset
	Line   int // 1-based line number
	Col    int // 1-based column number (in bytes, not runes)
}

// IsValid reports whether the position has been initialized. The zero
// Position is considered invalid so that "no position" can be
// represented without a pointer.
func (p Position) IsValid() bool {
	return p.Line > 0 && p.Col > 0
}

// String returns a compact "line:col" form suitable for error
// messages and log lines.
func (p Position) String() string {
	if !p.IsValid() {
		return "?:?"
	}
	return strconv.Itoa(p.Line) + ":" + strconv.Itoa(p.Col)
}