#!/usr/bin/env python3
"""Run MinifyJS on a source string and return the minified output.

Used by run_all.py to include MinifyJS in head-to-head comparisons.
This module uses the Python package, not the CLI, because the Python
package is what most users of this project will invoke.
"""

from __future__ import annotations

import time
from typing import Optional, Tuple

from minifyjs import optimize


def run(source: str, *, compress: bool = True, mangle: bool = True) -> str:
    """Minify source and return the output string."""
    if compress and mangle:
        return optimize(source).code
    from minifyjs import minify
    return minify(source, compress=compress, mangle=mangle).code


def run_timed(source: str) -> Tuple[str, float]:
    """Minify source, returning (output, elapsed_seconds)."""
    start = time.perf_counter()
    out = run(source)
    elapsed = time.perf_counter() - start
    return out, elapsed


if __name__ == "__main__":
    import sys
    if len(sys.argv) < 2:
        print("usage: compare_minifyjs.py INPUT.js", file=sys.stderr)
        raise SystemExit(2)
    src = open(sys.argv[1], "r", encoding="utf-8").read()
    out, t = run_timed(src)
    print(f"in:  {len(src)} bytes")
    print(f"out: {len(out)} bytes")
    print(f"t:   {t * 1000:.2f} ms")