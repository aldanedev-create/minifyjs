package source

// LineMap records where each line of a source file starts (as a byte
// offset). It lets Position and Offset be translated to each other in
// O(log n) time via binary search, which is what diagnostics and
// source maps need.
//
// LineMap is immutable after construction. Callers create it once per
// source file with NewLineMap and reuse it.
type LineMap struct {
	// offsets[i] is the byte offset of the first byte of line i+1.
	// offsets[0] is always 0. len(offsets) is the number of lines.
	offsets []int
	// length is the total byte length of the source, cached so the
	// last line can be bounded without re-reading the source.
	length int
}

// NewLineMap scans src and records the starting offset of every line.
// A line starts at offset 0 and after every '\n'. This matches how
// editors and most toolchains count lines.
func NewLineMap(src string) *LineMap {
	offsets := make([]int, 0, 64)
	offsets = append(offsets, 0)
	for i := 0; i < len(src); i++ {
		if src[i] == '\n' {
			offsets = append(offsets, i+1)
		}
	}
	return &LineMap{offsets: offsets, length: len(src)}
}

// LineCount returns the number of lines in the source.
func (m *LineMap) LineCount() int {
	return len(m.offsets)
}

// Position converts a byte offset into a Position. Offsets past the
// end of the source are clamped to the end.
func (m *LineMap) Position(offset int) Position {
	if offset < 0 {
		offset = 0
	}
	if offset > m.length {
		offset = m.length
	}
	// Binary search for the greatest i such that offsets[i] <= offset.
	lo, hi := 0, len(m.offsets)-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if m.offsets[mid] <= offset {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	lineStart := m.offsets[lo]
	return Position{
		Offset: offset,
		Line:   lo + 1,
		Col:    offset - lineStart + 1,
	}
}

// Offset returns the byte offset of the start of the given 1-based
// line. Returns -1 if line is out of range.
func (m *LineMap) Offset(line int) int {
	if line < 1 || line > len(m.offsets) {
		return -1
	}
	return m.offsets[line-1]
}