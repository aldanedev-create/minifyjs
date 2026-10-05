#!/usr/bin/env python3
"""Minify a single file from Python.

This is the smallest useful Python program that uses MinifyJS. It
reads input.js, runs it through the full optimizer, writes
input.min.js, and prints the compression ratio.

Run:

    python app.py
"""

from __future__ import annotations

from pathlib import Path

from minifyjs import optimize, minify


HERE = Path(__file__).parent


def main() -> int:
    source = (HERE / "input.js").read_text(encoding="utf-8")

    print(f"input:  {len(source.encode('utf-8'))} bytes")

    # Just remove whitespace and comments.
    basic = minify(source)
    print(f"minify: {basic.minified_bytes} bytes "
          f"({basic.ratio * 100:.1f}% of input)")

    # Full optimizer: also mangles identifiers and folds constants.
    full = optimize(source)
    print(f"optim:  {full.minified_bytes} bytes "
          f"({full.ratio * 100:.1f}% of input, "
          f"saved {full.bytes_saved} bytes)")

    out_path = HERE / "input.min.js"
    out_path.write_text(full.code, encoding="utf-8")
    print(f"wrote {out_path.name}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())