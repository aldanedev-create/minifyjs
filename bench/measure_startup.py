#!/usr/bin/env python3
"""Measure process startup time for each minifier.

Startup time matters when a build script invokes a minifier once per
file. It is the dominant cost for small inputs.

    python measure_startup.py [--iterations 20]
"""

from __future__ import annotations

import argparse
import shutil
import statistics
import subprocess
import sys
import time
from pathlib import Path
from typing import List, Optional

HERE = Path(__file__).parent


def timed_run(cmd: List[str], stdin: bytes = b"") -> float:
    start = time.perf_counter()
    subprocess.run(
        cmd,
        input=stdin,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
        check=False,
    )
    return time.perf_counter() - start


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--iterations", type=int, default=20)
    args = parser.parse_args()

    src = b"const x = 1;\n"

    entries = [
        ("minifyjs", ["minifyjs"]),
        ("esbuild", ["esbuild", "--minify", "--log-level=silent"]),
        ("terser", ["terser", "--compress", "--mangle"]),
        ("uglify-js", ["uglifyjs", "--compress", "--mangle"]),
        ("rjsmin", None),  # handled specially, it is not a CLI
    ]

    print(f"iterations: {args.iterations}")
    print()
    print(f"{'tool':<12} {'median ms':>12}")
    print("-" * 26)

    for name, cmd in entries:
        if cmd is None:
            # rjsmin is a Python library, not a CLI. Measure import.
            times: List[float] = []
            for _ in range(args.iterations):
                start = time.perf_counter()
                subprocess.run(
                    [sys.executable, "-c", "import rjsmin"],
                    stdout=subprocess.DEVNULL,
                    stderr=subprocess.DEVNULL,
                    check=False,
                )
                times.append(time.perf_counter() - start)
            print(f"{name:<12} {statistics.median(times) * 1000:>12.2f}")
            continue

        if shutil.which(cmd[0]) is None:
            print(f"{name:<12} {'(not installed)':>12}")
            continue

        times = [timed_run(cmd, src) for _ in range(args.iterations)]
        print(f"{name:<12} {statistics.median(times) * 1000:>12.2f}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())