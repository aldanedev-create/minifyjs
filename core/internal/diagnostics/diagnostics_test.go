package diagnostics

import (
	"testing"
)

func TestSeverityString(t *testing.T) {
	cases := map[Severity]string{
		SeverityInfo:    "info",
		SeverityWarning: "warning",
		SeverityError:   "error",
		Severity(99):    "unknown",
	}
	for sev, want := range cases {
		if got := sev.String(); got != want {
			t.Errorf("Severity(%d).String() = %q, want %q", sev, got, want)
		}
	}
}

func TestDiagnosticIsError(t *testing.T) {
	if !NewError("x", "m").IsError() {
		t.Fatal("NewError should be IsError")
	}
	if NewWarning("x", "m").IsError() {
		t.Fatal("NewWarning should not be IsError")
	}
	if NewInfo("x", "m").IsError() {
		t.Fatal("NewInfo should not be IsError")
	}
}

func TestWithLocation(t *testing.T) {
	d := NewError("x", "m").WithLocation("a.js", 2, 3)
	if !d.HasLocation() || d.File != "a.js" || d.Line != 2 || d.Col != 3 {
		t.Fatalf("WithLocation wrong: %+v", d)
	}
}

func TestSortStableOrdersByFileLineCol(t *testing.T) {
	ds := []Diagnostic{
		NewError("b", "b").WithLocation("b.js", 1, 1),
		NewError("a", "a").WithLocation("a.js", 2, 1),
		NewError("c", "c").WithLocation("a.js", 1, 5),
		NewError("d", "d").WithLocation("a.js", 1, 1),
	}
	SortStable(ds)
	wantOrder := []string{"d", "c", "a", "b"}
	for i, w := range wantOrder {
		if ds[i].Code != w {
			t.Fatalf("index %d: got %s, want %s (full: %v)", i, ds[i].Code, w, ds)
		}
	}
}