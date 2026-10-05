#!/usr/bin/env python3
"""Run terser on a source string.

terser is a Node.js CLI. If it is not on PATH, this module reports
that fact so run_all.py can skip it.
"""

from __future__ import annotations

import shutil
import subprocess
import time
from typing import Tuple
# --- prepend bench/node_modules/.bin to PATH -------------------
import os
from pathlib import Path

_NODE_BIN = Path(__file__).resolve().parent / "node_modules" / ".bin"
if _NODE_BIN.is_dir():
    os.environ["PATH"] = str(_NODE_BIN) + os.pathsep + os.environ.get("PATH", "")
# ---------------------------------------------------------------


def is_available() -> bool:
    return shutil.which("terser") is not None


def run(source: str) -> str:
    if not is_available():
        raise RuntimeError("terser CLI not found on PATH")
    proc = subprocess.run(
        ["terser", "--compress", "--mangle"],
        input=source.encode("utf-8"),
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
    )
    if proc.returncode != 0:
        raise RuntimeError(f"terser failed: {proc.stderr.decode()}")
    return proc.stdout.decode("utf-8")


def run_timed(source: str) -> Tuple[str, float]:
    start = time.perf_counter()
    out = run(source)
    return out, time.perf_counter() - start


if __name__ == "__main__":
    import sys
    if len(sys.argv) < 2:
        print("usage: compare_terser.py INPUT.js", file=sys.stderr)
        raise SystemExit(2)
    src = open(sys.argv[1], "r", encoding="utf-8").read()
    out, t = run_timed(src)
    print(f"in:  {len(src)} bytes")
    print(f"out: {len(out)} bytes")
    print(f"t:   {t * 1000:.2f} ms")