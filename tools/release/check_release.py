#!/usr/bin/env python3
"""Pre-release checklist.

Run this before tagging a release. It verifies:

  1. Version is consistent across every file that declares one.
  2. CHANGELOG.md has an entry for the release version.
  3. The tag format matches (vMAJOR.MINOR.PATCH).
  4. The working tree is clean.
  5. The current branch is main.

Exits 0 on success, 1 on any failure. Prints every failure it finds
rather than stopping at the first, so a maintainer can fix them all
in one pass.

Usage:

    python tools/release/check_release.py --version 0.1.0
"""

from __future__ import annotations

import argparse
import re
import subprocess
import sys
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[2]

# Files that declare a version, and the pattern to extract it.
VERSION_FILES = [
    ("core/internal/version/version.go", r'var Version = "([^"]+)"'),
    ("python/minifyjs/_version.py", r'__version__ = "([^"]+)"'),
    ("python/pyproject.toml", r'^version = "([^"]+)"'),
    ("bench/package.json", r'"version":\s*"([^"]+)"'),
]


def run(cmd, cwd=None):
    return subprocess.run(cmd, cwd=cwd, capture_output=True, text=True)


def check_versions(expected: str) -> list:
    problems = []
    for path, pattern in VERSION_FILES:
        full = REPO_ROOT / path
        if not full.is_file():
            continue
        text = full.read_text(encoding="utf-8")
        m = re.search(pattern, text, re.MULTILINE)
        if not m:
            problems.append(f"{path}: version not found")
            continue
        found = m.group(1)
        if found != expected:
            problems.append(f"{path}: version is {found}, expected {expected}")
    return problems


def check_changelog(expected: str) -> list:
    changelog = REPO_ROOT / "CHANGELOG.md"
    if not changelog.is_file():
        return ["CHANGELOG.md: file missing"]
    text = changelog.read_text(encoding="utf-8")
    if f"## [{expected}]" not in text:
        return [f"CHANGELOG.md: no section for [{expected}]"]
    return []


def check_tag_format(expected: str) -> list:
    if not re.fullmatch(r"\d+\.\d+\.\d+", expected):
        return [f"version {expected!r} is not MAJOR.MINOR.PATCH"]
    return []


def check_git_clean() -> list:
    result = run(["git", "status", "--porcelain"], cwd=REPO_ROOT)
    if result.returncode != 0:
        return ["git status failed (not a git repository?)"]
    if result.stdout.strip():
        lines = result.stdout.strip().splitlines()
        return [f"working tree is not clean ({len(lines)} changed file(s))"]
    return []


def check_branch() -> list:
    result = run(["git", "rev-parse", "--abbrev-ref", "HEAD"], cwd=REPO_ROOT)
    if result.returncode != 0:
        return ["could not determine current branch"]
    branch = result.stdout.strip()
    if branch != "main":
        return [f"current branch is {branch!r}, expected 'main'"]
    return []


def main() -> int:
    parser = argparse.ArgumentParser(description="Pre-release checklist")
    parser.add_argument("--version", required=True)
    parser.add_argument("--skip-git", action="store_true",
                        help="skip the git cleanliness and branch checks")
    args = parser.parse_args()

    problems = []
    problems += check_tag_format(args.version)
    problems += check_versions(args.version)
    problems += check_changelog(args.version)
    if not args.skip_git:
        problems += check_git_clean()
        problems += check_branch()

    if problems:
        print("release check FAILED:")
        for p in problems:
            print(f"  - {p}")
        return 1

    print(f"release check OK for {args.version}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())