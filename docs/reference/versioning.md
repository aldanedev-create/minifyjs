# Versioning

How MinifyJS assigns version numbers, and what a version change
promises.

## Semantic Versioning

MinifyJS follows [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** version for incompatible API changes.
- **MINOR** version for added functionality in a backwards-compatible way.
- **PATCH** version for backwards-compatible bug fixes.

A version looks like `MAJOR.MINOR.PATCH`, optionally with a
pre-release suffix: `0.2.0-rc1`.

## What "backwards compatible" means

The stable public surface is:

- **The CLI:** every flag, every exit code, the format of stdout
  and stderr.
- **The Python package:** every name in `minifyjs.__all__`.
- **The config file format:** every documented key.
- **The Go API:** every exported name in `core/api`.

A change is **backwards incompatible** if it:

- Removes or renames a public name.
- Changes the meaning of an existing flag.
- Changes the format of stdout or stderr in a way that a script
  parsing it would break.
- Changes an exit code.
- Changes the default value of a config key in a way that would
  change the output for existing users.

A change is **backwards compatible** if it:

- Adds a new flag, function, or config key.
- Changes an error message (as long as the error code is
  unchanged).
- Improves the minification (smaller output for the same input).
- Changes the output for a `--target` that was already in use (this
  is a "bug fix", not a breaking change, and it is documented in
  the changelog).

## Pre-1.0 versions

Before 1.0.0, the API is stable but the on-disk cache format and the
exact output bytes may change between minor versions. Specifically:

- The cache format may change in a **MINOR** release. Existing
  caches are invalidated by the version bump; users do not need to
  do anything.
- The exact output bytes may change in a **MINOR** release if the
  underlying esbuild version changes. Tests that compare against
  exact bytes should be written against the current release and
  regenerated when the release changes.

After 1.0.0, both of those become frozen. A cache format change
would be a **MAJOR** version. An output byte change for the same
input and options would be a **MAJOR** version.

## The version is stated in one place per language

- **Go:** `core/internal/version/version.go` (`var Version`)
- **Python:** `python/minifyjs/_version.py` (`__version__`)
- **`pyproject.toml`:** `version = "..."`

They must agree. `tools/release/check_release.py` verifies this
before a release.

## The esbuild version

MinifyJS pins a specific esbuild version in `core/go.mod`. The
version is exposed at runtime:

```console
$ minifyjs --version
minifyjs 0.2.0 (esbuild 0.28.2)
```

Upgrading esbuild is a **MINOR** version change for MinifyJS. It
can change the output for the same input and options.

Every esbuild upgrade is:

1. Announced in `CHANGELOG.md`.
2. Run against the full test suite before release.
3. Benchmarked, with results compared against the previous release.

An esbuild upgrade that makes output *larger* is investigated
before merging.

## Platform support

Adding support for a new platform is a **MINOR** version change.
Removing support is a **MAJOR** version change.

The platform support window is documented in
[../compatibility.md](../compatibility.md). Supported platforms do
not change within a **PATCH** release.

## Python version support

Adding support for a new Python version is a **MINOR** change.
Removing support for a Python version that was previously supported
is a **MAJOR** change.

The support matrix is:

| Python version | Status |
|---|---|
| 3.8 | Supported |
| 3.9 | Supported |
| 3.10 | Supported |
| 3.11 | Supported |
| 3.12 | Supported |
| 3.13+ | Not yet tested |

The minimum supported Python version is 3.8. Dropping 3.8 would be
a **MAJOR** change.

## Go version support

The minimum Go version is declared in `core/go.mod` (the `go 1.22`
line). Bumping it is a **MINOR** change for MinifyJS users, because
end users do not install Go; only contributors do.

## Deprecation

Before removing a public name, MinifyJS deprecates it:

1. The name continues to work, but its documentation is marked
   deprecated.
2. A warning is printed on stderr the first time it is used.
3. After one **MINOR** version, the name is removed in the next
   **MAJOR** version.

There is no automatic deprecation for config file keys. A key that
is removed is a **MAJOR** change, and the key's absence is what
users see; the change is documented in the changelog.

## Pre-release versions

Pre-release suffixes:

- `-alpha.N` — early, incomplete, expect changes.
- `-beta.N` — feature complete, expect bugs.
- `-rc.N` — release candidate, expect no changes.

`pip install minifyjs` does not install a pre-release. To install
one:

```console
pip install --pre minifyjs
```

Or explicitly:

```console
pip install minifyjs==0.2.0-rc1
```

## Yanking

If a release is broken and should not be installed:

1. **Yank on PyPI.** This makes `pip install minifyjs` (without a
   version pin) refuse to install it. Users who explicitly pin the
   yanked version can still install it.
2. **Edit the GitHub release** to add a warning.
3. **Cut a new release** with the fix.

The version number is never reused. A yanked `0.2.0` stays `0.2.0`;
the fix is `0.2.1`.

## See also

- [../../CHANGELOG.md](../../CHANGELOG.md) — the changelog
- [../development/release-process.md](../development/release-process.md) —
  how releases are cut
- [exit-codes.md](exit-codes.md) — the stable exit code contract