"""Tests that reading JavaScript from a file works end-to-end."""

from __future__ import annotations

import pytest

from minifyjs import Adapter, minify
from minifyjs.errors import MinifyError

from .conftest import BINARY_PATH


pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_read_file_and_minify(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("function add(a, b) { return a + b; }\n", encoding="utf-8")
    r = minify(src.read_text(encoding="utf-8"))
    assert not r.has_errors


def test_adapter_reads_file(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("const x = 1;\n", encoding="utf-8")
    a = Adapter()
    r = a.minify_file(src)
    assert not r.has_errors


def test_adapter_writes_minified_file(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("function add(a, b) { return a + b; }\n", encoding="utf-8")
    a = Adapter()
    a.minify_file(src)
    out = tmp_path / "app.min.js"
    assert out.is_file()
    assert "\n" not in out.read_text()


def test_adapter_missing_file_raises(tmp_path):
    a = Adapter()
    with pytest.raises(FileNotFoundError):
        a.minify_file(tmp_path / "nope.js")


def test_adapter_reads_unicode_file(tmp_path):
    src = tmp_path / "app.js"
    src.write_text('const s = "日本語";\n', encoding="utf-8")
    a = Adapter()
    r = a.minify_file(src)
    assert not r.has_errors


def test_adapter_reads_file_with_spaces_in_name(tmp_path):
    src = tmp_path / "with spaces.js"
    src.write_text("const x = 1;\n", encoding="utf-8")
    a = Adapter()
    r = a.minify_file(src)
    assert not r.has_errors
    assert (tmp_path / "with spaces.min.js").is_file()


def test_adapter_reads_file_with_unicode_name(tmp_path):
    src = tmp_path / "café.js"
    src.write_text("const x = 1;\n", encoding="utf-8")
    a = Adapter()
    r = a.minify_file(src)
    assert not r.has_errors


def test_read_empty_file(tmp_path):
    src = tmp_path / "empty.js"
    src.write_text("", encoding="utf-8")
    a = Adapter()
    r = a.minify_file(src)
    assert r.code == ""


def test_read_file_with_syntax_error_raises(tmp_path):
    src = tmp_path / "bad.js"
    src.write_text("function () { } }", encoding="utf-8")
    a = Adapter()
    with pytest.raises(MinifyError):
        a.minify_file(src)


def test_read_file_with_only_comments(tmp_path):
    src = tmp_path / "comments.js"
    src.write_text("// just a comment\n/* and a block */\n", encoding="utf-8")
    a = Adapter()
    r = a.minify_file(src)
    assert r.code == ""


def test_read_file_with_crlf(tmp_path):
    src = tmp_path / "crlf.js"
    src.write_bytes(b"const a = 1;\r\nconst b = 2;\r\n")
    a = Adapter()
    r = a.minify_file(src)
    assert "\r" not in r.code
    assert "\n" not in r.code


def test_read_file_with_bom(tmp_path):
    src = tmp_path / "bom.js"
    src.write_bytes(b"\xef\xbb\xbfconst x = 1;\n")
    a = Adapter()
    r = a.minify_file(src)
    assert not r.has_errors


def test_relative_path_via_adapter(tmp_path, monkeypatch):
    src = tmp_path / "app.js"
    src.write_text("const x = 1;\n", encoding="utf-8")
    monkeypatch.chdir(tmp_path)
    a = Adapter()
    r = a.minify_file("app.js")
    assert not r.has_errors


def test_adapter_directory_with_nested(tmp_path):
    sub = tmp_path / "nested"
    sub.mkdir()
    (sub / "a.js").write_text("const a = 1;\n", encoding="utf-8")
    (tmp_path / "b.js").write_text("const b = 2;\n", encoding="utf-8")
    a = Adapter()
    results = a.minify_directory(tmp_path)
    assert len(results) == 2
    assert (sub / "a.min.js").is_file()
    assert (tmp_path / "b.min.js").is_file()