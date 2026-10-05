"""Tests that malformed JavaScript is reported cleanly."""

from __future__ import annotations

import pytest

from minifyjs import MinifyError, minify

from .conftest import BINARY_PATH


pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_unbalanced_brace():
    with pytest.raises(MinifyError):
        minify("function f() { } }")


def test_unterminated_string():
    with pytest.raises(MinifyError):
        minify('const s = "unterminated')


def test_unterminated_template():
    with pytest.raises(MinifyError):
        minify("const s = `unterminated")


def test_unterminated_comment():
    with pytest.raises(MinifyError):
        minify("/* unterminated")


def test_reserved_word_as_identifier():
    with pytest.raises(MinifyError):
        minify("const if = 1;")


def test_duplicate_let():
    with pytest.raises(MinifyError):
        minify("let x = 1; let x = 2;")


def test_missing_value():
    with pytest.raises(MinifyError):
        minify("const x = ;")


# NEW tests


def test_invalid_empty_function_params():
    with pytest.raises(MinifyError):
        minify("function () { }")


def test_invalid_unclosed_paren():
    with pytest.raises(MinifyError):
        minify("f(1, 2")


def test_invalid_unclosed_bracket():
    with pytest.raises(MinifyError):
        minify("const a = [1, 2")


def test_invalid_unclosed_object():
    with pytest.raises(MinifyError):
        minify("const a = { x: 1")


def test_invalid_invalid_regex():
    with pytest.raises(MinifyError):
        minify("const re = /[unclosed/;")


def test_invalid_const_without_value():
    with pytest.raises(MinifyError):
        minify("const x;")


def test_invalid_let_without_value_is_ok():
    # let x; is legal in modern JS (defines undefined). Only const
    # requires a value.
    r = minify("let x;")
    assert not r.has_errors


def test_invalid_duplicate_const():
    with pytest.raises(MinifyError):
        minify("const x = 1; const x = 2;")


def test_invalid_import_missing_from():
    with pytest.raises(MinifyError):
        minify('import { x } "mod";')


def test_invalid_export_default_without_value():
    with pytest.raises(MinifyError):
        minify("export default;")


def test_invalid_arrow_missing_body():
    with pytest.raises(MinifyError):
        minify("const f = () =>;")


def test_invalid_class_missing_body():
    with pytest.raises(MinifyError):
        minify("class A")