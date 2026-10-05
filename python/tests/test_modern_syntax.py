"""Tests that modern JavaScript features are handled."""

from __future__ import annotations

import pytest

from minifyjs import minify

from .conftest import BINARY_PATH


pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_arrow_function():
    r = minify("const f = (x) => x * 2;")
    assert not r.has_errors


def test_template_literal():
    r = minify("const s = `hello ${name}`;")
    assert not r.has_errors


def test_destructuring_object():
    r = minify("const { a, b } = obj;")
    assert not r.has_errors


def test_destructuring_array():
    r = minify("const [x, y] = arr;")
    assert not r.has_errors


def test_spread():
    r = minify("const c = [...a, ...b];")
    assert not r.has_errors


def test_rest_parameters():
    r = minify("function f(...args) { return args.length; }")
    assert not r.has_errors


def test_default_parameters():
    r = minify("function f(x = 1) { return x; }")
    assert not r.has_errors


def test_class():
    r = minify("class A { constructor(x) { this.x = x; } }")
    assert not r.has_errors


def test_private_fields():
    r = minify("class A { #x = 1; get() { return this.#x; } }")
    assert not r.has_errors


def test_static_fields():
    r = minify("class A { static x = 1; }")
    assert not r.has_errors


def test_async_await():
    r = minify("async function f() { return await g(); }")
    assert not r.has_errors


def test_generators():
    r = minify("function* gen() { yield 1; }")
    assert not r.has_errors


def test_optional_chaining():
    r = minify("const x = a?.b?.c;")
    assert not r.has_errors


def test_nullish_coalescing():
    r = minify("const x = a ?? b;")
    assert not r.has_errors


def test_logical_assignment():
    r = minify("a ||= b; a &&= c; a ??= d;")
    assert not r.has_errors


def test_for_of():
    r = minify("for (const x of arr) { console.log(x); }")
    assert not r.has_errors


def test_top_level_await():
    r = minify("const x = await fetch('/');")
    assert not r.has_errors


def test_bigint():
    r = minify("const n = 123n;")
    assert not r.has_errors


# NEW tests


def test_arrow_with_block_body():
    r = minify("const f = (x) => { return x * 2; };")
    assert not r.has_errors


def test_arrow_implicit_return():
    r = minify("const f = (x) => x * 2;")
    assert not r.has_errors


def test_template_nested_interpolation():
    r = minify("const s = `a ${`b ${x}`} c`;")
    assert not r.has_errors


def test_object_computed_key():
    r = minify("const o = { [key]: 1 };")
    assert not r.has_errors


def test_object_shorthand():
    r = minify("const o = { x, y };")
    assert not r.has_errors


def test_object_method():
    r = minify("const o = { greet() { return 'hi'; } };")
    assert not r.has_errors


def test_class_getter_setter():
    r = minify("class A { get x() { return 1; } set x(v) {} }")
    assert not r.has_errors


def test_class_static_method():
    r = minify("class A { static create() { return new A(); } }")
    assert not r.has_errors


def test_class_extends():
    r = minify("class B extends A { constructor() { super(); } }")
    assert not r.has_errors


def test_class_private_method():
    r = minify("class A { #helper() { return 1; } call() { return this.#helper(); } }")
    assert not r.has_errors


def test_async_arrow():
    r = minify("const f = async (x) => await x;")
    assert not r.has_errors


def test_async_iife():
    r = minify("(async () => { await g(); })();")
    assert not r.has_errors


def test_await_in_try():
    r = minify("async function f() { try { await g(); } catch (e) {} }")
    assert not r.has_errors


def test_generator_yield_star():
    r = minify("function* g() { yield* other(); }")
    assert not r.has_errors


def test_for_await_of():
    r = minify("async function f() { for await (const x of xs) { use(x); } }")
    assert not r.has_errors


def test_optional_chaining_call():
    r = minify("const x = obj.method?.();")
    assert not r.has_errors


def test_optional_chaining_index():
    r = minify("const x = arr?.[0];")
    assert not r.has_errors


def test_nullish_with_chaining():
    r = minify("const x = a?.b?.c ?? d;")
    assert not r.has_errors


def test_numeric_separators():
    r = minify("const x = 1_000_000;")
    assert not r.has_errors


def test_bigint_hex():
    r = minify("const x = 0xFFn;")
    assert not r.has_errors


def test_import_meta():
    r = minify("const url = import.meta.url;", format="esm")
    assert not r.has_errors


def test_new_target():
    r = minify("function F() { if (new.target) return; }")
    assert not r.has_errors