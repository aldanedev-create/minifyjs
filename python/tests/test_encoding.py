"""Tests for UTF-8 encoding handling across the bridge."""

from __future__ import annotations

import pytest

from minifyjs import minify

from .conftest import BINARY_PATH

pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_ascii_round_trip():
    r = minify('const s = "hello";')
    assert "hello" in r.code


def test_latin1_accent():
    r = minify('const s = "café";')
    assert "café" in r.code


def test_japanese():
    r = minify('const s = "日本語";')
    assert "日本語" in r.code


def test_korean():
    r = minify('const s = "한국어";')
    assert "한국어" in r.code


def test_emoji():
    r = minify('const s = "🎉";')
    assert "🎉" in r.code


def test_emoji_in_string_not_comment():
    r = minify('const s = "🎉"; // 🎉 comment')
    assert "🎉" in r.code
    assert "comment" not in r.code


def test_unicode_identifier():
    r = minify("const café = 1;")
    assert not r.has_errors


def test_unicode_identifier_japanese():
    r = minify("const 名前 = 1;")
    assert not r.has_errors


def test_unicode_escape_sequence_preserved():
    r = minify('const s = "\\u00e9";')
    assert "\\u00e9" in r.code or "é" in r.code


def test_mixed_ascii_unicode():
    r = minify('const greeting = "Hello 世界";')
    assert "Hello" in r.code
    assert "世界" in r.code


def test_original_bytes_uses_utf8_length():
    # "café" is 4 characters but 5 UTF-8 bytes.
    src = 'const s = "café";'
    r = minify(src)
    assert r.original_bytes == len(src.encode("utf-8"))


def test_minified_bytes_uses_utf8_length():
    r = minify('const s = "café";')
    assert r.minified_bytes == len(r.code.encode("utf-8"))


def test_rtl_text():
    r = minify('const s = "مرحبا";')
    assert not r.has_errors


def test_surrogate_pair_emoji():
    # 🎉 is a single codepoint but two UTF-16 units. Python treats
    # it as one character; the encoding layer must preserve it.
    r = minify('const s = "🎉🎉";')
    assert r.code.count("🎉") == 2
