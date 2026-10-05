package source

import "testing"

func TestLineMapBasic(t *testing.T) {
	m := NewLineMap("a\nbb\nccc")
	if m.LineCount() != 3 {
		t.Fatalf("LineCount = %d, want 3", m.LineCount())
	}
	cases := []struct {
		off  int
		line int
		col  int
	}{
		{0, 1, 1},
		{1, 1, 2},
		{2, 2, 1},
		{3, 2, 2},
		{4, 2, 3},
		{5, 3, 1},
		{8, 3, 4},
	}
	for _, c := range cases {
		p := m.Position(c.off)
		if p.Line != c.line || p.Col != c.col {
			t.Errorf("Position(%d) = line %d col %d, want line %d col %d",
				c.off, p.Line, p.Col, c.line, c.col)
		}
	}
}

func TestLineMapOffset(t *testing.T) {
	m := NewLineMap("a\nbb\nccc")
	if m.Offset(1) != 0 {
		t.Fatalf("Offset(1) = %d", m.Offset(1))
	}
	if m.Offset(2) != 2 {
		t.Fatalf("Offset(2) = %d", m.Offset(2))
	}
	if m.Offset(3) != 5 {
		t.Fatalf("Offset(3) = %d", m.Offset(3))
	}
	if m.Offset(0) != -1 {
		t.Fatalf("Offset(0) should be -1")
	}
	if m.Offset(4) != -1 {
		t.Fatalf("Offset(4) should be -1")
	}
}

func TestLineMapClampsOutOfRange(t *testing.T) {
	m := NewLineMap("abc")
	if p := m.Position(100); p.Offset != 3 {
		t.Fatalf("Position(100).Offset = %d, want 3", p.Offset)
	}
	if p := m.Position(-5); p.Offset != 0 {
		t.Fatalf("Position(-5).Offset = %d, want 0", p.Offset)
	}
}

func TestLineMapEmpty(t *testing.T) {
	m := NewLineMap("")
	if m.LineCount() != 1 {
		t.Fatalf("empty source should report 1 line, got %d", m.LineCount())
	}
	p := m.Position(0)
	if p.Line != 1 || p.Col != 1 {
		t.Fatalf("Position(0) = %+v", p)
	}
}

func TestLineMapTrailingNewline(t *testing.T) {
	m := NewLineMap("a\n")
	if m.LineCount() != 2 {
		t.Fatalf("LineCount = %d, want 2", m.LineCount())
	}
	if m.Offset(2) != 2 {
		t.Fatalf("Offset(2) = %d, want 2", m.Offset(2))
	}
}