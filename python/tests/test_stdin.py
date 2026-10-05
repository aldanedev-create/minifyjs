"""Tests for stdin handling through the Python bridge."""

from __future__ import annotations

import pytest

from minifyjs import minify

from .conftest import BINARY_PATH


pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_stdin_simple():
    r = minify("function add(a, b) { return a + b; }")
    assert not r.has_errors
    assert "function" in r.code


def test_stdin_with_newlines():
    src = "function add(a, b) {\n    return a + b;\n}\n"
    r = minify(src)
    assert "\n" not in r.code


def test_stdin_empty():
    r = minify("")
    assert r.code == ""
    assert not r.has_errors


def test_stdin_whitespace_only():
    r = minify("   \n\t  ")
    assert r.code == ""


def test_stdin_comments_only():
    r = minify("// just a comment\n/* and a block */\n")
    assert r.code == ""


def test_stdin_unicode():
    r = minify("const café = 1;")
    assert not r.has_errors


def test_stdin_unicode_string():
    r = minify('const s = "日本語";')
    assert not r.has_errors


def test_stdin_emoji():
    r = minify('const s = "🎉";')
    assert not r.has_errors


def test_stdin_crlf():
    r = minify("function add(a, b) {\r\n    return a + b;\r\n}\r\n")
    assert "\r" not in r.code
    assert "\n" not in r.code


def test_stdin_bom():
    r = minify("\ufeffconst x = 1;")
    assert not r.has_errors


def test_stdin_large_input():
    src = "\n".join(f"const v{i} = {i};" for i in range(1000))
    r = minify(src)
    assert not r.has_errors


def test_stdin_long_single_line():
    src = "const x = " + "+".join(["1"] * 5000) + ";"
    r = minify(src)
    assert not r.has_errors


def test_stdin_syntax_error_raises():
    from minifyjs import MinifyError
    with pytest.raises(MinifyError):
        minify("function () { } }")


def test_stdin_preserves_shebang():
    r = minify("#!/usr/bin/env node\nconst x = 1;")
    assert "#!/usr/bin/env node" in r.code


def test_stdin_returns_original_bytes():
    src = "const x = 1;"
    r = minify(src)
    assert r.original_bytes == len(src.encode("utf-8"))


def test_stdin_returns_minified_bytes():
    src = "const x = 1;"
    r = minify(src)
    assert r.minified_bytes == len(r.code.encode("utf-8"))