# Binary staging

Where the raw native binaries are staged before they are wrapped in
wheels. One subdirectory per platform, plus a `checksums/`
directory for the checksum files that end up on the GitHub release
page.

## Layout

```
packaging/binaries/
├── linux_x86_64/
│   └── minifyjs
├── linux_aarch64/
│   └── minifyjs
├── macos_x86_64/
│   └── minifyjs
├── macos_arm64/
│   └── minifyjs
├── windows_x86_64/
│   └── minifyjs.exe
└── checksums/
    ├── linux_x86_64.sha256
    ├── linux_aarch64.sha256
    └── ...
```

## Nothing here is checked in

Everything in this tree is a build artifact and is listed in
`.gitignore`:

```
packaging/binaries/*/minifyjs
packaging/binaries/*/minifyjs.exe
packaging/binaries/checksums/*.sha256
```

Only the `README.md` and the `checksums/.gitkeep` are tracked.

## Why binaries are staged separately from wheels

The wheel is the shipped artifact. The bare binary is a convenience
for users who want MinifyJS without Python at all.

Users who install the wheel get the binary automatically. Users who
want to download the binary directly (for a Docker image, a
Homebrew formula, an OS package, or a machine without Python) can
grab it from the GitHub release page, where the checksum files
attached here are the integrity guarantee.

## The checksum files

Each `.sha256` file has one line, the SHA-256 of the binary:

```
a1b2c3d4e5f6...  minifyjs
```

This is the standard `sha256sum` format, so it can be verified
with:

```console
sha256sum -c packaging/binaries/checksums/linux_x86_64.sha256
```

On macOS, `shasum -a 256 -c` works the same way.

## What gets uploaded to GitHub releases

The release workflow attaches:

- All six binaries (one per platform).
- The six `.sha256` files.
- The `manifest.json` from `build/generate_manifest.py`.

The wheel `.whl` files are uploaded to PyPI, not to the GitHub
release page. They can be downloaded from PyPI:

```console
pip download minifyjs==0.1.0 --no-deps -d /tmp/wheels
```