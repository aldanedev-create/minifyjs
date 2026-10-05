"""Maps the running interpreter's OS and architecture to the name of
the precompiled native binary the package should use.

The Go engine is compiled once per (OS, architecture) pair during the
release process; each platform-specific wheel bundles exactly one
binary. This module is what the loader consults to find it.
"""

from __future__ import annotations

import platform
import sys


def binary_filename() -> str:
    """Return the expected binary filename for the current platform."""
    return "minifyjs.exe" if sys.platform.startswith("win") else "minifyjs"


def platform_tag() -> str:
    """Return a human-readable platform tag, used in error messages."""
    system = platform.system().lower()
    machine = platform.machine().lower()
    if system == "darwin":
        system = "macos"
    if machine in ("x86_64", "amd64"):
        machine = "x86_64"
    elif machine in ("arm64", "aarch64"):
        machine = "arm64"
    return f"{system}-{machine}"
