#!/usr/bin/env bash
# Runs once after the container is created, as the vscode user.
#
# Installs the project's own dependencies (Go modules, Python
# dev dependencies, Node modules for the benchmark suite) and
# installs the pre-commit hooks.
#
# This script is idempotent: running it twice is safe.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

echo "==> configuring git"
git config --global --add safe.directory "$REPO_ROOT" || true
git config --global core.autocrlf input || true
git config --global pull.rebase true || true

echo "==> syncing Go workspace"
go work sync

echo "==> downloading Go modules"
(cd core && go mod download)
(cd tools && go mod download)

echo "==> installing Python package (editable) with dev extras"
python -m pip install --user -e "python[dev]"

echo "==> building the CLI"
(cd core && go build -o bin/minifyjs ./cmd/minifyjs)

echo "==> copying the CLI into the Python package"
cp core/bin/minifyjs python/minifyjs/bin/minifyjs 2>/dev/null || \
    cp core/bin/minifyjs.exe python/minifyjs/bin/minifyjs.exe 2>/dev/null || \
    echo "    (binary copy skipped; the platform name did not match)"

echo "==> installing pre-commit hooks"
pre-commit install || true

echo "==> installing Node dependencies for benchmarks"
if [ -f bench/package.json ]; then
    (cd bench && npm install)
fi

echo
echo "==> development environment ready"
echo
echo "Common commands:"
echo "  make build      # build the CLI"
echo "  make test       # run the full test suite"
echo "  make lint       # run linters"
echo "  make bench      # run benchmarks"
echo
echo "Build the CLI manually with:"
echo "  cd core && go build -o bin/minifyjs ./cmd/minifyjs"