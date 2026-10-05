"""Tests for Diagnostic parsing and formatting."""

from __future__ import annotations

import dataclasses

from minifyjs import Diagnostic


def test_diagnostic_defaults():
    d = Diagnostic(severity="error", message="boom")
    assert d.code == ""
    assert d.file == ""
    assert d.line == 0
    assert d.col == 0


def test_diagnostic_is_error():
    assert Diagnostic(severity="error", message="x").is_error
    assert not Diagnostic(severity="warning", message="x").is_error
    assert not Diagnostic(severity="info", message="x").is_error


def test_diagnostic_has_location():
    d = Diagnostic(severity="error", message="x", line=3, col=5)
    assert d.has_location
    d2 = Diagnostic(severity="error", message="x")
    assert not d2.has_location


def test_diagnostic_str_with_location():
    d = Diagnostic(
        severity="error", message="bad",
        file="app.js", line=3, col=5,
    )
    assert str(d) == "app.js:3:5: error: bad"


def test_diagnostic_str_without_file():
    d = Diagnostic(severity="warning", message="hmm", line=2, col=1)
    assert str(d) == "2:1: warning: hmm"


def test_diagnostic_str_without_location():
    d = Diagnostic(severity="info", message="fyi")
    assert str(d) == "info: fyi"


# NEW tests


def test_diagnostic_is_dataclass():
    assert dataclasses.is_dataclass(Diagnostic)


def test_diagnostic_default_code_empty():
    d = Diagnostic(severity="error", message="x")
    assert d.code == ""


def test_diagnostic_str_error_with_location():
    d = Diagnostic(severity="error", message="bad", file="a.js", line=1, col=1)
    assert str(d) == "a.js:1:1: error: bad"


def test_diagnostic_str_warning_with_location():
    d = Diagnostic(severity="warning", message="hmm", file="b.js", line=3, col=5)
    assert str(d) == "b.js:3:5: warning: hmm"


def test_diagnostic_str_info_no_location():
    d = Diagnostic(severity="info", message="fyi")
    assert str(d) == "info: fyi"


def test_diagnostic_with_file_but_no_line():
    d = Diagnostic(severity="error", message="bad", file="a.js")
    assert not d.has_location
    assert str(d) == "error: bad"


def test_diagnostic_severity_values():
    for sev in ("info", "warning", "error"):
        d = Diagnostic(severity=sev, message="m")
        assert d.severity == sev


def test_diagnostic_is_error_other_values():
    assert not Diagnostic(severity="debug", message="x").is_error
    assert not Diagnostic(severity="", message="x").is_error


def test_diagnostic_equality():
    assert Diagnostic(severity="error", message="x") == \
        Diagnostic(severity="error", message="x")


def test_diagnostic_inequality_on_message():
    assert Diagnostic(severity="error", message="a") != \
        Diagnostic(severity="error", message="b")


def test_diagnostic_inequality_on_severity():
    assert Diagnostic(severity="error", message="x") != \
        Diagnostic(severity="warning", message="x")


def test_diagnostic_repr():
    assert repr(Diagnostic(severity="error", message="x"))


def test_diagnostic_zero_values():
    d = Diagnostic(severity="error", message="x")
    assert d.file == ""
    assert d.line == 0
    assert d.col == 0
    assert d.code == ""


def test_diagnostic_string_with_col_zero():
    d = Diagnostic(severity="error", message="x", file="a.js", line=3, col=0)
    assert not d.has_location