# Installation

## Requirements

- **Python 3.8 or newer.** Nothing else.
- **A supported platform:** Linux x86-64, Linux ARM64, macOS x86-64,
  macOS ARM64, or Windows x86-64.

## Install

```console
pip install minifyjs
```

That is the whole install. There is no Node.js, no npm, no Go
toolchain, no build step. The wheel contains a precompiled binary
for your platform.

Verify it works:

```console
minifyjs --version
# minifyjs 0.1.0 (esbuild 0.28.2)
```

## What just got installed

Three things:

1. The `minifyjs` command-line tool on your `PATH`.
2. The `minifyjs` Python package.
3. The native engine binary, inside the package at
   `minifyjs/bin/minifyjs` (or `minifyjs.exe` on Windows).

The Python package is a thin wrapper; it does not include a
JavaScript runtime. The native binary is the only thing that does
any work.

## Platform-specific notes

### Linux

The wheels are built on `manylinux2014`, which supports glibc 2.17
and newer. That covers RHEL 7, Ubuntu 14.04, Debian 8, and every
newer distribution.

If you are on a musl-based distribution (Alpine Linux), the
manylinux wheel will not work. Install from source instead:

```console
pip install --no-binary :all: minifyjs
```

That requires a Go toolchain. If you do not have Go, install it
first.

### macOS

Both Intel and Apple Silicon are supported. On macOS 11 or newer,
`pip install minifyjs` picks the arm64 wheel automatically.

The binary is not signed with an Apple Developer ID. If Gatekeeper
blocks it, run:

```console
xattr -d com.apple.quarantine $(python -c "from minifyjs._paths import binary_path; print(binary_path())")
```

### Windows

The wheel is for 64-bit Windows only. There is no 32-bit wheel.
The binary is not signed with an Authenticode certificate; Windows
SmartScreen may warn on first run.

## Installing from source

If you are on an unsupported platform or want to build the binary
yourself:

```console
git clone https://github.com/aldanedev-create/minifyjs
cd minifyjs

# Build the Go binary.
cd core
go build -o ../python/minifyjs/bin/minifyjs ./cmd/minifyjs

# Install the Python package pointing at the binary you just built.
cd ..
pip install -e python
```

The `-e` (editable) install means changes to the Python source take
effect without reinstalling. Changes to the Go source require
rerunning the `go build` command.

Requirements for a source install:

- **Go 1.22 or newer** ([download](https://go.dev/dl/))
- **Python 3.8 or newer**
- **Git** (to clone)

## Installing a specific version

```console
pip install minifyjs==0.1.0
```

## Upgrading

```console
pip install --upgrade minifyjs
```

The upgrade replaces both the Python package and the binary.
Existing caches are invalidated by the version bump.

## Uninstalling

```console
pip uninstall minifyjs
```

The cache directory (`~/.cache/minifyjs` on Linux and macOS,
`%LOCALAPPDATA%\minifyjs\Cache` on Windows) is left behind. Delete
it manually if you want a clean removal.

## Troubleshooting

### `pip install minifyjs` fails with "no matching distribution"

Your platform is not supported. Check the list above. If you are on
a supported platform and still see this, your `pip` is too old:

```console
pip install --upgrade pip
```

### `minifyjs: command not found`

The wheel installed, but the `minifyjs` script is not on your
`PATH`. This usually means `pip` installed to a user directory that
is not on `PATH`. On macOS and Linux:

```console
export PATH="$HOME/.local/bin:$PATH"
```

Add that line to your shell profile to make it permanent.

### `BinaryNotFoundError` when importing the Python package

The Python package is installed but the bundled binary is missing.
This can happen if you installed with `--no-binary` and the build
failed silently. Reinstall:

```console
pip install --force-reinstall minifyjs
```

### Antivirus blocks the binary on Windows

Some antivirus software flags unsigned binaries. The binary is
compiled from the source in this repository; you can verify this
yourself:

```console
pip download minifyjs --no-deps --no-binary :all: -d /tmp/src
tar -xzf /tmp/src/minifyjs-*.tar.gz -C /tmp/src
grep -r esbuild /tmp/src/minifyjs-*/core/go.mod
```

You will see the same `github.com/evanw/esbuild` dependency the
released wheel was built with.