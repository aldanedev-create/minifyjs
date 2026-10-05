"""Python-side CLI entry point.

``pip install minifyjs`` installs a ``minifyjs`` console script that
points here. This module forwards argv to the bundled native binary
so the documented CLI behavior is identical whether the binary is
invoked directly or through ``python -m minifyjs``.
"""

from __future__ import annotations

import subprocess
import sys

from ._binary import find_binary


def main() -> int:
    """Run the native CLI with the current argv. Returns its exit code."""
    binary = find_binary()
    proc = subprocess.run([binary, *sys.argv[1:]])
    return proc.returncode


if __name__ == "__main__":
    raise SystemExit(main())
