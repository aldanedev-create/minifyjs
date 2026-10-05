#!/usr/bin/env python3
"""Bundle a multi-file project from Python.

The Python API's `bundle()` function wraps esbuild's Build API. It
resolves imports, applies tree shaking, and writes one output file
per entry point.

Run:

    python main.py
"""

from __future__ import annotations

import shutil
from pathlib import Path

from minifyjs import bundle


HERE = Path(__file__).parent
SRC = HERE / "src"
DIST = HERE / "dist"


def main() -> int:
    # Clean the output directory so re-runs are deterministic.
    if DIST.exists():
        shutil.rmtree(DIST)
    DIST.mkdir(parents=True)

    print(f"bundling {SRC / 'main.js'} ...")

    result = bundle(
        entry_points=[str(SRC / "main.js")],
        outfile=str(DIST / "bundle.js"),
        format="esm",
        target="es2015",
        compress=True,
        mangle=True,
        sourcemap="external",
    )

    if result.has_errors:
        for d in result.diagnostics:
            print(f"error: {d}")
        return 1

    bundle_size = (DIST / "bundle.js").stat().st_size
    map_size = (DIST / "bundle.js.map").stat().st_size

    print(f"wrote {DIST / 'bundle.js'} ({bundle_size} bytes)")
    print(f"wrote {DIST / 'bundle.js.map'} ({map_size} bytes)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())