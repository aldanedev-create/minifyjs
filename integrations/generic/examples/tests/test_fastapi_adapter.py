"""Tests for the FastAPI adapter example."""

from __future__ import annotations

from pathlib import Path

import pytest

from examples.fastapi_adapter import FastAPIMinifyAdapter


pytestmark = pytest.mark.binary


def _write(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content)


def test_fastapi_build_copies_and_minifies(tmp_path):
    src = tmp_path / "static"
    dist = tmp_path / "static" / "dist"
    _write(src / "index.html", "<html></html>")
    _write(src / "app.js", "const x = 1;\n")

    adapter = FastAPIMinifyAdapter(source_dir=src, build_dir=dist)
    adapter.build()

    assert (dist / "index.html").is_file()
    assert (dist / "app.min.js").is_file()


def test_fastapi_build_copies_css(tmp_path):
    src = tmp_path / "static"
    dist = tmp_path / "dist"
    _write(src / "style.css", "body {}")

    adapter = FastAPIMinifyAdapter(source_dir=src, build_dir=dist)
    adapter.build()

    assert (dist / "style.css").is_file()
    assert (dist / "style.css").read_text() == "body {}"


def test_fastapi_build_preserves_layout(tmp_path):
    src = tmp_path / "static"
    dist = tmp_path / "dist"
    _write(src / "assets" / "js" / "app.js", "const x = 1;\n")
    _write(src / "assets" / "css" / "app.css", "body {}")

    adapter = FastAPIMinifyAdapter(source_dir=src, build_dir=dist)
    adapter.build()

    assert (dist / "assets" / "js" / "app.min.js").is_file()
    assert (dist / "assets" / "css" / "app.css").is_file()


def test_fastapi_build_empty_source(tmp_path):
    src = tmp_path / "static"
    src.mkdir()
    dist = tmp_path / "dist"

    adapter = FastAPIMinifyAdapter(source_dir=src, build_dir=dist)
    results = adapter.build()
    assert results == []


def test_fastapi_returns_results_only_for_js(tmp_path):
    src = tmp_path / "static"
    dist = tmp_path / "dist"
    _write(src / "app.js", "const x = 1;\n")
    _write(src / "index.html", "<html></html>")
    _write(src / "style.css", "body {}")

    adapter = FastAPIMinifyAdapter(source_dir=src, build_dir=dist)
    results = adapter.build()
    assert len(results) == 1
    assert results[0].code


def test_fastapi_file_outside_source(tmp_path):
    src = tmp_path / "static"
    src.mkdir()
    dist = tmp_path / "dist"
    outside = tmp_path / "outside"
    _write(outside / "x.js", "const x = 1;\n")

    adapter = FastAPIMinifyAdapter(source_dir=src, build_dir=dist)
    # Direct call, not through build().
    out = adapter.output_name(outside / "x.js")
    assert out.parent == dist