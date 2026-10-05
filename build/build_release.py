#!/usr/bin/env python3
"""The full release pipeline.

Runs, in order:

  1. Go tests
  2. Python tests
  3. Cross-compile binaries for every platform
  4. Verify each binary runs
  5. Build a wheel per platform
  6. Verify each wheel installs and runs
  7. Generate checksums and a release manifest
  8. (Optional) Publish to PyPI

Step 8 is only run when `--publish` is passed. Without it, artifacts
are staged in build/wheelhouse/ and left for review.

Usage:

    python build_release.py --version 0.1.0
    python build_release.py --version 0.1.0 --publish

The script is a thin orchestrator. It does not contain any build logic
itself; every step delegates to a script in this directory.
"""

from __future__ import annotations

import argparse
import shutil
import subprocess
import sys
from pathlib import Path
from typing import List, Optional


BUILD_DIR = Path(__file__).resolve().parent
REPO_ROOT = BUILD_DIR.parent
WHEELHOUSE = BUILD_DIR / "wheelhouse"
DIST = BUILD_DIR / "dist"


def run(cmd: List[str], cwd: Optional[Path] = None) -> None:
    print(f"$ {' '.join(cmd)}")
    subprocess.run(cmd, cwd=cwd, check=True)


def step_go_tests() -> None:
    print("==> step 1: go tests")
    run(["go", "test", "./..."], cwd=REPO_ROOT / "core")


def step_build_test_binary(version: str) -> None:
    """Refresh the source package binary before tests exercise new CLI flags."""
    import json
    platforms = json.loads((BUILD_DIR / "platforms.json").read_text())["platforms"]
    host = next((p for p in platforms if _is_host_platform(p["wheel_tag"])), None)
    if host is None:
        raise RuntimeError("Unsupported host platform for Python release tests")
    print("==> build current host binary for Python tests")
    run([
        sys.executable, str(BUILD_DIR / "build_core.py"),
        "--platform", host["wheel_tag"], "--version", version,
        "--out", str(REPO_ROOT / "python" / "minifyjs" / "bin"),
    ])


def step_python_tests() -> None:
    print("==> step 2: python tests")
    run([sys.executable, "-m", "pytest", "tests"], cwd=REPO_ROOT / "python")


def step_cross_compile(version: str) -> None:
    print("==> step 3: cross-compile")
    run([
        sys.executable, str(BUILD_DIR / "build_core.py"),
        "--all", "--out", str(DIST),
        "--version", version,
    ])


def step_verify_binaries() -> None:
    print("==> step 4: verify binaries")
    for platform_dir in sorted(DIST.iterdir()):
        if not platform_dir.is_dir():
            continue
        for binary in platform_dir.iterdir():
            # Only verify the host platform; the others can't run.
            if _is_host_platform(platform_dir.name):
                run([
                    sys.executable, str(BUILD_DIR / "scripts" / "verify_binary.py"),
                    str(binary),
                ])
            else:
                print(f"    {binary} (skipped: not host platform)")


def _is_host_platform(wheel_tag: str) -> bool:
    import platform
    system = platform.system().lower()
    machine = platform.machine().lower()
    if system == "linux" and machine in ("x86_64", "amd64"):
        return wheel_tag == "manylinux2014_x86_64"
    if system == "linux" and machine in ("aarch64", "arm64"):
        return wheel_tag == "manylinux2014_aarch64"
    if system == "darwin":
        if machine == "arm64":
            return wheel_tag == "macosx_11_0_arm64"
        return wheel_tag == "macosx_10_13_x86_64"
    if system == "windows":
        return wheel_tag == "win_amd64"
    return False


def step_build_wheels(version: str) -> None:
    print("==> step 5: build wheels")
    shutil.rmtree(WHEELHOUSE, ignore_errors=True)
    run([
        sys.executable, str(BUILD_DIR / "build_wheels.py"),
        "--all", "--version", version,
        "--out", str(WHEELHOUSE),
    ])


def step_verify_wheels() -> None:
    print("==> step 6: verify wheels (host only)")
    wheels = sorted(WHEELHOUSE.glob("*.whl"))
    host_wheels = [w for w in wheels if any(
        _is_host_platform(tag) and tag in w.name
        for tag in ("manylinux2014_x86_64", "manylinux2014_aarch64",
                    "macosx_10_13_x86_64", "macosx_11_0_arm64", "win_amd64")
    )]
    if not host_wheels:
        raise RuntimeError("No wheel available for this host; cannot verify release")
    for wheel in host_wheels:
        run([sys.executable, str(BUILD_DIR / "scripts" / "verify_wheel.py"), str(wheel)])


def step_checksums(version: str) -> None:
    print("==> step 7: checksums and manifest")
    run([
        sys.executable, str(BUILD_DIR / "generate_checksums.py"),
        "--dir", str(WHEELHOUSE),
        "--out", str(WHEELHOUSE / "checksums.txt"),
    ])
    run([
        sys.executable, str(BUILD_DIR / "generate_manifest.py"),
        "--dir", str(WHEELHOUSE),
        "--out", str(WHEELHOUSE / "manifest.json"),
        "--version", version,
    ])


def step_publish() -> None:
    print("==> step 8: publish to PyPI")
    wheels = sorted(WHEELHOUSE.glob("*.whl"))
    run([
        sys.executable, "-m", "twine", "upload",
        "--skip-existing",
        *(str(wheel) for wheel in wheels),
    ])


def main(argv: Optional[List[str]] = None) -> int:
    parser = argparse.ArgumentParser(description="Run the full MinifyJS release pipeline")
    parser.add_argument("--version", required=True, help="version string (e.g. 0.1.0)")
    parser.add_argument("--publish", action="store_true", help="upload to PyPI")
    parser.add_argument("--skip-tests", action="store_true", help="skip Go and Python tests")
    args = parser.parse_args(argv)

    if not args.skip_tests:
        step_go_tests()
        step_build_test_binary(args.version)
        step_python_tests()

    step_cross_compile(args.version)
    step_verify_binaries()
    step_build_wheels(args.version)
    step_verify_wheels()
    step_checksums(args.version)

    if args.publish:
        step_publish()

    print(f"\nartifacts staged in {WHEELHOUSE}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
