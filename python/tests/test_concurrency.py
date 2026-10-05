"""Tests that the Python bridge is safe under concurrency."""

from __future__ import annotations

import concurrent.futures

import pytest

from minifyjs import minify, optimize

from .conftest import BINARY_PATH


pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_concurrent_same_input():
    src = "function add(a, b) { return a + b; }"
    with concurrent.futures.ThreadPoolExecutor(max_workers=8) as ex:
        futures = [ex.submit(minify, src) for _ in range(40)]
        results = [f.result() for f in futures]
    first = results[0].code
    for r in results:
        assert r.code == first


def test_concurrent_different_inputs():
    inputs = [f"const v{i} = {i};" for i in range(40)]
    with concurrent.futures.ThreadPoolExecutor(max_workers=8) as ex:
        futures = [ex.submit(minify, s) for s in inputs]
        results = [f.result() for f in futures]
    assert len(results) == 40
    for r in results:
        assert not r.has_errors


def test_concurrent_optimize():
    src = "const x = 1 + 2 + 3;"
    with concurrent.futures.ThreadPoolExecutor(max_workers=8) as ex:
        futures = [ex.submit(optimize, src) for _ in range(20)]
        results = [f.result() for f in futures]
    for r in results:
        assert "6" in r.code


def test_concurrent_syntax_errors():
    src = "function () { } }"
    from minifyjs import MinifyError
    errors = 0
    with concurrent.futures.ThreadPoolExecutor(max_workers=8) as ex:
        futures = [ex.submit(minify, src) for _ in range(20)]
        for f in futures:
            try:
                f.result()
            except MinifyError:
                errors += 1
    assert errors == 20


def test_concurrent_with_varied_options():
    def run(i):
        return minify(
            f"const v{i} = {i} + 1;",
            compress=(i % 2 == 0),
            target="es5" if i % 3 == 0 else None,
        )
    with concurrent.futures.ThreadPoolExecutor(max_workers=8) as ex:
        futures = [ex.submit(run, i) for i in range(30)]
        results = [f.result() for f in futures]
    for r in results:
        assert not r.has_errors


def test_concurrent_results_are_independent():
    src = "const x = 1;"
    with concurrent.futures.ThreadPoolExecutor(max_workers=8) as ex:
        futures = [ex.submit(minify, src) for _ in range(20)]
        results = [f.result() for f in futures]
    # Mutating one Result must not affect another.
    results[0].diagnostics.append("fake")
    for r in results[1:]:
        assert "fake" not in r.diagnostics