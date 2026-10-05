"""Tests for the watch-directory example."""

from __future__ import annotations

import time
from pathlib import Path

import pytest

from examples.watch_directory import WatchingAdapter
from minifyjs import Options


pytestmark = pytest.mark.binary


def _write(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content)


def test_watch_first_scan_processes_all(tmp_path):
    _write(tmp_path / "a.js", "const a = 1;\n")
    _write(tmp_path / "b.js", "const b = 2;\n")

    adapter = WatchingAdapter(root=tmp_path, options=Options(compress=True, mangle=True))
    changed = adapter.scan()
    assert len(changed) == 2


def test_watch_second_scan_processes_nothing(tmp_path):
    _write(tmp_path / "a.js", "const a = 1;\n")

    adapter = WatchingAdapter(root=tmp_path, options=Options(compress=True, mangle=True))
    adapter.scan()
    # Second scan without touching files.
    changed = adapter.scan()
    assert changed == []


def test_watch_detects_change(tmp_path):
    p = tmp_path / "a.js"
    _write(p, "const a = 1;\n")

    adapter = WatchingAdapter(root=tmp_path, options=Options(compress=True, mangle=True))
    adapter.scan()

    # Modify the file with a later mtime.
    time.sleep(0.01)
    p.write_text("const a = 2;\n")
    # Force mtime to be newer than the recorded one.
    future = time.time() + 1
    import os
    os.utime(p, (future, future))

    changed = adapter.scan()
    assert p in changed


def test_watch_run_once_returns_zero_on_success(tmp_path):
    _write(tmp_path / "a.js", "const a = 1;\n")

    adapter = WatchingAdapter(root=tmp_path, options=Options(compress=True, mangle=True))
    rc = adapter.run_once()
    assert rc == 0
    assert (tmp_path / "a.min.js").is_file()


def test_watch_run_once_returns_failures(tmp_path):
    _write(tmp_path / "bad.js", "function () { } }")

    adapter = WatchingAdapter(root=tmp_path, options=Options(compress=True, mangle=True))
    rc = adapter.run_once()
    assert rc == 1


def test_watch_output_name(tmp_path):
    adapter = WatchingAdapter(root=tmp_path, options=Options())
    assert adapter.output_name(tmp_path / "a.js") == tmp_path / "a.min.js"