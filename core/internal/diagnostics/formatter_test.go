package diagnostics

import (
	"bytes"
	"strings"
	"testing"
)

func TestTerminalFormatterWithLocation(t *testing.T) {
	var buf bytes.Buffer
	f := TerminalFormatter{Prefix: "minifyjs: "}
	err := f.Format(&buf, []Diagnostic{
		NewError("lex.bad", "unexpected token").WithLocation("app.js", 3, 5),
	})
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	want := "minifyjs: app.js:3:5: error: unexpected token\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestTerminalFormatterWithoutLocation(t *testing.T) {
	var buf bytes.Buffer
	f := TerminalFormatter{}
	_ = f.Format(&buf, []Diagnostic{NewWarning("x", "no location")})
	if !strings.Contains(buf.String(), "warning: no location") {
		t.Fatalf("got %q", buf.String())
	}
}

func TestJSONFormatterEmitsArray(t *testing.T) {
	var buf bytes.Buffer
	f := JSONFormatter{}
	if err := f.Format(&buf, []Diagnostic{NewError("x", "m").WithLocation("a.js", 1, 2)}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{`"severity":"error"`, `"code":"x"`, `"file":"a.js"`, `"line":1`, `"col":2`} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %s: %s", want, out)
		}
	}
}

func TestJSONFormatterEmpty(t *testing.T) {
	var buf bytes.Buffer
	if err := (JSONFormatter{}).Format(&buf, nil); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(buf.String()) != "[]" {
		t.Fatalf("got %q, want []", buf.String())
	}
}

func TestHasErrors(t *testing.T) {
	if HasErrors(nil) {
		t.Fatal("nil should have no errors")
	}
	if HasErrors([]Diagnostic{NewWarning("x", "m")}) {
		t.Fatal("only warnings -> false")
	}
	if !HasErrors([]Diagnostic{NewError("x", "m")}) {
		t.Fatal("error -> true")
	}
}