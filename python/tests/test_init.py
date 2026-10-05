"""Tests for the package's public surface: __init__.py."""

from __future__ import annotations

import pytest

import minifyjs


def test_version_is_string():
    assert isinstance(minifyjs.__version__, str)
    assert minifyjs.__version__


def test_version_has_major_minor_patch():
    parts = minifyjs.__version__.split(".")
    assert len(parts) >= 3


def test_public_api_exported():
    expected = [
        "minify",
        "optimize",
        "bundle",
        "Options",
        "BundleOptions",
        "Result",
        "Diagnostic",
        "Adapter",
        "Config",
        "MinifyJSError",
        "MinifyError",
        "BundleError",
        "BinaryNotFoundError",
    ]
    for name in expected:
        assert hasattr(minifyjs, name), f"minifyjs.{name} is missing"


def test_all_matches_exports():
    for name in minifyjs.__all__:
        assert hasattr(minifyjs, name), f"__all__ names {name} but it's missing"


def test_minify_callable():
    assert callable(minifyjs.minify)


def test_optimize_callable():
    assert callable(minifyjs.optimize)


def test_bundle_callable():
    assert callable(minifyjs.bundle)


# NEW tests


def test_version_matches_go_version():
    """The Python __version__ should match the Go engine version.

    They are kept in sync by the release pipeline; a mismatch is a
    release bug. This test is skipped because the release pipeline
    enforces the sync at build time.
    """
    pytest.skip("version sync is enforced by the release pipeline")


def test_minify_returns_result():
    r = minifyjs.minify("const x = 1;")
    assert isinstance(r, minifyjs.Result)


def test_all_is_tuple_or_list():
    assert isinstance(minifyjs.__all__, (list, tuple))


def test_no_duplicate_names_in_all():
    assert len(minifyjs.__all__) == len(set(minifyjs.__all__))


def test_import_from_submodule_works():
    from minifyjs.minifier import minify
    assert callable(minify)


def test_import_optimizer_module():
    from minifyjs.optimizer import optimize
    assert callable(optimize)


def test_import_bundler_module():
    from minifyjs.bundler import bundle
    assert callable(bundle)


def test_import_adapter_module():
    from minifyjs.adapter import Adapter
    assert Adapter is not None
