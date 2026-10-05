#!/usr/bin/env python3
"""Write a JSON manifest describing every artifact in a release.

The manifest is attached to the GitHub release page so that a
downstream consumer (or a mirror) can enumerate what was published
without parsing the release notes.

The manifest contains, for each file:

  - its path relative to the release directory
  - its byte size
  - its SHA-256 checksum
  - its platform tag, if the filename ends with a wheel tag

Usage:

    python generate_manifest.py --dir wheelhouse --out wheelhouse/manifest.json --version 0.1.0
"""

from __future__ import annotations

import argparse
import hashlib
import json
from datetime import datetime, timezone
from pathlib import Path
from typing import Dict, List, Optional


def sha256_of(path: Path) -> str:
    h = hashlib.sha256()
    with path.open("rb") as f:
        for chunk in iter(lambda: f.read(65536), b""):
            h.update(chunk)
    return h.hexdigest()


def guess_platform(filename: str) -> Optional[str]:
    for tag in (
        "manylinux2014_x86_64",
        "manylinux2014_aarch64",
        "macosx_10_13_x86_64",
        "macosx_11_0_arm64",
        "win_amd64",
    ):
        if tag in filename:
            return tag
    return None


def main(argv: Optional[List[str]] = None) -> int:
    parser = argparse.ArgumentParser(description="Generate release manifest")
    parser.add_argument("--dir", required=True)
    parser.add_argument("--out", required=True)
    parser.add_argument("--version", default="unknown")
    args = parser.parse_args(argv)

    root = Path(args.dir).resolve()
    out = Path(args.out).resolve()

    if not root.is_dir():
        raise SystemExit(f"not a directory: {root}")

    entries: List[Dict[str, object]] = []
    for path in sorted(root.rglob("*")):
        if not path.is_file():
            continue
        if path.resolve() == out:
            continue
        rel = str(path.relative_to(root)).replace("\\", "/")
        entries.append({
            "path": rel,
            "size": path.stat().st_size,
            "sha256": sha256_of(path),
            "platform": guess_platform(path.name),
        })

    manifest = {
        "name": "minifyjs",
        "version": args.version,
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "files": entries,
    }

    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(manifest, indent=2), encoding="utf-8")
    print(f"wrote {out} ({len(entries)} files)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())