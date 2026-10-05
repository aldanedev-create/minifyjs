#!/usr/bin/env python3
"""Verify that a staged release contains every platform wheel.

Reads build/platforms.json to get the expected list, then checks that
a wheel exists in the wheelhouse for each one.

Usage:

    python tools/release/verify_artifacts.py --wheelhouse build/wheelhouse --version 0.1.0
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[2]


def load_platforms() -> list:
    path = REPO_ROOT / "build" / "platforms.json"
    with path.open("r", encoding="utf-8") as f:
        return json.load(f)["platforms"]


def main() -> int:
    parser = argparse.ArgumentParser(description="Verify release artifacts")
    parser.add_argument("--wheelhouse", required=True)
    parser.add_argument("--version", required=True)
    args = parser.parse_args()

    wheelhouse = Path(args.wheelhouse).resolve()
    if not wheelhouse.is_dir():
        print(f"error: {wheelhouse} is not a directory", file=sys.stderr)
        return 1

    platforms = load_platforms()
    found = set()

    for whl in wheelhouse.glob("*.whl"):
        m = re.match(r"minifyjs-([^-]+)-py3-none-(.+)\.whl", whl.name)
        if not m:
            continue
        ver, tag = m.group(1), m.group(2)
        if ver != args.version:
            print(f"error: {whl.name} has wrong version {ver}",
                  file=sys.stderr)
            return 1
        found.add(tag)

    missing = [p["wheel_tag"] for p in platforms if p["wheel_tag"] not in found]

    if missing:
        print("release artifact check FAILED")
        print(f"  found:   {sorted(found)}")
        print(f"  missing: {missing}")
        return 1

    print(f"release artifact check OK ({len(found)} wheel(s))")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())