#!/bin/sh
# Linux-specific build wrapper. Called by CI on ubuntu runners.
# Ensures Go and Python are on PATH and delegates to build_core.py.
#
# Local users normally do not need this; `python build_core.py` from
# the build/ directory works on any platform.

set -eu

cd "$(dirname "$0")/.."

: "${VERSION:=$(git describe --tags --always --dirty 2>/dev/null || echo dev)}"
echo "building MinifyJS for linux (version $VERSION)"

python build_core.py --platform linux-x86_64 --out dist --version "$VERSION"
python build_core.py --platform linux-aarch64 --out dist --version "$VERSION"

python scripts/verify_binary.py dist/manylinux2014_x86_64/minifyjs