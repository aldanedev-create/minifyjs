"""Tests for Options and BundleOptions dataclasses."""

from __future__ import annotations

import dataclasses

from minifyjs import BundleOptions, Options


def test_options_defaults():
    o = Options()
    assert o.compress is False
    assert o.mangle is False
    assert o.target is None
    assert o.format is None
    assert o.sourcemap is None


def test_options_can_be_constructed_with_keywords():
    o = Options(compress=True, mangle=True, target="es2015")
    assert o.compress is True
    assert o.mangle is True
    assert o.target == "es2015"


def test_options_is_mutable():
    o = Options()
    o.compress = True
    assert o.compress is True


def test_bundle_options_inherits_from_options():
    b = BundleOptions(compress=True)
    assert b.compress is True
    assert b.platform == "browser"
    assert b.splitting is False


def test_bundle_options_defaults():
    b = BundleOptions()
    assert b.entry_points is None
    assert b.outdir is None
    assert b.outfile is None
    assert b.platform == "browser"
    assert b.splitting is False


# NEW tests


def test_options_is_dataclass():
    assert dataclasses.is_dataclass(Options)


def test_options_repr_does_not_crash():
    assert repr(Options())


def test_options_equality():
    assert Options() == Options()
    assert Options(compress=True) != Options()


def test_options_all_fields_none_by_default():
    o = Options()
    assert o.target is None
    assert o.format is None
    assert o.sourcemap is None
    assert o.banner is None
    assert o.footer is None
    assert o.legal_comments is None
    assert o.source_name is None


def test_options_bool_fields_false_by_default():
    o = Options()
    assert o.compress is False
    assert o.mangle is False


def test_options_accepts_all_fields():
    o = Options(
        compress=True,
        mangle=True,
        target="es2015",
        format="esm",
        sourcemap="inline",
        banner="/*b*/",
        footer="/*f*/",
        legal_comments="eof",
        source_name="app.js",
    )
    assert o.target == "es2015"
    assert o.source_name == "app.js"


def test_bundle_options_is_dataclass():
    assert dataclasses.is_dataclass(BundleOptions)


def test_bundle_options_all_fields():
    b = BundleOptions(
        entry_points=["a.js"],
        outdir="dist",
        platform="node",
        splitting=True,
    )
    assert b.entry_points == ["a.js"]
    assert b.outdir == "dist"
    assert b.platform == "node"
    assert b.splitting is True


def test_bundle_options_repr():
    assert repr(BundleOptions())


def test_bundle_options_equality():
    assert BundleOptions() == BundleOptions()
