"""Tests for the public exception hierarchy."""

from __future__ import annotations

import pickle

import pytest

from minifyjs import (
    BinaryNotFoundError,
    BundleError,
    MinifyError,
    MinifyJSError,
)


def test_hierarchy():
    assert issubclass(MinifyError, MinifyJSError)
    assert issubclass(BundleError, MinifyError)
    assert issubclass(BinaryNotFoundError, MinifyJSError)


def test_raise_and_catch_base():
    with pytest.raises(MinifyJSError):
        raise MinifyError("boom")


def test_raise_and_catch_specific():
    with pytest.raises(MinifyError):
        raise MinifyError("boom")


def test_bundle_error_is_minify_error():
    with pytest.raises(MinifyError):
        raise BundleError("bundle failed")


def test_binary_not_found_message_preserved():
    try:
        raise BinaryNotFoundError("specific message")
    except BinaryNotFoundError as e:
        assert str(e) == "specific message"


# NEW tests


def test_minify_error_carries_message():
    e = MinifyError("specific message")
    assert str(e) == "specific message"


def test_bundle_error_carries_message():
    e = BundleError("specific message")
    assert str(e) == "specific message"


def test_binary_not_found_error_carries_message():
    e = BinaryNotFoundError("missing")
    assert str(e) == "missing"


def test_catching_bundle_error_as_minify_error():
    with pytest.raises(MinifyError):
        raise BundleError("x")


def test_catching_binary_error_as_base():
    with pytest.raises(MinifyJSError):
        raise BinaryNotFoundError("x")


def test_errors_can_be_pickled():
    for cls in (MinifyJSError, MinifyError, BundleError, BinaryNotFoundError):
        e = cls("msg")
        restored = pickle.loads(pickle.dumps(e))
        assert str(restored) == "msg"
