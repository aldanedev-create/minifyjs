package source

import "testing"

func TestNewStoresNameAndText(t *testing.T) {
	s := New("app.js", "let x = 1;\n")
	if s.Name() != "app.js" {
		t.Fatalf("Name = %q", s.Name())
	}
	if s.Text() != "let x = 1;\n" {
		t.Fatalf("Text = %q", s.Text())
	}
	if s.Len() != 11 {
		t.Fatalf("Len = %d, want 11", s.Len())
	}
}

func TestPositionAt(t *testing.T) {
	s := New("t.js", "ab\ncd\n")
	cases := map[int]Position{
		0: {Offset: 0, Line: 1, Col: 1},
		1: {Offset: 1, Line: 1, Col: 2},
		3: {Offset: 3, Line: 2, Col: 1},
		4: {Offset: 4, Line: 2, Col: 2},
	}
	for off, want := range cases {
		got := s.PositionAt(off)
		if got != want {
			t.Errorf("PositionAt(%d) = %+v, want %+v", off, got, want)
		}
	}
}

func TestSpanFromOffsets(t *testing.T) {
	s := New("t.js", "hello world")
	sp := s.SpanFromOffsets(6, 11)
	if sp.File != "t.js" {
		t.Fatalf("File = %q", sp.File)
	}
	if sp.Start.Col != 7 || sp.End.Col != 12 {
		t.Fatalf("span = %+v", sp)
	}
	if sp.Len() != 5 {
		t.Fatalf("Len = %d", sp.Len())
	}
}

func TestPositionIsValid(t *testing.T) {
	if (Position{}).IsValid() {
		t.Fatal("zero Position should be invalid")
	}
	if !(Position{Line: 1, Col: 1}).IsValid() {
		t.Fatal("1:1 should be valid")
	}
}

func TestPositionString(t *testing.T) {
	p := Position{Line: 3, Col: 5}
	if p.String() != "3:5" {
		t.Fatalf("String = %q", p.String())
	}
	if (Position{}).String() != "?:?" {
		t.Fatalf("zero String = %q", (Position{}).String())
	}
}

func TestLocationString(t *testing.T) {
	l := Location{File: "a.js", Position: Position{Line: 1, Col: 2}}
	if l.String() != "a.js:1:2" {
		t.Fatalf("String = %q", l.String())
	}
}