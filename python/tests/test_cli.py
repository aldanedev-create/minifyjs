"""Tests for the python -m minifyjs entry point."""

from __future__ import annotations

import subprocess
import sys

import pytest

from .conftest import BINARY_PATH

pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_module_version():
    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs", "--version"],
        capture_output=True,
        check=False,
    )
    assert proc.returncode == 0
    assert b"minifyjs" in proc.stdout
    assert b"esbuild" in proc.stdout


def test_module_help():
    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs", "--help"],
        capture_output=True,
        check=False,
    )
    assert proc.returncode == 0
    assert b"Usage:" in proc.stdout


def test_module_stdin_to_stdout():
    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs"],
        input=b"function add(a, b) { return a + b; }\n",
        capture_output=True,
        check=False,
    )
    assert proc.returncode == 0
    assert b"function" in proc.stdout


def test_module_unknown_flag():
    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs", "--nope"],
        capture_output=True,
        check=False,
    )
    assert proc.returncode == 2


# NEW tests


def test_module_no_args_empty_stdin():
    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs"],
        input=b"",
        capture_output=True,
        check=False,
    )
    assert proc.returncode == 0
    assert proc.stdout == b""


def test_module_with_compress_flag():
    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs", "--compress"],
        input=b"const x = 1 + 2 + 3;",
        capture_output=True,
        check=False,
    )
    assert proc.returncode == 0
    assert b"6" in proc.stdout


def test_module_with_target_es5():
    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs", "--target", "es5"],
        input=b"const f = (a) => a;",
        capture_output=True,
        check=False,
    )
    assert proc.returncode == 0
    assert b"=>" not in proc.stdout


def test_module_output_to_stdout_exact():
    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs"],
        input=b"const x = 1;",
        capture_output=True,
        check=False,
    )
    assert proc.stdout == b"const x=1;"


def test_module_error_to_stderr():
    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs"],
        input=b"function () { } }",
        capture_output=True,
        check=False,
    )
    assert proc.returncode != 0
    assert proc.stdout == b""
    assert len(proc.stderr) > 0


def test_module_exit_code_on_syntax_error():
    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs"],
        input=b"function () { } }",
        capture_output=True,
        check=False,
    )
    assert proc.returncode == 1


def test_module_exit_code_on_unknown_flag():
    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs", "--not-a-flag"],
        capture_output=True,
        check=False,
    )
    assert proc.returncode == 2


def test_module_quiet_suppresses_stderr():
    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs", "--quiet"],
        input=b"const x = 1;",
        capture_output=True,
        check=False,
    )
    assert proc.returncode == 0
    assert proc.stderr == b""
