# Python: platform support

Which platforms have prebuilt wheels, and what to do on the ones
that do not.

## Supported platforms

| Platform | Wheel tag | Status |
|---|---|---|
| Linux x86-64 | `manylinux2014_x86_64` | ✅ |
| Linux ARM64 | `manylinux2014_aarch64` | ✅ |
| macOS x86-64 | `macosx_10_13_x86_64` | ✅ |
| macOS ARM64 | `macosx_11_0_arm64` | ✅ |
| Windows x86-64 | `win_amd64` | ✅ |
| Alpine Linux (musl) | — | ❌ No wheel |
| FreeBSD | — | ❌ No wheel |
| Windows ARM64 | — | ❌ No wheel |
| 32-bit anything | — | ❌ No wheel |

## What "no wheel" means

`pip install minifyjs` on an unsupported platform fails with:

```
ERROR: Could not find a version that satisfies the requirement minifyjs
```

The package is not on PyPI for your platform. You have two options:

1. Install from source (requires Go).
2. Build a wheel yourself and install it.

## Installing from source

Requires Go 1.22 or newer.

```console
git clone https://github.com/minifyjs/minifyjs
cd minifyjs
cd core
go build -o ../python/minifyjs/bin/minifyjs ./cmd/minifyjs
cd ..
pip install -e python
```

The `-e` install makes the Python package editable, so changes to
Python source take effect immediately. Changes to the Go source
require re-running `go build`.

On Windows, the binary is named `minifyjs.exe`:

```powershell
cd core
go build -o ..\python\minifyjs\bin\minifyjs.exe .\cmd\minifyjs
```

On Alpine Linux (musl), the Go build must be a static build to work
in the musl environment. Go produces static binaries by default with
`CGO_ENABLED=0`, which is what the build script uses.

## Building a wheel

To produce a wheel for your platform:

```console
cd minifyjs
python build/build_wheels.py --platform linux-x86_64 --version 0.1.0 --out wheelhouse
```

The `--platform` value comes from `build/platforms.json`. If your
platform is not listed, add it there first (see
[../development/adding-platform.md](../development/adding-platform.md)).

## Verifying the binary

After a source install:

```console
python -c "from minifyjs._binary import find_binary; print(find_binary())"
```

This prints the absolute path of the binary. If it raises
`BinaryNotFoundError`, the binary was not built or not found in the
expected location.

Then a functional check:

```console
python -c "from minifyjs import minify; print(minify('const x = 1;').code)"
# const x=1;
```

## Architecture detection

The Python package picks the binary by looking at:

- `platform.system()` — returns `Linux`, `Darwin`, or `Windows`
- `platform.machine()` — returns `x86_64`, `arm64`, `aarch64`, etc.

And mapping those to a filename: `minifyjs` on Linux and macOS,
`minifyjs.exe` on Windows.

There is only one binary per wheel, so the detection is mostly for
error messages. If you installed the wrong wheel, the error will
say so.

## Cross-platform notes

### Linux

The wheels are built on `manylinux2014`, which requires glibc 2.17
or newer. This covers:

- RHEL/CentOS 7+
- Ubuntu 14.04+
- Debian 8+
- Every modern distribution

It does **not** cover Alpine Linux, which uses musl libc. For
Alpine, build from source.

### macOS

Both Intel and Apple Silicon are supported. On macOS 11 or newer,
`pip install minifyjs` picks the arm64 wheel automatically.

The binary is not signed with an Apple Developer ID. Gatekeeper may
block it on first run:

```console
xattr -d com.apple.quarantine \
    $(python -c "from minifyjs._paths import binary_path; print(binary_path())")
```

After this, the binary runs normally. The quarantine attribute is
only set the first time a binary is downloaded.

### Windows

The wheel is for 64-bit Windows 10 or newer.

The binary is not signed with an Authenticode certificate. Windows
SmartScreen may warn the first time it runs. The warning can be
dismissed by clicking "More info" → "Run anyway".

Some antivirus software flags unsigned binaries. The binary is
compiled from the source in the repository; you can verify this
yourself:

```powershell
pip download minifyjs --no-deps --no-binary :all: -d C:\tmp\src
tar -xzf C:\tmp\src\minifyjs-*.tar.gz -C C:\tmp\src
Get-Content C:\tmp\src\minifyjs-*\core\go.mod
```

The `require` block will name `github.com/evanw/esbuild`, the same
dependency the released wheel was built with.

## Python version support

| Python version | Status |
|---|---|
| 3.8 | ✅ |
| 3.9 | ✅ |
| 3.10 | ✅ |
| 3.11 | ✅ |
| 3.12 | ✅ |
| 3.13+ | Not yet tested |

MinifyJS ships `py.typed`. Type checkers with strict mode work on
every supported version.

## Docker

The wheel installs cleanly inside `python:3.11-slim`. No C library
dependencies, no toolchain.

```dockerfile
FROM python:3.11-slim
RUN pip install minifyjs
RUN echo "const x = 1;" | minifyjs --compress
```

## CI matrix

The CI workflow at `.github/workflows/python-test.yml` runs the test
suite on every supported Python version on Linux. The cross-platform
tests run in the wheel-build workflows, one runner per platform.

If you find a platform that should be supported but is not, open an
issue. Adding a platform is a checklist of about 6 steps; see
[../development/adding-platform.md](../development/adding-platform.md).

## See also

- [../compatibility.md](../compatibility.md) — the OS and Python
  version matrix, in one place
- [errors.md](errors.md) — `BinaryNotFoundError` and how to
  diagnose it