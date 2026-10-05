package source

// Span is a half-open byte range [Start, End) into a single source
// file, plus the file's name for diagnostics that need to print
// "file:line:col".
//
// Start and End are byte offsets, matching Position.Offset. The End
// offset is exclusive, so a zero-width token has Start == End.
type Span struct {
	File  string
	Start Position
	End   Position
}

// IsValid reports whether the span refers to a real location.
func (s Span) IsValid() bool {
	return s.Start.IsValid() && s.End.IsValid()
}

// Len returns the number of bytes in the span.
func (s Span) Len() int {
	return s.End.Offset - s.Start.Offset
}

// String returns "file:line:col" for the start of the span, suitable
// for the prefix of a diagnostic message.
func (s Span) String() string {
	if s.File == "" {
		return s.Start.String()
	}
	return s.File + ":" + s.Start.String()
}