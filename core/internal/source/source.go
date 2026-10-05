package source

// Source is a named, immutable JavaScript source file together with
// its precomputed LineMap.
//
// Every stage of the pipeline (transform, diagnostics, source-map
// generation) carries a *Source so that any diagnostic can be
// translated back to a real "file:line:col" without re-scanning the
// text.
type Source struct {
	name string
	text string
	line *LineMap
}

// New creates a Source from a name (used in diagnostics; may be empty
// for stdin) and the raw text.
func New(name, text string) *Source {
	return &Source{
		name: name,
		text: text,
		line: NewLineMap(text),
	}
}

// Name returns the source name (file path, "<stdin>", etc).
func (s *Source) Name() string { return s.name }

// Text returns the raw source text.
func (s *Source) Text() string { return s.text }

// Len returns the byte length of the source.
func (s *Source) Len() int { return len(s.text) }

// LineMap returns the precomputed line map.
func (s *Source) LineMap() *LineMap { return s.line }

// PositionAt converts a byte offset into a Position.
func (s *Source) PositionAt(offset int) Position {
	return s.line.Position(offset)
}

// LocationAt builds a Location for a byte offset.
func (s *Source) LocationAt(offset int) Location {
	return Location{File: s.name, Position: s.line.Position(offset)}
}

// SpanFromOffsets builds a Span from byte offsets.
func (s *Source) SpanFromOffsets(start, end int) Span {
	return Span{
		File:  s.name,
		Start: s.line.Position(start),
		End:   s.line.Position(end),
	}
}