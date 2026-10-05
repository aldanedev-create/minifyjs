"""Tests for the Django adapter example."""

from __future__ import annotations

from pathlib import Path

import pytest

from minifyjs import MinifyJSError

from examples.django_adapter import DjangoMinifyAdapter


pytestmark = pytest.mark.binary


def _write(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content)


def test_django_output_inside_staticfiles_dir(tmp_path):
    static_src = tmp_path / "static" / "app"
    static_root = tmp_path / "static_root"
    _write(static_src / "main.js", "const x = 1;\n")

    adapter = DjangoMinifyAdapter(
        static_root=static_root,
        staticfiles_dirs=[tmp_path / "static"],
    )
    results = adapter.minify_directory(tmp_path / "static")

    assert len(results) == 1
    out = static_root / "app" / "main.min.js"
    assert out.is_file()


def test_django_output_preserves_relative_path(tmp_path):
    static_dir = tmp_path / "static"
    _write(static_dir / "js" / "a.js", "const a = 1;\n")
    _write(static_dir / "js" / "sub" / "b.js", "const b = 2;\n")

    adapter = DjangoMinifyAdapter(
        static_root=tmp_path / "static_root",
        staticfiles_dirs=[static_dir],
    )
    adapter.minify_directory(static_dir)

    assert (tmp_path / "static_root" / "js" / "a.min.js").is_file()
    assert (tmp_path / "static_root" / "js" / "sub" / "b.min.js").is_file()


def test_django_file_outside_any_dir_drops_at_root(tmp_path):
    static_root = tmp_path / "static_root"
    outside = tmp_path / "outside"
    _write(outside / "x.js", "const x = 1;\n")

    adapter = DjangoMinifyAdapter(
        static_root=static_root,
        staticfiles_dirs=[tmp_path / "not_this_dir"],
    )
    adapter.minify_directory(outside)

    assert (static_root / "x.min.js").is_file()


def test_django_creates_output_dirs(tmp_path):
    static_dir = tmp_path / "static"
    _write(static_dir / "deeply" / "nested" / "x.js", "const x = 1;\n")

    adapter = DjangoMinifyAdapter(
        static_root=tmp_path / "static_root",
        staticfiles_dirs=[static_dir],
    )
    adapter.minify_directory(static_dir)

    assert (tmp_path / "static_root" / "deeply" / "nested" / "x.min.js").is_file()


def test_django_empty_dir(tmp_path):
    empty = tmp_path / "empty"
    empty.mkdir()

    adapter = DjangoMinifyAdapter(
        static_root=tmp_path / "static_root",
        staticfiles_dirs=[empty],
    )
    results = adapter.minify_directory(empty)
    assert results == []


def test_django_non_recursive(tmp_path):
    static_dir = tmp_path / "static"
    _write(static_dir / "a.js", "const a = 1;\n")
    _write(static_dir / "sub" / "b.js", "const b = 2;\n")

    adapter = DjangoMinifyAdapter(
        static_root=tmp_path / "static_root",
        staticfiles_dirs=[static_dir],
    )
    results = adapter.minify_directory(static_dir, recursive=False)
    assert len(results) == 1
    assert (tmp_path / "static_root" / "a.min.js").is_file()
    assert not (tmp_path / "static_root" / "sub" / "b.min.js").exists()


def test_django_minifies_content(tmp_path):
    static_dir = tmp_path / "static"
    _write(static_dir / "app.js", "function add(a, b) {\n    return a + b;\n}\n")

    adapter = DjangoMinifyAdapter(
        static_root=tmp_path / "static_root",
        staticfiles_dirs=[static_dir],
    )
    adapter.minify_directory(static_dir)

    body = (tmp_path / "static_root" / "app.min.js").read_text()
    assert "\n" not in body


def test_django_invalid_js_raises(tmp_path):
    static_dir = tmp_path / "static"
    _write(static_dir / "bad.js", "function () { } }")

    adapter = DjangoMinifyAdapter(
        static_root=tmp_path / "static_root",
        staticfiles_dirs=[static_dir],
    )
    with pytest.raises(MinifyJSError):
        adapter.minify_directory(static_dir)