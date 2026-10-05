"""Tests for the subprocess bridge."""

from __future__ import annotations

import pytest

from minifyjs import MinifyError
from minifyjs._runner import _parse_diagnostics, _parse_one
from minifyjs.diagnostics import Diagnostic

from .conftest import BINARY_PATH


pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_parse_one_file_line_col():
    d = _parse_one("app.js:3:5: error: bad token")
    assert d is not None
    assert d.file == "app.js"
    assert d.line == 3
    assert d.col == 5
    assert d.severity == "error"
    assert d.message == "bad token"


def test_parse_one_no_location():
    d = _parse_one("warning: something")
    assert d is not None
    assert d.severity == "warning"
    assert d.message == "something"
    assert d.file == ""


def test_parse_one_rejects_garbage():
    assert _parse_one("this is not a diagnostic") is None


def test_parse_diagnostics_multiline():
    stderr = (
        "minifyjs: app.js:1:1: error: bad\n"
        "minifyjs: warning: something\n"
        "not a diagnostic line\n"
    )
    diags = _parse_diagnostics(stderr)
    assert len(diags) == 2
    assert diags[0].severity == "error"
    assert diags[1].severity == "warning"


def test_parse_diagnostics_empty():
    assert _parse_diagnostics("") == []


def test_run_minify_simple():
    from minifyjs import minify
    r = minify("const x = 1;")
    assert "const x" in r.code
    assert not r.has_errors


def test_run_minify_syntax_error_raises():
    from minifyjs import minify
    with pytest.raises(MinifyError):
        minify("function () { } }")


# NEW tests


def test_parse_one_error_with_empty_message():
    d = _parse_one("app.js:1:1: error: ")
    assert d is not None
    assert d.severity == "error"
    assert d.message == ""


def test_parse_one_warning_without_location():
    d = _parse_one("warning: deprecated")
    assert d is not None
    assert d.severity == "warning"
    assert d.message == "deprecated"


def test_parse_one_info_without_location():
    d = _parse_one("info: note")
    assert d is not None
    assert d.severity == "info"


def test_parse_one_rejects_unknown_severity():
    assert _parse_one("bogus: not a severity") is None


def test_parse_one_rejects_non_numeric_position():
    assert _parse_one("app.js:x:y: error: msg") is None


def test_parse_diagnostics_ignores_blank_lines():
    stderr = "\n\nminifyjs: warning: something\n\n"
    diags = _parse_diagnostics(stderr)
    assert len(diags) == 1


def test_parse_diagnostics_ignores_lines_without_prefix():
    stderr = "some random output\nmore output\n"
    assert _parse_diagnostics(stderr) == []


def test_parse_diagnostics_multiple_errors():
    stderr = (
        "minifyjs: a.js:1:1: error: first\n"
        "minifyjs: a.js:2:1: error: second\n"
    )
    diags = _parse_diagnostics(stderr)
    assert len(diags) == 2
    assert diags[0].message == "first"
    assert diags[1].message == "second"