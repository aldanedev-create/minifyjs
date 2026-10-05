"""Tests for the CI-pipeline example."""

from __future__ import annotations

from pathlib import Path

import pytest

from examples.ci_pipeline import CIPipelineAdapter, _main
from minifyjs import Options


pytestmark = pytest.mark.binary


def _write(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content)


def test_ci_minifies_in_place(tmp_path):
    _write(tmp_path / "a.js", "const a = 1;\n")
    _write(tmp_path / "b.js", "const b = 2;\n")

    adapter = CIPipelineAdapter(root=tmp_path, options=Options(compress=True, mangle=True))
    rc = adapter.run()

    assert rc == 0
    assert (tmp_path / "a.min.js").is_file()
    assert (tmp_path / "b.min.js").is_file()


def test_ci_returns_one_on_failure(tmp_path):
    _write(tmp_path / "good.js", "const a = 1;\n")
    _write(tmp_path / "bad.js", "function () { } }")

    adapter = CIPipelineAdapter(root=tmp_path, options=Options(compress=True, mangle=True))
    rc = adapter.run()

    assert rc == 1
    assert len(adapter.failures) == 1


def test_ci_skips_already_minified(tmp_path):
    _write(tmp_path / "app.js", "const x = 1;\n")
    _write(tmp_path / "app.min.js", "const x=1;")

    adapter = CIPipelineAdapter(root=tmp_path, options=Options(compress=True, mangle=True))
    adapter.run()

    assert not (tmp_path / "app.min.min.js").exists()


def test_ci_empty_dir(tmp_path):
    adapter = CIPipelineAdapter(root=tmp_path, options=Options(compress=True, mangle=True))
    rc = adapter.run()
    assert rc == 0
    assert adapter.failures == []


def test_ci_main_missing_arg(capsys):
    rc = _main(["ci_pipeline.py"])
    assert rc == 2


def test_ci_main_bad_dir(tmp_path, capsys):
    rc = _main(["ci_pipeline.py", str(tmp_path / "nope")])
    assert rc == 1


def test_ci_main_success(tmp_path):
    _write(tmp_path / "app.js", "const x = 1;\n")
    rc = _main(["ci_pipeline.py", str(tmp_path)])
    assert rc == 0