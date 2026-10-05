# Compatibility

What MinifyJS runs on, what it produces, and what it is tested
against.

## Python

MinifyJS supports Python 3.8 and newer. Wheels are published for
every minor version from 3.8 through 3.12 on the day of release.

| Python version | Status |
|---|---|
| 3.8 | Supported |
| 3.9 | Supported |
| 3.10 | Supported |
| 3.11 | Supported |
| 3.12 | Supported |
| 3.13 | Metadata and compile checked; runtime depends on the host |
| 3.14+ | Metadata and compile checked; runtime depends on the host |

## Platforms

Six platforms are officially supported. Each has its own wheel on
PyPI.

| Platform | Wheel tag | Status |
|---|---|---|
| Linux x86-64 | `manylinux2014_x86_64` | Supported |
| Linux ARM64 | `manylinux2014_aarch64` | Supported |
| macOS x86-64 | `macosx_10_13_x86_64` | Supported |
| macOS ARM64 (Apple Silicon) | `macosx_11_0_arm64` | Supported |
| Windows x86-64 | `win_amd64` | Supported |
| Anything else | — | Not supported |

"Not supported" means "no wheel". If `pip install minifyjs` fails on
your platform, it is because no wheel exists for it. See
[development/adding-platform.md](development/adding-platform.md) if
you want to contribute one.

## Operating systems

The wheels cover:

| OS | Minimum version |
|---|---|
| Linux | glibc 2.17 (RHEL 7, Ubuntu 14.04, Debian 8) |
| macOS | 10.13 (High Sierra) |
| Windows | 10 |

The Linux wheels use the `manylinux2014` baseline. Older glibc
versions are not supported.

## JavaScript targets

The output of MinifyJS can target any ECMAScript version esbuild
supports, from `es5` through `esnext`. See
[supported-syntax.md](supported-syntax.md).

The *input* is always ES2024 plus JSX plus TypeScript type syntax.
There is no mode that accepts ES3 input.

## esbuild versions

MinifyJS pins a specific esbuild version in `core/go.mod`. The
version is exposed at runtime:

```console
$ minifyjs --version
minifyjs 0.1.0 (esbuild 0.28.2)
```

Upgrading esbuild is a minor-version change for MinifyJS. Every
upgrade is announced in `CHANGELOG.md` and the test suite is run
against the new version before release.

## Reproducibility

Given the same input, the same options, and the same MinifyJS
version, MinifyJS produces byte-identical output. This is enforced
by tests in `core/integration/minify/idempotence_test.go`.

The output is not guaranteed to be identical *between* MinifyJS
versions. A new esbuild version may produce different (usually
smaller) output for the same input. This is intentional; the whole
point of upgrading is to get better output.

## Things MinifyJS is compatible with

- **Django staticfiles.** See [integrations/django.md](integrations/django.md).
- **Flask static folder.** See [integrations/flask.md](integrations/flask.md).
- **FastAPI StaticFiles.** See [integrations/fastapi.md](integrations/fastapi.md).
- **GitHub Actions, GitLab CI, CircleCI.** Shell out to `minifyjs`.
- **Docker.** The wheel installs cleanly inside `python:3.11-slim`.
- **Homebrew.** A formula is published under `packaging/homebrew/`.

## Things MinifyJS is not compatible with

- **Node.js build pipelines that expect a plugin API.** Use esbuild
  directly.
- **Bundlers other than esbuild's.** MinifyJS does not produce
  output that webpack, Rollup, or Vite can consume as a plugin.
- **The `terser` API.** MinifyJS's API is different. There is no
  compatibility shim.