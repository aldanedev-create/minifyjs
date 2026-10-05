"""Tests for the CLI argument building in _protocol.py."""

from __future__ import annotations

from minifyjs._protocol import build_args, build_bundle_args
from minifyjs.options import BundleOptions, Options


def test_build_args_default_is_empty():
    assert build_args(Options()) == ["--no-minify", "--no-mangle"]


def test_build_args_compress():
    assert "--compress" in build_args(Options(compress=True))


def test_build_args_mangle():
    assert "--mangle" in build_args(Options(mangle=True))


def test_build_args_compress_mangle():
    args = build_args(Options(compress=True, mangle=True))
    assert "--compress" in args
    assert "--mangle" in args


def test_build_args_target():
    args = build_args(Options(target="es2015"))
    assert "--target" in args
    assert "es2015" in args


def test_build_args_format():
    args = build_args(Options(format="esm"))
    assert "--format" in args
    assert "esm" in args


def test_build_args_sourcemap_uses_equals():
    args = build_args(Options(sourcemap="inline"))
    assert "--sourcemap=inline" in args


def test_build_args_banner():
    args = build_args(Options(banner="/*b*/"))
    assert "--banner" in args
    assert "/*b*/" in args


def test_build_args_footer():
    args = build_args(Options(footer="/*f*/"))
    assert "--footer" in args
    assert "/*f*/" in args


def test_build_args_legal_comments():
    args = build_args(Options(legal_comments="none"))
    assert "--legal-comments" in args
    assert "none" in args


def test_build_bundle_args_includes_bundle_flag():
    args = build_bundle_args(BundleOptions(outfile="out.js"))
    assert "--bundle" in args


def test_build_bundle_args_outdir():
    args = build_bundle_args(BundleOptions(outdir="dist"))
    assert "--outdir" in args
    assert "dist" in args


def test_build_bundle_args_outfile():
    args = build_bundle_args(BundleOptions(outfile="bundle.js"))
    assert "--outfile" in args
    assert "bundle.js" in args


def test_build_bundle_args_platform():
    args = build_bundle_args(BundleOptions(platform="node"))
    assert "--platform" in args
    assert "node" in args
