"""Tests for the Flask adapter example."""

from __future__ import annotations

from pathlib import Path

import pytest

from examples.flask_adapter import FlaskMinifyAdapter


pytestmark = pytest.mark.binary


def _write(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content)


def test_flask_output_next_to_source(tmp_path):
    static = tmp_path / "static"
    _write(static / "app.js", "const x = 1;\n")

    adapter = FlaskMinifyAdapter(static_folder=static)
    adapter.minify_all()

    assert (static / "app.min.js").is_file()


def test_flask_recursive(tmp_path):
    static = tmp_path / "static"
    _write(static / "a.js", "const a = 1;\n")
    _write(static / "sub" / "b.js", "const b = 2;\n")

    adapter = FlaskMinifyAdapter(static_folder=static)
    adapter.minify_all()

    assert (static / "a.min.js").is_file()
    assert (static / "sub" / "b.min.js").is_file()


def test_flask_skips_already_minified(tmp_path):
    static = tmp_path / "static"
    _write(static / "app.js", "const x = 1;\n")
    _write(static / "app.min.js", "const x=1;")

    adapter = FlaskMinifyAdapter(static_folder=static)
    adapter.minify_all()

    assert (static / "app.min.js").is_file()
    assert not (static / "app.min.min.js").exists()


def test_flask_empty_dir(tmp_path):
    static = tmp_path / "static"
    static.mkdir()

    adapter = FlaskMinifyAdapter(static_folder=static)
    results = adapter.minify_all()
    assert results == []


def test_flask_skips_non_js(tmp_path):
    static = tmp_path / "static"
    _write(static / "app.js", "const x = 1;\n")
    _write(static / "style.css", "body {}")

    adapter = FlaskMinifyAdapter(static_folder=static)
    adapter.minify_all()

    assert (static / "app.min.js").is_file()
    assert not (static / "style.min.css").exists()


def test_flask_minifies_content(tmp_path):
    static = tmp_path / "static"
    _write(static / "app.js", "function f() {\n    return 1;\n}\n")

    adapter = FlaskMinifyAdapter(static_folder=static)
    adapter.minify_all()

    body = (static / "app.min.js").read_text()
    assert "\n" not in body