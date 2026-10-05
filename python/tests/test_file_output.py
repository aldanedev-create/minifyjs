"""Tests that writing minified output to a file works end-to-end."""

from __future__ import annotations

import pytest

from minifyjs import Adapter, Options

from .conftest import BINARY_PATH

pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_write_minified_file(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("function add(a, b) { return a + b; }\n", encoding="utf-8")
    a = Adapter()
    a.minify_file(src)
    out = tmp_path / "app.min.js"
    assert out.is_file()
    assert "function" in out.read_text()


def test_write_to_explicit_path(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("const x = 1;\n", encoding="utf-8")
    dst = tmp_path / "custom.js"
    a = Adapter()
    a.minify_file(src, output=dst)
    assert dst.is_file()


def test_overwrite_existing_output(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("const x = 1;\n", encoding="utf-8")
    out = tmp_path / "app.min.js"
    out.write_text("// old contents\n", encoding="utf-8")
    a = Adapter()
    a.minify_file(src)
    assert "old contents" not in out.read_text()


def test_output_is_valid_utf8(tmp_path):
    src = tmp_path / "app.js"
    src.write_text('const s = "日本語";\n', encoding="utf-8")
    a = Adapter()
    a.minify_file(src)
    out = tmp_path / "app.min.js"
    # read_text with utf-8 will raise if the bytes are not valid utf-8.
    text = out.read_text(encoding="utf-8")
    assert "日本語" in text


def test_output_no_trailing_newline(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("const x = 1;\n", encoding="utf-8")
    a = Adapter()
    a.minify_file(src)
    out = tmp_path / "app.min.js"
    text = out.read_text(encoding="utf-8")
    assert not text.endswith("\n")


def test_output_single_line(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("function f() {\n    return 1;\n}\n", encoding="utf-8")
    a = Adapter()
    a.minify_file(src)
    text = (tmp_path / "app.min.js").read_text(encoding="utf-8")
    assert "\n" not in text


def test_output_bytes_match_result(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("function add(a, b) { return a + b; }\n", encoding="utf-8")
    a = Adapter()
    r = a.minify_file(src)
    out = tmp_path / "app.min.js"
    assert out.stat().st_size == r.minified_bytes


def test_output_relative_path(tmp_path, monkeypatch):
    src = tmp_path / "app.js"
    src.write_text("const x = 1;\n", encoding="utf-8")
    monkeypatch.chdir(tmp_path)
    a = Adapter()
    a.minify_file("app.js", output="out.min.js")
    assert (tmp_path / "out.min.js").is_file()


def test_output_with_target_es5(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("const f = (a) => a + 1;\n", encoding="utf-8")
    a = Adapter(options=Options(compress=True, mangle=True, target="es5"))
    a.minify_file(src)
    text = (tmp_path / "app.min.js").read_text(encoding="utf-8")
    assert "=>" not in text


def test_output_directory_creates_files(tmp_path):
    (tmp_path / "a.js").write_text("const a = 1;\n", encoding="utf-8")
    (tmp_path / "b.js").write_text("const b = 2;\n", encoding="utf-8")
    (tmp_path / "c.js").write_text("const c = 3;\n", encoding="utf-8")
    a = Adapter()
    results = a.minify_directory(tmp_path)
    assert len(results) == 3
    for name in ("a.min.js", "b.min.js", "c.min.js"):
        assert (tmp_path / name).is_file()


def test_output_directory_does_not_rewrite_already_minified(tmp_path):
    # A .min.js file should not be re-minified into .min.min.js.
    (tmp_path / "app.js").write_text("const x = 1;\n", encoding="utf-8")
    (tmp_path / "app.min.js").write_text("const x=1;", encoding="utf-8")
    a = Adapter()
    results = a.minify_directory(tmp_path)
    # Only app.js should be processed.
    assert len(results) == 1
    assert not (tmp_path / "app.min.min.js").exists()


def test_output_preserves_source_when_output_is_input(tmp_path):
    # Writing to the same path as the input is legal.
    src = tmp_path / "app.js"
    src.write_text("function add(a, b) {\n    return a + b;\n}\n", encoding="utf-8")
    a = Adapter()
    a.minify_file(src, output=src)
    text = src.read_text(encoding="utf-8")
    assert "\n" not in text
    assert "function" in text


def test_output_directory_nested(tmp_path):
    sub = tmp_path / "a" / "b"
    sub.mkdir(parents=True)
    (sub / "app.js").write_text("const x = 1;\n", encoding="utf-8")
    a = Adapter()
    a.minify_directory(tmp_path)
    assert (sub / "app.min.js").is_file()


def test_output_directory_skips_non_js(tmp_path):
    (tmp_path / "app.js").write_text("const x = 1;\n", encoding="utf-8")
    (tmp_path / "style.css").write_text("body {}\n", encoding="utf-8")
    (tmp_path / "data.json").write_text("{}\n", encoding="utf-8")
    a = Adapter()
    results = a.minify_directory(tmp_path)
    assert len(results) == 1
    assert not (tmp_path / "style.min.css").exists()
    assert not (tmp_path / "data.min.json").exists()


def test_output_directory_recursive_false(tmp_path):
    sub = tmp_path / "sub"
    sub.mkdir()
    (sub / "a.js").write_text("const a = 1;\n", encoding="utf-8")
    (tmp_path / "b.js").write_text("const b = 2;\n", encoding="utf-8")
    a = Adapter()
    results = a.minify_directory(tmp_path, recursive=False)
    assert len(results) == 1
    assert (tmp_path / "b.min.js").is_file()
    assert not (sub / "a.min.js").exists()


def test_output_directory_with_custom_suffix(tmp_path):
    class CustomAdapter(Adapter):
        output_suffix = ".compressed"

    (tmp_path / "app.js").write_text("const x = 1;\n", encoding="utf-8")
    a = CustomAdapter()
    a.minify_directory(tmp_path)
    assert (tmp_path / "app.compressed.js").is_file()
    assert not (tmp_path / "app.min.js").exists()
