"""End-to-end tests for optimize()."""

from __future__ import annotations

import pytest

from minifyjs import Result, optimize

from .conftest import BINARY_PATH

pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_optimize_folds():
    r = optimize("const x = 1 + 2 + 3;")
    assert "6" in r.code


def test_optimize_mangles():
    r = optimize("function f(longName) { return longName; }")
    assert r.code.count("longName") < 2


def test_optimize_with_target():
    r = optimize("const f = (a) => a + 1;", target="es5")
    assert "=>" not in r.code


def test_optimize_with_banner():
    r = optimize("const x = 1;", banner="/*b*/")
    assert r.code.startswith("/*b*/")


def test_optimize_with_sourcemap():
    r = optimize("const x = 1;", sourcemap="external")
    assert r.map


def test_optimize_empty():
    r = optimize("")
    assert r.code == ""


def test_optimize_syntax_error():
    from minifyjs import MinifyError
    with pytest.raises(MinifyError):
        optimize("function () { } }")


def test_optimize_bytes_saved():
    src = "function add(a, b) {\n    return a + b;\n}\n"
    r = optimize(src)
    assert r.bytes_saved > 0


# NEW tests


def test_optimize_returns_result():
    assert isinstance(optimize("const x = 1;"), Result)


def test_optimize_compress_enabled():
    r = optimize("const x = 1 + 2 + 3;")
    assert "6" in r.code


def test_optimize_mangle_enabled():
    r = optimize("function f(longName) { return longName; }")
    assert r.code.count("longName") < 2


def test_optimize_both():
    r = optimize("function f(x) { const y = 1 + 2; return x + y; }")
    assert not r.has_errors


def test_optimize_with_format():
    r = optimize("export const x = 1;", format="esm")
    assert "export" in r.code


def test_optimize_with_legal_comments():
    r = optimize("/*! (c) */ const x = 1;", legal_comments="inline")
    assert "(c)" in r.code


def test_optimize_with_source_name():
    r = optimize("const x = 1;", source_name="app.js")
    assert not r.has_errors


def test_optimize_class():
    r = optimize("class A { constructor(x) { this.x = x; } }")
    assert not r.has_errors


def test_optimize_nested_function():
    r = optimize("function outer(x) { return function inner(y) { return x + y; }; }")
    assert not r.has_errors


def test_optimize_with_footer():
    r = optimize("const x = 1;", footer="/*f*/")
    assert r.code.endswith("/*f*/")
