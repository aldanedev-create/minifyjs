"""End-to-end tests for bundle()."""

from __future__ import annotations

import pytest

from minifyjs import BundleError, Result, bundle

from .conftest import BINARY_PATH

pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_bundle_requires_outdir_or_outfile(tmp_path):
    with pytest.raises(ValueError):
        bundle([str(tmp_path / "main.js")])


def test_bundle_rejects_both_outdir_and_outfile(tmp_path):
    with pytest.raises(ValueError):
        bundle(
            [str(tmp_path / "main.js")],
            outdir=str(tmp_path / "dist"),
            outfile=str(tmp_path / "bundle.js"),
        )


def test_bundle_simple_esm(tmp_path):
    (tmp_path / "lib.js").write_text("export const v = 1;\n")
    (tmp_path / "main.js").write_text(
        'import { v } from "./lib.js"; console.log(v);\n'
    )
    out = tmp_path / "bundle.js"
    r = bundle([str(tmp_path / "main.js")], outfile=str(out), format="esm")
    assert out.is_file()
    assert not r.has_errors


def test_bundle_missing_entry_raises(tmp_path):
    with pytest.raises(BundleError):
        bundle(
            [str(tmp_path / "nope.js")],
            outfile=str(tmp_path / "bundle.js"),
        )


def test_bundle_with_minify(tmp_path):
    (tmp_path / "main.js").write_text("function add(a, b) { return a + b; }\n")
    out = tmp_path / "bundle.js"
    r = bundle(
        [str(tmp_path / "main.js")],
        outfile=str(out),
        format="esm",
        compress=True,
        mangle=True,
    )
    assert not r.has_errors
    body = out.read_text()
    assert "\n    " not in body


def test_bundle_outdir_creates_multiple_files(tmp_path):
    (tmp_path / "shared.js").write_text("export const shared = 1;\n")
    (tmp_path / "a.js").write_text(
        'import { shared } from "./shared.js"; console.log("a", shared);\n'
    )
    (tmp_path / "b.js").write_text(
        'import { shared } from "./shared.js"; console.log("b", shared);\n'
    )
    outdir = tmp_path / "dist"
    r = bundle(
        [str(tmp_path / "a.js"), str(tmp_path / "b.js")],
        outdir=str(outdir),
        format="esm",
        splitting=True,
    )
    assert not r.has_errors
    assert outdir.is_dir()
    assert len(list(outdir.iterdir())) >= 2


# NEW tests


def test_bundle_returns_result(tmp_path):
    (tmp_path / "main.js").write_text("const x = 1;\n")
    out = tmp_path / "bundle.js"
    r = bundle([str(tmp_path / "main.js")], outfile=str(out), format="esm")
    assert isinstance(r, Result)


def test_bundle_default_format_is_esm(tmp_path):
    (tmp_path / "main.js").write_text("export const x = 1;\n")
    out = tmp_path / "bundle.js"
    r = bundle([str(tmp_path / "main.js")], outfile=str(out))
    assert not r.has_errors


def test_bundle_multiple_entry_points(tmp_path):
    (tmp_path / "a.js").write_text("const a = 1;\n")
    (tmp_path / "b.js").write_text("const b = 2;\n")
    outdir = tmp_path / "dist"
    r = bundle(
        [str(tmp_path / "a.js"), str(tmp_path / "b.js")],
        outdir=str(outdir),
        format="esm",
    )
    assert not r.has_errors


def test_bundle_with_custom_banner(tmp_path):
    (tmp_path / "main.js").write_text("const x = 1;\n")
    out = tmp_path / "bundle.js"
    r = bundle(
        [str(tmp_path / "main.js")],
        outfile=str(out),
        format="esm",
        banner="/* (c) */",
    )
    assert not r.has_errors


def test_bundle_target_es2015(tmp_path):
    (tmp_path / "main.js").write_text("const f = (a) => a;\n")
    out = tmp_path / "bundle.js"
    r = bundle(
        [str(tmp_path / "main.js")],
        outfile=str(out),
        format="esm",
        target="es2015",
    )
    assert not r.has_errors


def test_bundle_platform_browser(tmp_path):
    (tmp_path / "main.js").write_text("const x = 1;\n")
    out = tmp_path / "bundle.js"
    r = bundle(
        [str(tmp_path / "main.js")],
        outfile=str(out),
        format="esm",
        platform="browser",
    )
    assert not r.has_errors


def test_bundle_platform_node(tmp_path):
    (tmp_path / "main.js").write_text("const x = 1;\n")
    out = tmp_path / "bundle.js"
    r = bundle(
        [str(tmp_path / "main.js")],
        outfile=str(out),
        format="esm",
        platform="node",
    )
    assert not r.has_errors


def test_bundle_sourcemap_external(tmp_path):
    (tmp_path / "main.js").write_text("const x = 1;\n")
    out = tmp_path / "bundle.js"
    r = bundle(
        [str(tmp_path / "main.js")],
        outfile=str(out),
        format="esm",
        sourcemap="external",
    )
    assert not r.has_errors
