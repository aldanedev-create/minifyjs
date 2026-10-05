"""Resolves the on-disk location of the bundled native binary."""

from __future__ import annotations

import os

from ._platform import binary_filename


def bin_dir() -> str:
    """Return the directory the package stores its native binary in."""
    return os.path.join(os.path.dirname(__file__), "bin")


def binary_path() -> str:
    """Return the expected full path of the native binary."""
    return os.path.join(bin_dir(), binary_filename())
