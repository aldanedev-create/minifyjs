#!/usr/bin/env python3
"""Measure throughput in bytes per second for each minifier.

Runs each tool N times on a fixture and reports the median.

    python measure_throughput.py [--iterations 10]
"""

from __future__ import annotations

import argparse
import statistics
import sys
from pathlib import Path
from typing import Callable, List

HERE = Path(__file__).parent
FIXTURES = HERE / "fixtures"

sys.path.insert(0, str(HERE))

from compare_minifyjs import run_timed as minifyjs_timed, is_available as has_minifyjs
from compare_esbuild import run_timed as esbuild_timed, is_available as has_esbuild
from compare_terser import run_timed as terser_timed, is_available as has_terser
from compare_uglifyjs import run_timed as uglifyjs_timed, is_available as has_uglifyjs
from compare_rjsmin import run_timed as rjsmin_timed, is_available as has_rjsmin


TOOLS: List[tuple] = [
    ("minifyjs", minifyjs_timed, has_minifyjs),
    ("esbuild", esbuild_timed, has_esbuild),
    ("terser", terser_timed, has_terser),
    ("uglify-js", uglifyjs_timed, has_uglifyjs),
    ("rjsmin", rjsmin_timed, has_rjsmin),
]


def measure(fn: Callable, src: str, iterations: int) -> float:
    times: List[float] = []
    for _ in range(iterations):
        _, t = fn(src)
        times.append(t)
    return statistics.median(times)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--iterations", type=int, default=10)
    parser.add_argument("--fixture", default="medium.js")
    args = parser.parse_args()

    path = FIXTURES / args.fixture
    if not path.is_file():
        print(f"fixture not found: {path}", file=sys.stderr)
        return 1

    src = path.read_text(encoding="utf-8")
    size = len(src.encode("utf-8"))
    print(f"fixture: {path.name} ({size} B)")
    print(f"iterations: {args.iterations}")
    print()

    print(f"{'tool':<12} {'median ms':>12} {'MB/s':>10}")
    print("-" * 36)

    for name, fn, avail in TOOLS:
        if not avail():
            print(f"{name:<12} {'(not installed)':>12}")
            continue
        try:
            median_s = measure(fn, src, args.iterations)
            mb_per_s = (size / 1_000_000) / median_s if median_s else 0
            print(f"{name:<12} {median_s * 1000:>12.2f} {mb_per_s:>10.1f}")
        except Exception as e:
            print(f"{name:<12} {'error: ' + str(e)[:30]:>12}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())