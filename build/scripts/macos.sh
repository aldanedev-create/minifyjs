#!/bin/sh
# macOS-specific build wrapper. Called by CI on macos runners.
#
# Builds both x86-64 and ARM64 binaries. The x86-64 binary can be
# built on either architecture; the ARM64 binary requires an ARM
# runner.

set -eu

cd "$(dirname "$0")/.."

: "${VERSION:=$(git describe --tags --always --dirty 2>/dev/null || echo dev)}"
echo "building MinifyJS for macOS (version $VERSION)"

ARCH="$(uname -m)"
if [ "$ARCH" = "arm64" ]; then
    python build_core.py --platform macos-arm64 --out dist --version "$VERSION"
    python scripts/verify_binary.py dist/macosx_11_0_arm64/minifyjs
else
    python build_core.py --platform macos-x86_64 --out dist --version "$VERSION"
    python scripts/verify_binary.py dist/macosx_10_13_x86_64/minifyjs
fi