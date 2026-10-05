#!/usr/bin/env python3
"""Verify a checksums.txt against the files it names.

The file format is sha256sum's:

    <hex>  <filename>

Filenames are resolved relative to the checksums file's directory.

Usage:

    python tools/release/verify_checksums.py build/wheelhouse/checksums.txt
"""

from __future__ import annotations

import argparse
import hashlib
import sys
from pathlib import Path


def sha256_of(path: Path) -> str:
    h = hashlib.sha256()
    with path.open("rb") as f:
        for chunk in iter(lambda: f.read(65536), b""):
            h.update(chunk)
    return h.hexdigest()


def main() -> int:
    parser = argparse.ArgumentParser(description="Verify SHA-256 checksums")
    parser.add_argument("checksums")
    args = parser.parse_args()

    checksums = Path(args.checksums).resolve()
    if not checksums.is_file():
        print(f"error: {checksums} not found", file=sys.stderr)
        return 1

    base = checksums.parent
    checked = 0
    failed = 0

    for lineno, line in enumerate(checksums.read_text(encoding="utf-8").splitlines(), 1):
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        parts = line.split(None, 1)
        if len(parts) != 2:
            print(f"line {lineno}: malformed: {line!r}", file=sys.stderr)
            failed += 1
            continue
        expected_hash, filename = parts
        target = base / filename
        if not target.is_file():
            print(f"line {lineno}: missing file: {filename}", file=sys.stderr)
            failed += 1
            continue
        actual_hash = sha256_of(target)
        if actual_hash != expected_hash:
            print(f"line {lineno}: hash mismatch for {filename}", file=sys.stderr)
            failed += 1
            continue
        checked += 1

    if failed:
        print(f"checksum verification FAILED ({failed} of {checked + failed})",
              file=sys.stderr)
        return 1

    print(f"checksum verification OK ({checked} file(s))")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())