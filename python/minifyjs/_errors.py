"""Internal error-message builders shared by _binary.py and _runner.py.

Public exception types live in errors.py. This module only formats
strings so the CLI and the Python API can present the same
diagnostics.
"""

from __future__ import annotations


def binary_not_found_message(path: str, platform_tag: str) -> str:
    return (
        f"minifyjs: no native binary found at {path!r} for platform "
        f"{platform_tag!r}.\n"
        "This normally means the platform-specific wheel for your OS/"
        "architecture hasn't been built yet, or you are running from a "
        "source checkout without having built core/cmd/minifyjs. See "
        "docs/getting-started/installation.md."
    )


def binary_not_executable_message(path: str) -> str:
    return (
        f"minifyjs: native binary at {path!r} is not executable and "
        "could not be made executable. Check the file permissions."
    )
