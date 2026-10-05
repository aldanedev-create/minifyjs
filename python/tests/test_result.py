"""Tests for Result dataclass and its derived properties."""

from __future__ import annotations

import dataclasses

from minifyjs import Diagnostic, Result


def test_result_defaults():
    r = Result(code="const x=1;")
    assert r.code == "const x=1;"
    assert r.map == ""
    assert r.diagnostics == []


def test_result_ratio():
    r = Result(code="x", original_bytes=100, minified_bytes=25)
    assert r.ratio == 0.25


def test_result_ratio_empty():
    r = Result(code="", original_bytes=0, minified_bytes=0)
    assert r.ratio == 0.0


def test_result_bytes_saved():
    r = Result(code="x", original_bytes=100, minified_bytes=25)
    assert r.bytes_saved == 75


def test_result_bytes_saved_negative():
    r = Result(code="x", original_bytes=10, minified_bytes=20)
    assert r.bytes_saved == -10


def test_result_has_errors_false():
    r = Result(code="")
    assert r.has_errors is False


def test_result_has_errors_true():
    r = Result(code="", diagnostics=[
        Diagnostic(severity="error", message="boom"),
    ])
    assert r.has_errors is True


def test_result_has_warnings():
    r = Result(code="", diagnostics=[
        Diagnostic(severity="warning", message="hmm"),
    ])
    assert r.has_warnings is True
    assert r.has_errors is False


# NEW tests


def test_result_is_dataclass():
    assert dataclasses.is_dataclass(Result)


def test_result_map_default_empty():
    r = Result(code="x")
    assert r.map == ""


def test_result_diagnostics_default_empty():
    r = Result(code="x")
    assert r.diagnostics == []


def test_result_diagnostics_are_independent_lists():
    r1 = Result(code="x")
    r2 = Result(code="y")
    r1.diagnostics.append(Diagnostic(severity="error", message="a"))
    assert r2.diagnostics == []


def test_result_ratio_of_full_size():
    r = Result(code="x", original_bytes=100, minified_bytes=100)
    assert r.ratio == 1.0


def test_result_ratio_zero_when_empty():
    r = Result(code="", original_bytes=0, minified_bytes=0)
    assert r.ratio == 0.0


def test_result_bytes_saved_zero():
    r = Result(code="x", original_bytes=100, minified_bytes=100)
    assert r.bytes_saved == 0


def test_result_bytes_saved_when_grew():
    r = Result(code="x", original_bytes=10, minified_bytes=30)
    assert r.bytes_saved == -20


def test_result_has_errors_with_multiple_diagnostics():
    r = Result(code="", diagnostics=[
        Diagnostic(severity="warning", message="w"),
        Diagnostic(severity="error", message="e"),
        Diagnostic(severity="info", message="i"),
    ])
    assert r.has_errors is True


def test_result_has_warnings_only_warnings():
    r = Result(code="", diagnostics=[
        Diagnostic(severity="warning", message="w"),
        Diagnostic(severity="info", message="i"),
    ])
    assert r.has_warnings is True
    assert r.has_errors is False


def test_result_has_warnings_false_when_error():
    r = Result(code="", diagnostics=[
        Diagnostic(severity="error", message="e"),
    ])
    assert r.has_warnings is False


def test_result_repr():
    assert repr(Result(code="x"))


def test_result_equality():
    assert Result(code="x") == Result(code="x")
    assert Result(code="x") != Result(code="y")


def test_result_field_types():
    r = Result(code="x", map="y", original_bytes=1, minified_bytes=1)
    assert isinstance(r.code, str)
    assert isinstance(r.map, str)
    assert isinstance(r.original_bytes, int)
    assert isinstance(r.minified_bytes, int)