"""Tests for unusual but valid inputs."""

from __future__ import annotations

import pytest

from minifyjs import Adapter, minify

from .conftest import BINARY_PATH

pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_single_character_identifier():
    r = minify("const a = 1;")
    assert not r.has_errors


def test_single_character_property():
    r = minify("const o = { a: 1 };")
    assert not r.has_errors


def test_very_long_identifier():
    name = "a" * 5000
    r = minify(f"const {name} = 1;")
    assert not r.has_errors


def test_very_long_string_literal():
    s = "x" * 100000
    r = minify(f'const s = "{s}";')
    assert not r.has_errors


def test_deeply_nested_arrays():
    src = "const a = " + "[" * 100 + "1" + "]" * 100 + ";"
    r = minify(src)
    assert not r.has_errors


def test_deeply_nested_objects():
    src = "const a = " + "{ x: " * 50 + "1" + " }" * 50 + ";"
    r = minify(src)
    assert not r.has_errors


def test_many_top_level_variables():
    src = "\n".join(f"const v{i} = {i};" for i in range(500))
    r = minify(src)
    assert not r.has_errors


def test_only_semicolons():
    r = minify(";;;;;")
    assert not r.has_errors


def test_only_whitespace_and_semicolons():
    r = minify("  ;\n  ;\n  ;  ")
    assert not r.has_errors


def test_single_expression_statement():
    r = minify("1 + 2;")
    assert not r.has_errors


def test_single_string_statement():
    r = minify('"use strict";')
    assert not r.has_errors


def test_directive_prologue():
    r = minify('"use strict";\nconst x = 1;')
    assert not r.has_errors


def test_empty_class_body():
    r = minify("class A {}")
    assert not r.has_errors


def test_empty_function_body():
    r = minify("function f() {}")
    assert not r.has_errors


def test_empty_arrow_body():
    r = minify("const f = () => {};")
    assert not r.has_errors


def test_adapter_handles_filename_with_no_extension(tmp_path):
    src = tmp_path / "app"
    src.write_text("const x = 1;")
    a = Adapter()
    # should_process returns False for extension-less files, so
    # minify_directory skips it. Direct minify_file still works.
    r = a.minify_file(src)
    assert not r.has_errors
