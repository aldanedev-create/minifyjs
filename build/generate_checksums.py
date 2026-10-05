#!/usr/bin/env python3
"""Write SHA-256 checksums of every file in a directory.

Output is the standard `sha256sum` format:

    <hash>  <filename>

so the file can be verified with `sha256sum -c checksums.txt`.

Usage:

    python generate_checksums.py --dir wheelhouse --out wheelhouse/checksums.txt
"""

from __future__ import annotations

import argparse
import hashlib
from pathlib import Path
from typing import List, Optional


def sha256_of(path: Path) -> str:
    h = hashlib.sha256()
    with path.open("rb") as f:
        for chunk in iter(lambda: f.read(65536), b""):
            h.update(chunk)
    return h.hexdigest()


def main(argv: Optional[List[str]] = None) -> int:
    parser = argparse.ArgumentParser(description="Generate SHA-256 checksums")
    parser.add_argument("--dir", required=True, help="directory to walk")
    parser.add_argument("--out", required=True, help="output file path")
    parser.add_argument("--exclude", action="append", default=[],
                        help="filename to skip (repeatable)")
    args = parser.parse_args(argv)

    root = Path(args.dir).resolve()
    out = Path(args.out).resolve()

    if not root.is_dir():
        raise SystemExit(f"not a directory: {root}")

    lines: List[str] = []
    for path in sorted(root.rglob("*")):
        if not path.is_file():
            continue
        if path.resolve() == out:
            continue
        if path.name in args.exclude:
            continue
        rel = path.relative_to(root)
        lines.append(f"{sha256_of(path)}  {rel}")

    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text("\n".join(lines) + "\n", encoding="utf-8")
    print(f"wrote {out} ({len(lines)} entries)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())