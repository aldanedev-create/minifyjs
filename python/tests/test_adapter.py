"""Tests for the generic Adapter base class."""

from __future__ import annotations

from pathlib import Path

import pytest

from minifyjs import Adapter, Options, Result

from .conftest import BINARY_PATH

pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_output_name_default(tmp_path):
    a = Adapter()
    assert a.output_name(tmp_path / "app.js") == tmp_path / "app.min.js"


def test_should_process_js():
    a = Adapter()
    assert a.should_process(Path("app.js"))


def test_should_process_skips_already_minified():
    a = Adapter()
    assert not a.should_process(Path("app.min.js"))


def test_should_process_skips_other_extensions():
    a = Adapter()
    assert not a.should_process(Path("app.css"))


def test_minify_file_writes_output(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("function add(a, b) { return a + b; }\n")
    a = Adapter()
    r = a.minify_file(src)
    assert r.code
    out = tmp_path / "app.min.js"
    assert out.is_file()
    assert "\n" not in out.read_text()


def test_minify_file_explicit_output(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("const x = 1;\n")
    dst = tmp_path / "custom.js"
    a = Adapter()
    a.minify_file(src, output=dst)
    assert dst.is_file()


def test_minify_directory(tmp_path):
    (tmp_path / "a.js").write_text("const a = 1;\n")
    (tmp_path / "b.js").write_text("const b = 2;\n")
    (tmp_path / "ignore.css").write_text("body{}")
    a = Adapter()
    results = a.minify_directory(tmp_path)
    assert len(results) == 2
    assert (tmp_path / "a.min.js").is_file()
    assert (tmp_path / "b.min.js").is_file()
    assert not (tmp_path / "ignore.min.css").exists()


def test_minify_directory_recursive(tmp_path):
    sub = tmp_path / "sub"
    sub.mkdir()
    (sub / "a.js").write_text("const a = 1;\n")
    a = Adapter()
    results = a.minify_directory(tmp_path, recursive=True)
    assert len(results) == 1
    assert (sub / "a.min.js").is_file()


def test_minify_directory_non_recursive(tmp_path):
    sub = tmp_path / "sub"
    sub.mkdir()
    (sub / "a.js").write_text("const a = 1;\n")
    (tmp_path / "b.js").write_text("const b = 2;\n")
    a = Adapter()
    results = a.minify_directory(tmp_path, recursive=False)
    assert len(results) == 1
    assert (tmp_path / "b.min.js").is_file()
    assert not (sub / "a.min.js").exists()


def test_minify_directory_skips_node_modules(tmp_path):
    nm = tmp_path / "node_modules"
    nm.mkdir()
    (nm / "lib.js").write_text("const x = 1;\n")
    (tmp_path / "app.js").write_text("const y = 2;\n")
    a = Adapter()
    results = a.minify_directory(tmp_path)
    assert len(results) == 1


def test_custom_output_suffix(tmp_path):
    class CustomAdapter(Adapter):
        output_suffix = ".compressed"

    a = CustomAdapter()
    assert a.output_name(tmp_path / "app.js") == tmp_path / "app.compressed.js"
    (tmp_path / "app.js").write_text("const x = 1;\n")
    a.minify_file(tmp_path / "app.js")
    assert (tmp_path / "app.compressed.js").is_file()


def test_custom_should_process(tmp_path):
    class OnlyA(Adapter):
        def should_process(self, path):
            return path.name == "a.js"

    (tmp_path / "a.js").write_text("const a = 1;\n")
    (tmp_path / "b.js").write_text("const b = 2;\n")
    a = OnlyA()
    results = a.minify_directory(tmp_path)
    assert len(results) == 1
    assert (tmp_path / "a.min.js").is_file()
    assert not (tmp_path / "b.min.js").exists()


def test_adapter_with_options(tmp_path):
    (tmp_path / "app.js").write_text("const f = (a) => a + 1;\n")
    a = Adapter(options=Options(target="es5"))
    r = a.minify_file(tmp_path / "app.js")
    assert "=>" not in r.code


# NEW tests


def test_adapter_default_options():
    a = Adapter()
    assert a.default_options.compress is True
    assert a.default_options.mangle is True


def test_adapter_custom_default_options():
    opts = Options(compress=False, mangle=False)
    a = Adapter(options=opts)
    assert a.default_options.compress is False


def test_adapter_extension_default():
    a = Adapter()
    assert a.extension == ".js"


def test_adapter_output_suffix_default():
    a = Adapter()
    assert a.output_suffix == ".min"


def test_adapter_should_process_skips_directories(tmp_path):
    a = Adapter()
    assert not a.should_process(tmp_path / "subdir")


def test_adapter_should_process_handles_no_extension():
    a = Adapter()
    assert not a.should_process(Path("Makefile"))


def test_adapter_should_process_multiple_dots():
    a = Adapter()
    assert a.should_process(Path("app.something.js"))
    assert not a.should_process(Path("app.something.min.js"))


def test_adapter_output_name_multiple_dots():
    a = Adapter()
    out = a.output_name(Path("app.something.js"))
    assert out.name == "app.something.min.js"


def test_adapter_minify_file_no_trailing_newline(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("const x = 1;\n")
    a = Adapter()
    a.minify_file(src)
    text = (tmp_path / "app.min.js").read_text()
    assert not text.endswith("\n")


def test_adapter_minify_file_returns_result(tmp_path):
    src = tmp_path / "app.js"
    src.write_text("const x = 1;\n")
    a = Adapter()
    r = a.minify_file(src)
    assert isinstance(r, Result)


def test_adapter_minify_directory_empty(tmp_path):
    a = Adapter()
    results = a.minify_directory(tmp_path)
    assert results == []


def test_adapter_minify_directory_no_matching_files(tmp_path):
    (tmp_path / "a.css").write_text("body {}")
    (tmp_path / "b.txt").write_text("hi")
    a = Adapter()
    results = a.minify_directory(tmp_path)
    assert results == []


def test_adapter_custom_extension(tmp_path):
    class CssAdapter(Adapter):
        extension = ".css"
        output_suffix = ".min"

    (tmp_path / "style.css").write_text("body { margin: 0; }")
    (tmp_path / "app.js").write_text("const x = 1;")
    a = CssAdapter()
    assert a.should_process(tmp_path / "style.css")
    assert not a.should_process(tmp_path / "app.js")


def test_adapter_skips_hidden_directories(tmp_path):
    hidden = tmp_path / ".git"
    hidden.mkdir()
    (hidden / "config.js").write_text("const x = 1;")
    (tmp_path / "app.js").write_text("const x = 1;")
    a = Adapter()
    a.minify_directory(tmp_path)
    assert not (hidden / "config.min.js").exists()


def test_adapter_skips_pycache(tmp_path):
    cache = tmp_path / "__pycache__"
    cache.mkdir()
    (cache / "something.js").write_text("const x = 1;")
    (tmp_path / "app.js").write_text("const x = 1;")
    a = Adapter()
    a.minify_directory(tmp_path)
    assert not (cache / "something.min.js").exists()


def test_adapter_skips_node_modules_recursively(tmp_path):
    nested = tmp_path / "src" / "node_modules" / "pkg"
    nested.mkdir(parents=True)
    (nested / "lib.js").write_text("const x = 1;")
    (tmp_path / "src" / "app.js").write_text("const x = 1;")
    a = Adapter()
    results = a.minify_directory(tmp_path)
    assert len(results) == 1


def test_adapter_multiple_files_same_directory(tmp_path):
    for name in ("a.js", "b.js", "c.js"):
        (tmp_path / name).write_text("const x = 1;")
    a = Adapter()
    results = a.minify_directory(tmp_path)
    assert len(results) == 3


def test_adapter_result_order_is_stable(tmp_path):
    for name in ("a.js", "b.js", "c.js"):
        (tmp_path / name).write_text("const x = 1;")
    a = Adapter()
    results1 = a.minify_directory(tmp_path)
    results2 = a.minify_directory(tmp_path)
    assert [r.code for r in results1] == [r.code for r in results2]
