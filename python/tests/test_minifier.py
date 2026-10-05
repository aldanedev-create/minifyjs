"""End-to-end tests for minify()."""

from __future__ import annotations

import pytest

from minifyjs import Result, minify

from .conftest import BINARY_PATH


pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_basic_minify():
    src = "function add(a, b) {\n    return a + b;\n}\n"
    r = minify(src)
    assert r.code == "function add(a,b){return a+b;}"
    assert r.minified_bytes < r.original_bytes


def test_comments_removed():
    r = minify("// hello\nconst x = 1;\n")
    assert "//" not in r.code
    assert r.code == "const x=1;"


def test_compress_folds_constants():
    r = minify("const x = 1 + 2 + 3;", compress=True)
    assert "6" in r.code


def test_mangle_renames_locals():
    r = minify("function f(longName) { return longName; }", mangle=True)
    assert r.code.count("longName") < 2


def test_target_es5_lowers_arrows():
    r = minify("const f = (a) => a + 1;", target="es5")
    assert "=>" not in r.code


def test_target_esnext_keeps_arrows():
    r = minify("const f = (a) => a + 1;", target="esnext")
    assert "=>" in r.code


def test_banner_prepended():
    r = minify("const x = 1;", banner="/*b*/")
    assert r.code.startswith("/*b*/")


def test_footer_appended():
    r = minify("const x = 1;", footer="/*f*/")
    assert r.code.endswith("/*f*/")


def test_sourcemap_inline():
    r = minify("const x = 1;", sourcemap="inline")
    assert "sourceMappingURL=data:application/json" in r.code


def test_sourcemap_external_returns_map():
    r = minify("const x = 1;", sourcemap="external")
    assert r.map
    assert '"version":3' in r.map


def test_source_name_in_diagnostic():
    from minifyjs import MinifyError
    try:
        minify("const = 1;", source_name="myfile.js")
    except MinifyError as e:
        assert "myfile.js" in str(e)
    else:
        pytest.fail("expected MinifyError")


def test_ratio_between_zero_and_one():
    r = minify("function add(a, b) { return a + b; }")
    assert 0 < r.ratio <= 1


def test_empty_input():
    r = minify("")
    assert r.code == ""


def test_whitespace_only_input():
    r = minify("   \n\t  ")
    assert r.code == ""


def test_result_diagnostics_list():
    r = minify("const x = 1;")
    assert isinstance(r.diagnostics, list)


def test_syntax_error_raises():
    from minifyjs import MinifyError
    with pytest.raises(MinifyError):
        minify("function () { } }")


def test_legal_comments_none():
    r = minify("/*! (c) */ const x = 1;", legal_comments="none")
    assert "(c)" not in r.code


def test_legal_comments_eof():
    r = minify("/*! (c) */ const x = 1;", legal_comments="eof")
    assert "(c)" in r.code


# NEW tests


def test_minify_returns_result_instance():
    assert isinstance(minify("const x = 1;"), Result)


def test_minify_source_name_kwarg():
    r = minify("const x = 1;", source_name="app.js")
    assert not r.has_errors


def test_minify_compress_only():
    r = minify("const x = 1 + 2 + 3;", compress=True)
    assert "6" in r.code


def test_minify_mangle_only():
    r = minify("function f(longName) { return longName; }", mangle=True)
    assert r.code.count("longName") < 2


def test_minify_compress_and_mangle():
    r = minify(
        "function f(x) { const y = 1 + 2; return x + y; }",
        compress=True, mangle=True,
    )
    assert not r.has_errors


def test_minify_with_format_esm():
    r = minify("export const x = 1;", format="esm")
    assert "export" in r.code


def test_minify_with_format_iife():
    r = minify("const x = 1;", format="iife")
    assert not r.has_errors


def test_minify_with_legal_comments_inline():
    r = minify("/*! (c) */ const x = 1;", legal_comments="inline")
    assert "(c)" in r.code


def test_minify_with_legal_comments_eof():
    r = minify("const x = 1; /*! tail */", legal_comments="eof")
    assert "tail" in r.code


def test_minify_legal_comments_default():
    r = minify("/*! (c) */ const x = 1;")
    assert "(c)" in r.code


def test_minify_banner_and_footer_together():
    r = minify("const x = 1;", banner="/*b*/", footer="/*f*/")
    assert r.code.startswith("/*b*/")
    assert r.code.endswith("/*f*/")


def test_minify_ratio():
    r = minify("function add(a, b) { return a + b; }")
    assert 0 < r.ratio < 1


def test_minify_bytes_saved_positive():
    r = minify("function add(a, b) { return a + b; }")
    assert r.bytes_saved > 0


def test_minify_original_bytes_int():
    r = minify("const x = 1;")
    assert isinstance(r.original_bytes, int)


def test_minify_minified_bytes_int():
    r = minify("const x = 1;")
    assert isinstance(r.minified_bytes, int)


def test_minify_diagnostics_is_list():
    r = minify("const x = 1;")
    assert isinstance(r.diagnostics, list)


def test_minify_shebang_preserved():
    r = minify("#!/usr/bin/env node\nconst x = 1;")
    assert "#!/usr/bin/env node" in r.code


def test_minify_async():
    r = minify("async function f() { return await g(); }")
    assert not r.has_errors


def test_minify_generator():
    r = minify("function* gen() { yield 1; }")
    assert not r.has_errors


def test_minify_class():
    r = minify("class A { constructor(x) { this.x = x; } }")
    assert not r.has_errors