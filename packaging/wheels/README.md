# Wheel staging

This directory is where the release pipeline stages the wheels it
builds, one subdirectory per platform. The subdirectories exist so
that a maintainer reviewing a release can see exactly what was built
for each platform before it is uploaded to PyPI.

## What lives here

Each subdirectory corresponds to a platform tag from
`build/platforms.json`:

| Directory | Wheel platform tag |
|---|---|
| `linux_x86_64/` | `manylinux2014_x86_64` |
| `linux_aarch64/` | `manylinux2014_aarch64` |
| `macos_x86_64/` | `macosx_10_13_x86_64` |
| `macos_arm64/` | `macosx_11_0_arm64` |
| `windows_x86_64/` | `win_amd64` |

Each subdirectory holds:

- The `.whl` file produced by `build/build_wheels.py`.
- A `checksums.txt` with the SHA-256 of the wheel.

## Nothing here is checked in

The `.gitkeep` files exist only so git tracks the empty directories.
Every other file in this tree is a build artifact and is listed in
`.gitignore`:

```
packaging/wheels/*/*.whl
packaging/wheels/*/checksums.txt
```

## When this directory is populated

Only during a release. The pipeline:

1. Cross-compiles the binary for each platform.
2. Builds a wheel with the binary inside.
3. Copies the wheel into the matching subdirectory.
4. Generates `checksums.txt`.

The staging directory is then reviewed by a maintainer before
upload. Nothing here is uploaded to PyPI directly; the pipeline
uploads from `build/wheelhouse/`.

## Why this directory exists

Two reasons:

1. **Reviewability.** A maintainer can `ls packaging/wheels/` and
   see, at a glance, which platforms were built. A missing platform
   is obvious; a wrong-sized wheel is obvious.

2. **CI artifact separation.** GitHub Actions uploads build
   artifacts into per-job directories. This tree gives each job a
   well-known place to drop its output, so the wheels step can
   collect them without a search step.

## Verifying a staged wheel

```console
cd packaging/wheels/linux_x86_64
sha256sum -c checksums.txt
```

The same verification is available for every platform.