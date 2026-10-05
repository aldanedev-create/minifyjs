"""Tests for the static-site build example."""

from __future__ import annotations

from pathlib import Path

import pytest

from examples.static_site_build import StaticSiteAdapter


pytestmark = pytest.mark.binary


def _write(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content)


def test_static_site_build_creates_dist(tmp_path):
    src = tmp_path / "src"
    dist = tmp_path / "dist"
    _write(src / "index.html", "<html></html>")
    _write(src / "app.js", "const x = 1;\n")

    adapter = StaticSiteAdapter(source_dir=src, dist_dir=dist)
    adapter.build()

    assert dist.is_dir()
    assert (dist / "index.html").is_file()
    assert (dist / "app.js").is_file()


def test_static_site_build_minifies_js_in_place(tmp_path):
    src = tmp_path / "src"
    dist = tmp_path / "dist"
    _write(src / "app.js", "function add(a, b) {\n    return a + b;\n}\n")

    adapter = StaticSiteAdapter(source_dir=src, dist_dir=dist)
    adapter.build()

    body = (dist / "app.js").read_text()
    assert "\n" not in body


def test_static_site_build_copies_non_js_verbatim(tmp_path):
    src = tmp_path / "src"
    dist = tmp_path / "dist"
    html = "<html><body>content</body></html>"
    _write(src / "index.html", html)

    adapter = StaticSiteAdapter(source_dir=src, dist_dir=dist)
    adapter.build()

    assert (dist / "index.html").read_text() == html


def test_static_site_build_clears_existing_dist(tmp_path):
    src = tmp_path / "src"
    dist = tmp_path / "dist"
    dist.mkdir()
    (dist / "old.txt").write_text("stale")
    _write(src / "app.js", "const x = 1;\n")

    adapter = StaticSiteAdapter(source_dir=src, dist_dir=dist)
    adapter.build()

    assert not (dist / "old.txt").exists()


def test_static_site_build_preserves_nested_layout(tmp_path):
    src = tmp_path / "src"
    dist = tmp_path / "dist"
    _write(src / "js" / "app.js", "const x = 1;\n")
    _write(src / "css" / "app.css", "body {}")

    adapter = StaticSiteAdapter(source_dir=src, dist_dir=dist)
    adapter.build()

    assert (dist / "js" / "app.js").is_file()
    assert (dist / "css" / "app.css").is_file()


def test_static_site_build_returns_only_js_results(tmp_path):
    src = tmp_path / "src"
    dist = tmp_path / "dist"
    _write(src / "app.js", "const x = 1;\n")
    _write(src / "index.html", "<html></html>")

    adapter = StaticSiteAdapter(source_dir=src, dist_dir=dist)
    results = adapter.build()
    assert len(results) == 1