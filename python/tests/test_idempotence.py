"""Tests that minifying already-minified code is a no-op."""

from __future__ import annotations

import pytest

from minifyjs import minify, optimize

from .conftest import BINARY_PATH

pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def _assert_idempotent(src, **opts):
    r1 = minify(src, **opts)
    r2 = minify(r1.code, **opts)
    assert r1.code == r2.code, (
        f"not idempotent for input {src!r}\n"
        f"first:  {r1.code!r}\n"
        f"second: {r2.code!r}"
    )


def test_idempotent_simple():
    _assert_idempotent("function add(a, b) { return a + b; }")


def test_idempotent_with_newlines():
    _assert_idempotent("function f() {\n    return 1;\n}\n")


def test_idempotent_with_comments():
    _assert_idempotent("// comment\nconst x = 1;\n")


def test_idempotent_compress():
    _assert_idempotent("const x = 1 + 2 + 3;", compress=True)


def test_idempotent_mangle():
    _assert_idempotent("function f(longName) { return longName; }", mangle=True)


def test_idempotent_compress_mangle():
    _assert_idempotent(
        "function f(x) { const y = 1 + 2; return x + y; }",
        compress=True, mangle=True,
    )


def test_idempotent_modern_syntax():
    _assert_idempotent("const f = (x) => x?.y ?? 0;")


def test_idempotent_class():
    _assert_idempotent("class A { constructor(x) { this.x = x; } }")


def test_idempotent_target_es5():
    _assert_idempotent("const f = (a) => a + 1;", target="es5")


def test_optimize_idempotent():
    src = "const x = 1 + 2 + 3;"
    r1 = optimize(src)
    r2 = optimize(r1.code)
    assert r1.code == r2.code
