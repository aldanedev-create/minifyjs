#!/bin/sh
# Cross-compile the MinifyJS binary for every platform in
# platforms.json. Output goes to build/dist/<wheel_tag>/minifyjs[.exe].
#
# This script is idempotent: re-running it replaces existing binaries.
#
# On CI, each platform is built on its own runner (see
# .github/workflows/build-*.yml) so that macOS codesigning and Windows
# version resources can be applied. This script is for local
# development and for reproducing a build on a single machine.

set -eu

cd "$(dirname "$0")"

# Locate the Python interpreter.
PY="${PYTHON:-python3}"
command -v "$PY" >/dev/null 2>&1 || PY=python
command -v "$PY" >/dev/null 2>&1 || { echo "no python interpreter found" >&2; exit 1; }

# Ask platforms.json for the list of platforms.
PLATFORMS=$("$PY" -c '
import json, sys
data = json.load(open("platforms.json"))
for p in data["platforms"]:
    print(p["wheel_tag"], p["goos"], p["goarch"], p["binary_name"])
')

DIST="dist"
mkdir -p "$DIST"

echo "==> building binaries into $DIST"
echo "$PLATFORMS" | while read -r tag goos goarch binname; do
    outdir="$DIST/$tag"
    outfile="$outdir/$binname"
    mkdir -p "$outdir"

    echo "    $tag -> $goos/$goarch"
    GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 \
        go build -trimpath \
            -ldflags "-s -w -X github.com/minifyjs/minifyjs/core/internal/version.Version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" \
            -o "$outfile" \
            ../core/cmd/minifyjs
done

echo "==> done"
ls -la "$DIST"