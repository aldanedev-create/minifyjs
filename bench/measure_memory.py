#!/usr/bin/env python3
"""Measure peak RSS for each minifier on the large fixture.

    python measure_memory.py
"""

from __future__ import annotations

import shutil
import subprocess
import sys
from pathlib import Path

try:
    import psutil
except ImportError:
    print("psutil is required: pip install psutil", file=sys.stderr)
    raise SystemExit(1)

HERE = Path(__file__).parent
FIXTURES = HERE / "fixtures"


def peak_rss_of(cmd, input_bytes: bytes) -> int:
    """Run cmd with input_bytes on stdin and return peak RSS in bytes."""
    proc = subprocess.Popen(
        cmd,
        stdin=subprocess.PIPE,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )
    proc.stdin.write(input_bytes)
    proc.stdin.close()

    p = psutil.Process(proc.pid)
    peak = 0
    while proc.poll() is None:
        try:
            info = p.memory_info()
            if info.rss > peak:
                peak = info.rss
        except psutil.NoSuchProcess:
            break
    proc.wait()
    return peak


def main() -> int:
    fixture = FIXTURES / "large.js"
    if not fixture.is_file():
        fixture = FIXTURES / "medium.js"
    if not fixture.is_file():
        print("no fixture found", file=sys.stderr)
        return 1

    src = fixture.read_bytes()
    print(f"fixture: {fixture.name} ({len(src)} B)")
    print()
    print(f"{'tool':<12} {'peak RSS':>12}")
    print("-" * 26)

    tools = [
        ("minifyjs", ["minifyjs"]),
        ("esbuild", ["esbuild", "--minify", "--log-level=silent"]),
        ("terser", ["terser", "--compress", "--mangle"]),
        ("uglify-js", ["uglifyjs", "--compress", "--mangle"]),
    ]

    for name, cmd in tools:
        if shutil.which(cmd[0]) is None:
            print(f"{name:<12} {'(not installed)':>12}")
            continue
        peak = peak_rss_of(cmd, src)
        print(f"{name:<12} {peak / 1_000_000:>9.1f} MB")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())