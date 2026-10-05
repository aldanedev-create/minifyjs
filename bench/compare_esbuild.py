#!/usr/bin/env python3
"""Run the esbuild CLI on a source string.

esbuild is a Node.js binary. If it is not on PATH, this module
reports that fact so run_all.py can skip it cleanly.
"""

from __future__ import annotations

import shutil
import subprocess
import time
from typing import Optional, Tuple
# --- prepend bench/node_modules/.bin to PATH -------------------
import os
from pathlib import Path

_NODE_BIN = Path(__file__).resolve().parent / "node_modules" / ".bin"
if _NODE_BIN.is_dir():
    os.environ["PATH"] = str(_NODE_BIN) + os.pathsep + os.environ.get("PATH", "")
# ---------------------------------------------------------------


def is_available() -> bool:
    return shutil.which("esbuild") is not None


def run(source: str) -> str:
    """Run esbuild's --minify on source and return the output."""
    if not is_available():
        raise RuntimeError("esbuild CLI not found on PATH")
    proc = subprocess.run(
        ["esbuild", "--minify", "--log-level=silent"],
        input=source.encode("utf-8"),
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
    )
    if proc.returncode != 0:
        raise RuntimeError(f"esbuild failed: {proc.stderr.decode()}")
    return proc.stdout.decode("utf-8")


def run_timed(source: str) -> Tuple[str, float]:
    start = time.perf_counter()
    out = run(source)
    return out, time.perf_counter() - start


if __name__ == "__main__":
    import sys
    if len(sys.argv) < 2:
        print("usage: compare_esbuild.py INPUT.js", file=sys.stderr)
        raise SystemExit(2)
    src = open(sys.argv[1], "r", encoding="utf-8").read()
    out, t = run_timed(src)
    print(f"in:  {len(src)} bytes")
    print(f"out: {len(out)} bytes")
    print(f"t:   {t * 1000:.2f} ms")