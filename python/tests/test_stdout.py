"""Tests for stdout handling through the Python bridge."""

from __future__ import annotations

import pytest

from minifyjs import minify

from .conftest import BINARY_PATH


pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_stdout_exact_bytes():
    r = minify("const x = 1;")
    assert r.code == "const x=1;"


def test_stdout_no_trailing_newline():
    r = minify("const x = 1;")
    assert not r.code.endswith("\n")


def test_stdout_no_leading_newline():
    r = minify("const x = 1;")
    assert not r.code.startswith("\n")


def test_stdout_ends_with_semicolon():
    r = minify("const x = 1;")
    assert r.code.endswith(";")


def test_stdout_unicode_preserved():
    r = minify('const s = "café";', compress=False)
    assert "café" in r.code


def test_stdout_emoji_preserved():
    r = minify('const s = "🎉";', compress=False)
    assert "🎉" in r.code


def test_stdout_no_bom():
    r = minify("const x = 1;")
    assert not r.code.startswith("\ufeff")


def test_stdout_empty_on_empty_input():
    r = minify("")
    assert r.code == ""


def test_stdout_no_minify_flag():
    src = "const x = 1;\n"
    r = minify(src, compress=False, mangle=False)
    assert "const x" in r.code


def test_stdout_idempotent():
    r1 = minify("function f() { return 1; }")
    r2 = minify(r1.code)
    assert r1.code == r2.code


def test_stdout_exact_when_target_preserves():
    r = minify("const x = 1;", target="esnext")
    assert r.code == "const x=1;"


def test_stdout_line_count_one():
    r = minify("function add(a, b) {\n    return a + b;\n}\n")
    assert r.code.count("\n") == 0


def test_stdout_bytes_match_code():
    r = minify("const x = 1;")
    assert r.minified_bytes == len(r.code.encode("utf-8"))


def test_stdout_ratio_matches_lengths():
    r = minify("function add(a, b) { return a + b; }")
    expected = r.minified_bytes / r.original_bytes
    assert abs(r.ratio - expected) < 1e-9


def test_stdout_bytes_saved_matches():
    r = minify("function add(a, b) { return a + b; }")
    assert r.bytes_saved == r.original_bytes - r.minified_bytes