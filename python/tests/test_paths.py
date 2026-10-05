"""Tests for the path resolution helpers."""

from __future__ import annotations

import os

from minifyjs._paths import bin_dir, binary_path


def test_bin_dir_is_absolute():
    assert os.path.isabs(bin_dir())


def test_bin_dir_ends_with_bin():
    assert bin_dir().endswith("bin")


def test_bin_dir_is_directory():
    assert os.path.isdir(bin_dir())


def test_binary_path_is_absolute():
    assert os.path.isabs(binary_path())


def test_binary_path_is_inside_bin_dir():
    assert binary_path().startswith(bin_dir())


def test_binary_path_is_string():
    assert isinstance(binary_path(), str)
