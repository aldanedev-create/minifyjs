#!/usr/bin/env python3
"""Run rjsmin on a source string.

rjsmin is a pure-Python JavaScript minifier. It is much smaller than
the other tools and is the baseline most Python users reach for
today, so it is useful for comparison even though it does far less.
"""

from __future__ import annotations

import time
from typing import Tuple


def is_available() -> bool:
    try:
        import rjsmin  # noqa: F401
        return True
    except ImportError:
        return False


def run(source: str) -> str:
    if not is_available():
        raise RuntimeError("rjsmin not installed (pip install rjsmin)")
    import rjsmin
    return rjsmin.jsmin(source)


def run_timed(source: str) -> Tuple[str, float]:
    start = time.perf_counter()
    out = run(source)
    return out, time.perf_counter() - start


if __name__ == "__main__":
    import sys
    if len(sys.argv) < 2:
        print("usage: compare_rjsmin.py INPUT.js", file=sys.stderr)
        raise SystemExit(2)
    src = open(sys.argv[1], "r", encoding="utf-8").read()
    out, t = run_timed(src)
    print(f"in:  {len(src)} bytes")
    print(f"out: {len(out)} bytes")
    print(f"t:   {t * 1000:.2f} ms")