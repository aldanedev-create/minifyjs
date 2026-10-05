"""Tests for the platform detection helpers."""

from __future__ import annotations

from minifyjs._platform import binary_filename, platform_tag


def test_binary_filename_nonempty():
    name = binary_filename()
    assert name
    assert name in ("minifyjs", "minifyjs.exe")


def test_platform_tag_shape():
    tag = platform_tag()
    assert "-" in tag
    system, machine = tag.split("-", 1)
    assert system in ("linux", "macos", "windows")
    assert machine in ("x86_64", "arm64") or machine


# NEW tests


def test_platform_tag_lowercase():
    tag = platform_tag()
    assert tag == tag.lower()


def test_platform_tag_no_spaces():
    assert " " not in platform_tag()


def test_binary_filename_ascii_only():
    name = binary_filename()
    assert name.isascii()


def test_platform_tag_ascii_only():
    assert platform_tag().isascii()
