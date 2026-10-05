#!/usr/bin/env python3
"""Measure compression ratio of each minifier on each fixture.

Output is a markdown table suitable for pasting into a README.

    python measure_compression.py
"""

from __future__ import annotations

import sys
from pathlib import Path

HERE = Path(__file__).parent
FIXTURES = HERE / "fixtures"

sys.path.insert(0, str(HERE))

from compare_minifyjs import run as run_minifyjs, is_available as has_minifyjs
from compare_esbuild import run as run_esbuild, is_available as has_esbuild
from compare_terser import run as run_terser, is_available as has_terser
from compare_uglifyjs import run as run_uglifyjs, is_available as has_uglifyjs
from compare_rjsmin import run as run_rjsmin, is_available as has_rjsmin


TOOLS = [
    ("minifyjs", run_minifyjs, has_minifyjs),
    ("esbuild", run_esbuild, has_esbuild),
    ("terser", run_terser, has_terser),
    ("uglify-js", run_uglifyjs, has_uglifyjs),
    ("rjsmin", run_rjsmin, has_rjsmin),
]


def main() -> int:
    fixture_paths = sorted(FIXTURES.glob("*.js"))
    if not fixture_paths:
        print("no fixtures found", file=sys.stderr)
        return 1

    active = [(name, fn) for name, fn, avail in TOOLS if avail()]

    # Header
    print("| Fixture | Input |", end="")
    for name, _ in active:
        print(f" {name} |", end="")
    print()

    print("|---|---|", end="")
    for _ in active:
        print("---|", end="")
    print()

    # Rows
    for path in fixture_paths:
        src = path.read_text(encoding="utf-8")
        in_size = len(src.encode("utf-8"))
        print(f"| `{path.name}` | {in_size} B |", end="")

        for name, fn in active:
            try:
                out = fn(src)
                out_size = len(out.encode("utf-8"))
                pct = (1 - out_size / in_size) * 100 if in_size else 0
                print(f" {out_size} B ({pct:.1f}%) |", end="")
            except Exception as e:
                print(f" error ({e}) |", end="")
        print()

    return 0


if __name__ == "__main__":
    raise SystemExit(main())