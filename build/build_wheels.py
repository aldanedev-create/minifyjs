#!/usr/bin/env python3
"""Assemble Python wheels with a precompiled binary inside.

Each platform gets its own wheel. The wheel contains:

    minifyjs/                  the Python package (same on every platform)
    minifyjs/bin/minifyjs      the native binary for this platform
    minifyjs-<ver>.dist-info/  metadata

The metadata's WHEEL file declares the platform tag, so pip will only
install a wheel whose tag matches the user's machine.

Usage:

    python build_wheels.py --platform linux_x86_64 --version 1.2.3 --out wheelhouse
    python build_wheels.py --all --version 1.2.3 --out wheelhouse

The --all mode expects binaries at build/dist/<wheel_tag>/minifyjs
as produced by cross_compile.sh.
"""

from __future__ import annotations

import argparse
import base64
import csv
import hashlib
import io
import json
import os
import shutil
import subprocess
import sys
import zipfile
from pathlib import Path
from typing import Any, Dict, List, Optional


REPO_ROOT = Path(__file__).resolve().parent.parent
PYTHON_PKG = REPO_ROOT / "python" / "minifyjs"
PLATFORMS_FILE = Path(__file__).resolve().parent / "platforms.json"
DIST_DIR = Path(__file__).resolve().parent / "dist"

# Wheel metadata. The version comes from the caller.
DIST_NAME = "minifyjs"
PYTHON_TAG = "py3"
ABI_TAG = "none"


def load_platforms() -> List[Dict[str, Any]]:
    with PLATFORMS_FILE.open("r", encoding="utf-8") as f:
        return json.load(f)["platforms"]


def wheel_filename(version: str, wheel_tag: str) -> str:
    return f"{DIST_NAME}-{version}-{PYTHON_TAG}-{ABI_TAG}-{wheel_tag}.whl"


def wheel_metadata(version: str) -> str:
    return (
        f"Metadata-Version: 2.1\n"
        f"Name: {DIST_NAME}\n"
        f"Version: {version}\n"
        f"Summary: A native JavaScript minifier and optimizer. No Node.js required.\n"
        f"License: MIT\n"
        f"Requires-Python: >=3.8\n"
    )


def wheel_wheel_file(wheel_tag: str) -> str:
    return (
        f"Wheel-Version: 1.0\n"
        f"Generator: minifyjs-build {1}\n"
        f"Root-Is-Purelib: false\n"
        f"Tag: {PYTHON_TAG}-{ABI_TAG}-{wheel_tag}\n"
    )


def record_hash(data: bytes) -> str:
    """Return the base64url-encoded SHA-256 of data, without padding."""
    digest = hashlib.sha256(data).digest()
    return base64.urlsafe_b64encode(digest).rstrip(b"=").decode("ascii")


def iter_python_files() -> List[Path]:
    """Every file in python/minifyjs except bin/."""
    files: List[Path] = []
    for path in PYTHON_PKG.rglob("*"):
        if path.is_file() and "bin" not in path.relative_to(PYTHON_PKG).parts:
            files.append(path)
    return files


def build_wheel(
    version: str,
    platform: Dict[str, Any],
    binary_path: Path,
    out_dir: Path,
) -> Path:
    wheel_tag = platform["wheel_tag"]
    wheel_name = wheel_filename(version, wheel_tag)
    wheel_path = out_dir / wheel_name
    dist_info = f"{DIST_NAME}-{version}.dist-info"

    # (arcname, bytes) list to write.
    entries: List[tuple] = []

    # Python package files.
    for src in iter_python_files():
        rel = src.relative_to(PYTHON_PKG.parent)
        arcname = str(rel).replace(os.sep, "/")
        entries.append((arcname, src.read_bytes()))

    # The compiled binary.
    bin_arcname = f"minifyjs/bin/{platform['binary_name']}"
    entries.append((bin_arcname, binary_path.read_bytes()))

    # Metadata.
    entries.append((f"{dist_info}/METADATA", wheel_metadata(version).encode("utf-8")))
    entries.append((f"{dist_info}/WHEEL", wheel_wheel_file(wheel_tag).encode("utf-8")))

    # Build the RECORD file, which lists every other file with its
    # hash and size. RECORD itself has an empty hash and size.
    record_rows: List[List[str]] = []
    for arcname, data in entries:
        record_rows.append([
            arcname,
            f"sha256={record_hash(data)}",
            str(len(data)),
        ])

    record_arcname = f"{dist_info}/RECORD"
    record_rows.append([record_arcname, "", ""])

    record_buf = io.StringIO()
    record_writer = csv.writer(record_buf, lineterminator="\n")
    for row in record_rows:
        record_writer.writerow(row)
    record_bytes = record_buf.getvalue().encode("utf-8")

    # Write the wheel (which is just a zip).
    out_dir.mkdir(parents=True, exist_ok=True)
    with zipfile.ZipFile(wheel_path, "w", zipfile.ZIP_DEFLATED) as zf:
        for arcname, data in entries:
            zf.writestr(arcname, data)
        zf.writestr(record_arcname, record_bytes)

    return wheel_path


def main(argv: Optional[List[str]] = None) -> int:
    parser = argparse.ArgumentParser(description="Build MinifyJS Python wheels")
    parser.add_argument("--platform", help="wheel_tag to build")
    parser.add_argument("--all", action="store_true", help="build every platform")
    parser.add_argument("--version", required=True, help="version string")
    parser.add_argument("--out", required=True, help="output wheelhouse directory")
    parser.add_argument(
        "--dist", default=str(DIST_DIR),
        help="directory holding prebuilt binaries (default: build/dist)",
    )
    args = parser.parse_args(argv)

    if not args.platform and not args.all:
        parser.error("one of --platform or --all is required")

    platforms = load_platforms()
    out_dir = Path(args.out).resolve()
    dist_dir = Path(args.dist).resolve()

    if not PYTHON_PKG.is_dir():
        raise SystemExit(f"python package not found at {PYTHON_PKG}")

    targets = platforms if args.all else [p for p in platforms if p["wheel_tag"] == args.platform]
    if not targets:
        raise SystemExit(f"unknown platform: {args.platform}")

    for p in targets:
        binary_path = dist_dir / p["wheel_tag"] / p["binary_name"]
        if not binary_path.is_file():
            raise SystemExit(
                f"binary missing for {p['wheel_tag']}: {binary_path}\n"
                f"run build_core.py --platform {p['wheel_tag']} --out {dist_dir} first"
            )

        wheel_path = build_wheel(args.version, p, binary_path, out_dir)
        print(f"built {wheel_path.name} ({wheel_path.stat().st_size} bytes)")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())