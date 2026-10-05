#!/usr/bin/env python3
"""Compile the MinifyJS binary for a single platform.

This is the primitive that build_wheels.py and cross_compile.sh both
rely on. It shells out to `go build` with the right GOOS, GOARCH, and
ldflags, and writes the resulting binary to a known path.

Usage:

    python build_core.py --platform linux-x86_64 --out /tmp/minifyjs
    python build_core.py --all --out build/dist

The script does not touch the network. If `go` is not on PATH, it
fails with a clear message.
"""

from __future__ import annotations

import argparse
import json
import os
import shutil
import subprocess
import sys
from pathlib import Path
from typing import Any, Dict, List, Optional


REPO_ROOT = Path(__file__).resolve().parent.parent
CORE_DIR = REPO_ROOT / "core"
PLATFORMS_FILE = Path(__file__).resolve().parent / "platforms.json"


def load_platforms() -> List[Dict[str, Any]]:
    with PLATFORMS_FILE.open("r", encoding="utf-8") as f:
        return json.load(f)["platforms"]


def platform_by_wheel_tag(platforms: List[Dict[str, Any]], tag: str) -> Dict[str, Any]:
    for p in platforms:
        if p["wheel_tag"] == tag:
            return p
        if p["human"].lower().replace(" ", "-") == tag.lower():
            return p
    raise SystemExit(f"unknown platform: {tag}")


def check_go_available() -> None:
    if shutil.which("go") is None:
        raise SystemExit(
            "go is not on PATH. Install Go from https://go.dev/dl/ "
            "or run this script inside the build container "
            "(build/Dockerfile.build)."
        )


def get_version() -> str:
    """Return the version string to bake into the binary."""
    try:
        out = subprocess.run(
            ["git", "describe", "--tags", "--always", "--dirty"],
            cwd=REPO_ROOT,
            capture_output=True,
            text=True,
            check=True,
        )
        return out.stdout.strip() or "dev"
    except (subprocess.CalledProcessError, FileNotFoundError):
        return "dev"


def build_one(platform: Dict[str, Any], out_dir: Path, version: str) -> Path:
    out_dir.mkdir(parents=True, exist_ok=True)
    out_file = out_dir / platform["binary_name"]

    env = os.environ.copy()
    env["GOOS"] = platform["goos"]
    env["GOARCH"] = platform["goarch"]
    env["CGO_ENABLED"] = "0"

    ldflags = (
        f"-s -w "
        f"-X github.com/minifyjs/minifyjs/core/internal/version.Version={version}"
    )

    cmd = [
        "go", "build",
        "-trimpath",
        "-ldflags", ldflags,
        "-o", str(out_file),
        "./cmd/minifyjs",
    ]

    print(f"  building {platform['human']} ({platform['goos']}/{platform['goarch']})")
    subprocess.run(cmd, cwd=CORE_DIR, env=env, check=True)

    return out_file


def main(argv: Optional[List[str]] = None) -> int:
    parser = argparse.ArgumentParser(description="Build MinifyJS for one or all platforms")
    parser.add_argument("--platform", help="wheel_tag to build")
    parser.add_argument("--all", action="store_true", help="build every platform")
    parser.add_argument("--out", required=True, help="output directory")
    parser.add_argument("--version", default=None, help="version string (default: git describe)")
    args = parser.parse_args(argv)

    if not args.platform and not args.all:
        parser.error("one of --platform or --all is required")

    check_go_available()

    platforms = load_platforms()
    version = args.version or get_version()
    out_dir = Path(args.out).resolve()

    print(f"MinifyJS build ({version})")

    if args.all:
        for p in platforms:
            build_one(p, out_dir / p["wheel_tag"], version)
    else:
        p = platform_by_wheel_tag(platforms, args.platform)
        build_one(p, out_dir, version)

    print(f"wrote binaries under {out_dir}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())