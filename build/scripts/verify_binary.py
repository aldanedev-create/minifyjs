#!/usr/bin/env python3
"""Smoke-test a compiled MinifyJS binary.

Given a path to a `minifyjs` binary, this script:

  1. Runs `minifyjs --version` and checks the output looks right.
  2. Pipes a small JavaScript file through stdin and checks the output.
  3. Runs a syntax-error input and checks the exit code.

It is used by the release pipeline to confirm that a cross-compiled
binary actually works before it is packaged into a wheel.

Usage:

    python verify_binary.py /path/to/minifyjs
"""

from __future__ import annotations

import argparse
import subprocess
import sys
from pathlib import Path
from typing import List, Optional


def run(binary: Path, args: List[str], stdin: bytes = b"") -> subprocess.CompletedProcess:
    return subprocess.run(
        [str(binary), *args],
        input=stdin,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
    )


def fail(msg: str) -> None:
    print(f"FAIL: {msg}", file=sys.stderr)
    sys.exit(1)


def main(argv: Optional[List[str]] = None) -> int:
    parser = argparse.ArgumentParser(description="Smoke-test a MinifyJS binary")
    parser.add_argument("binary", help="path to the minifyjs executable")
    args = parser.parse_args(argv)

    binary = Path(args.binary).resolve()
    if not binary.is_file():
        fail(f"binary not found: {binary}")

    # Test 1: --version.
    proc = run(binary, ["--version"])
    if proc.returncode != 0:
        fail(f"--version exited {proc.returncode}: {proc.stderr!r}")
    out = proc.stdout.decode("utf-8", "replace")
    if "minifyjs" not in out:
        fail(f"--version output missing 'minifyjs': {out!r}")
    if "esbuild" not in out:
        fail(f"--version output missing 'esbuild': {out!r}")

    # Test 2: minify a small input.
    proc = run(binary, [], b"function add(a, b) { return a + b; }\n")
    if proc.returncode != 0:
        fail(f"minify exited {proc.returncode}: {proc.stderr!r}")
    if b"function add(a,b){return a+b;}" not in proc.stdout:
        fail(f"unexpected minified output: {proc.stdout!r}")

    # Test 3: syntax error.
    proc = run(binary, [], b"function () { } }")
    if proc.returncode == 0:
        fail("syntax error input exited 0")
    if not proc.stderr:
        fail("syntax error produced no stderr")

    # Test 4: --help.
    proc = run(binary, ["--help"])
    if proc.returncode != 0:
        fail(f"--help exited {proc.returncode}")
    if b"Usage:" not in proc.stdout:
        fail("--help output missing 'Usage:'")

    print(f"OK: {binary}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())