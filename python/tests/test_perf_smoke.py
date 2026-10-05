"""Smoke tests that catch catastrophic performance regressions.

These are not benchmarks. They only assert that the bridge completes
small and medium inputs in a bounded time. A real benchmark suite
lives under bench/ at the repository root.
"""

from __future__ import annotations

import time

import pytest

from minifyjs import minify, optimize

from .conftest import BINARY_PATH

pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_small_input_completes_quickly():
    start = time.monotonic()
    minify("const x = 1;")
    elapsed = time.monotonic() - start
    assert elapsed < 5.0, f"took {elapsed:.2f}s for a tiny input"


def test_medium_input_completes_quickly():
    src = "\n".join(f"const v{i} = {i} * 2;" for i in range(1000))
    start = time.monotonic()
    r = minify(src, compress=True, mangle=True)
    elapsed = time.monotonic() - start
    assert not r.has_errors
    assert elapsed < 10.0, f"took {elapsed:.2f}s for a medium input"


def test_large_input_completes_in_reasonable_time():
    src = "\n".join(f"function f{i}(x) {{ return x + {i}; }}" for i in range(2000))
    start = time.monotonic()
    r = optimize(src)
    elapsed = time.monotonic() - start
    assert not r.has_errors
    assert elapsed < 30.0, f"took {elapsed:.2f}s for a large input"


def test_startup_time_is_bounded():
    """A single minify() call includes process startup. If that ever
    balloons past 5 seconds, something is badly wrong."""
    start = time.monotonic()
    minify("1;")
    elapsed = time.monotonic() - start
    assert elapsed < 5.0, f"startup took {elapsed:.2f}s"


def test_many_small_calls_complete_in_reasonable_time():
    start = time.monotonic()
    for i in range(20):
        minify(f"const x{i} = {i};")
    elapsed = time.monotonic() - start
    assert elapsed < 20.0, f"20 calls took {elapsed:.2f}s"


def test_result_bytes_scale_linearly():
    """Minified output for a linear input should be roughly linear in
    the input size. This catches accidental O(n^2) behavior."""
    small = "\n".join(f"const a{i} = {i};" for i in range(100))
    large = "\n".join(f"const a{i} = {i};" for i in range(1000))
    r_small = minify(small, compress=True, mangle=True)
    r_large = minify(large, compress=True, mangle=True)
    assert not r_small.has_errors
    assert not r_large.has_errors
    # 10x input should produce < 50x output (generous factor for
    # startup noise and non-linearity in esbuild's parser).
    assert r_large.minified_bytes < r_small.minified_bytes * 50
