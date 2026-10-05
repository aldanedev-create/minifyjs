# Build and release automation

This directory contains everything that turns the Go source in
`core/` into a signed, tested, PyPI-ready set of platform wheels.

Nothing here runs during normal development. `go build ./...` in
`core/` is enough to get a working binary. This directory exists for
the maintainer cutting a release.

## The pipeline

```
core/  ──▶  go build  ──▶  minifyjs (one binary)
                              │
                              ├── linux x86_64   ─┐
                              ├── linux aarch64   │
                              ├── macos x86_64    ├──▶ Python wheel per platform
                              ├── macos arm64     │
                              └── windows x86_64 ─┘
                                        │
                                        ▼
                                    PyPI (minifyjs)
```

End users get:

```console
pip install minifyjs
```

No Go toolchain. No npm. No Node.js. One binary per platform, chosen
by the wheel tag.

## Entry points

| Script | What it does |
|---|---|
| `build_core.py` | Compiles the Go binary for one platform. Called once per target. |
| `build_wheels.py` | Assembles a platform-specific wheel with the binary inside. |
| `build_release.py` | Full pipeline: build, test, wheel, checksum, publish. |
| `generate_checksums.py` | SHA-256 checksums of every artifact, for the release page. |
| `generate_manifest.py` | A JSON manifest describing every file in a release. |
| `cross_compile.sh` | One-shot script to build all binaries locally. |
| `scripts/verify_binary.py` | Runs a smoke test on a built binary. |

## Platforms

`platforms.json` is the source of truth for what we ship. Every
entry maps a Python wheel tag to a Go `GOOS`/`GOARCH` pair.

`targets.json` is the source of truth for what versions we build
against (Go, Python, esbuild).

Adding a new platform means:

1. Add an entry to `platforms.json`.
2. Add a cross-compilation step to `cross_compile.sh`.
3. Add the platform to `.github/workflows/build-*.yml`.
4. Update `python/minifyjs/_platform.py` to recognize the new tag.
5. Update `docs/reference/versioning.md` with the support window.

## Running the pipeline locally

```console
cd build
python build_core.py --platform linux-x86_64 --out /tmp/minifyjs
python scripts/verify_binary.py /tmp/minifyjs
```

The full release pipeline runs only in CI. It is triggered by pushing
a tag of the form `v1.2.3`.

## What is checked in vs generated

**Checked in:** every script and JSON file in this directory.

**Generated, never checked in:**

- `build/dist/` — intermediate binaries
- `build/wheelhouse/` — built wheels
- `build/release-manifest.json` — release output

These are in `.gitignore`.