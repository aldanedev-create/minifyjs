# Release process

How to cut a MinifyJS release.

## Who can cut a release

Maintainers. If you are not a maintainer and want a release, open an
issue. See [../../GOVERNANCE.md](../../GOVERNANCE.md).

## Prerequisites

- The `main` branch is green on CI.
- Every change you want in the release has been merged.
- You have push access to `main` and permission to publish to PyPI.
- You have the PyPI API token in your environment or in a
  `.pypirc`.

## The one-command release

From a clean `main`:

```console
python build/build_release.py --version 0.1.0
```

That runs the whole pipeline. See the sections below for what each
step does and how to verify it.

## The pipeline

```text
1. Go tests                    go test ./... in core/
2. Python tests                pytest tests in python/
3. Cross-compile               build_core.py --all
4. Verify native host binary   verify_binary.py
5. Build wheels                build_wheels.py --all
6. Verify host wheel           fresh venv, CLI and Python bundle smoke tests
7. Generate checksums/manifest
8. Optional publication        twine upload (only with --publish)
```

The wheel workflow also builds and installs each platform wheel on a matching
Linux, macOS, or Windows runner. musllinux wheels are installed and exercised
in Alpine containers for x86-64 and ARM64. Host selection reads pip compatibility
tags so musl environments verify musllinux wheels and glibc environments verify
manylinux wheels. Wait for all platform jobs to pass before
publishing. Cross-compiling alone does not verify that a foreign-platform binary
runs. This CI tests modern operating systems; it does not establish compatibility
with every older OS implied by a wheel tag.

The local verification step fails when no supported host wheel is available.
It installs without dependencies or a source checkout, verifies both CLI entry
points, and checks bundle splitting, source maps, metadata, tree shaking,
working-directory handling, and explicit dependency externalization.

## Step by step

### 1. Update the version

The version lives in four places:

- `core/internal/version/version.go` — `var Version`
- `python/minifyjs/_version.py` — `__version__`
- `python/pyproject.toml` — `version = "..."`
- `bench/package.json` — `"version"` (cosmetic; does not affect
  output)

Change all four. The release check will verify that they agree.

### 2. Update `CHANGELOG.md`

Add an entry under a new `## [x.y.z]` heading:

```markdown
## [0.2.0] - 2026-01-20

### Added

- New `--sourcemap=both` mode (#42).

### Fixed

- Preserve shebang lines through minification (#55).

## [0.1.0] - 2026-01-01
```

The format is [Keep a Changelog](https://keepachangelog.com/).

### 3. Run the release check

```console
python tools/release/check_release.py --version 0.2.0
```

Verifies:

- The version is `MAJOR.MINOR.PATCH`.
- Every file that declares a version agrees.
- `CHANGELOG.md` has a section for the version.
- The working tree is clean.
- The current branch is `main`.

Fix any problem it reports before continuing.

For a local rehearsal with uncommitted changes, use
`--skip-git`. A real release must still be made from a clean `main`
worktree.

### 4. Commit the version bump and changelog

```console
git add -A
git commit -m "release: v0.2.0"
git push origin main
```

Wait for CI to pass. Do not proceed until it is green.

### 5. Tag the release

```console
git tag -a v0.2.0 -m "v0.2.0"
git push origin v0.2.0
```

The tag triggers the release workflow at
`.github/workflows/release.yml`. That workflow:

1. Runs the release check.
2. Builds wheels for every platform.
3. Verifies each wheel.
4. Generates checksums and a manifest.
5. Publishes to PyPI.
6. Creates a GitHub release with the artifacts attached.

The workflow takes about 15 minutes. Watch it at
`https://github.com/minifyjs/minifyjs/actions`.

### 6. Verify the release

Wait for the workflow to finish. Then:

```console
# In a fresh venv.
python -m venv /tmp/release-check
. /tmp/release-check/bin/activate
pip install minifyjs==0.2.0
minifyjs --version
# minifyjs 0.2.0 (esbuild 0.28.2)

echo "const x = 1 + 2 + 3;" | minifyjs --compress
# const x=6;
```

If that works, the release is live.

## What can go wrong

### CI fails on the tag

The tag is pushed; the release workflow fails.

Fix the problem on `main`, then delete the tag and re-tag:

```console
git tag -d v0.2.0
git push origin :refs/tags/v0.2.0
# Fix main.
git tag -a v0.2.0 -m "v0.2.0"
git push origin v0.2.0
```

Do not push a second tag with a different version. `0.2.0` is
either the release or it is not.

### PyPI rejects the upload

Most common causes:

- **The version already exists.** PyPI does not allow re-uploading
  a version. Bump to a new patch (`0.2.1`) and re-tag.
- **The wheel is malformed.** Check `python -m twine check
  build/wheelhouse/*.whl` locally.
- **The API token is wrong.** Regenerate it on PyPI.

### The wheel installs but `minifyjs --version` fails

The binary is missing or not executable. See
[adding-platform.md](adding-platform.md) for the diagnosis.

### The GitHub release is empty

The workflow created a release but did not upload artifacts. Check
the workflow logs for `actions/upload-artifact` errors. The most
common cause is a mismatch between the artifact path in the
workflow and the actual output directory.

## Post-release

After the workflow is green:

1. **Update the docs** if the release added a new flag or API. The
   `docs/cli/reference.md` file is generated; regenerate it:

   ```console
   python tools/docs/generate_cli_reference.py --output docs/cli/reference.md
   git add docs/cli/reference.md
   git commit -m "docs: regenerate CLI reference for v0.2.0"
   git push origin main
   ```

2. **Update the benchmark results** if the release is a minor
   version bump:

   ```console
   cd bench
   python run_all.py
   git add results/latest.md
   git commit -m "bench: v0.2.0 results"
   git push origin main
   ```

3. **Announce the release** in a GitHub Discussion. The release
   notes are auto-generated from the changelog; the discussion is
   for context and questions.

## Hotfix releases

For a critical bug that cannot wait for the next minor release:

1. Branch from the release tag:

   ```console
   git checkout -b hotfix/0.2.1 v0.2.0
   ```

2. Fix the bug. Add a test. Update `CHANGELOG.md`.

3. Bump the version to `0.2.1` in the four places.

4. Merge to `main` first, then tag from the branch:

   ```console
   git checkout main
   git merge --no-ff hotfix/0.2.1
   git push origin main

   git tag -a v0.2.1 -m "v0.2.1"
   git push origin v0.2.1
   ```

The release workflow runs from the tag.

## Yanking a release

If a release is broken and users should not install it:

1. **Yank on PyPI.** This does not delete the files; it makes
   `pip install minifyjs` without a version pin refuse to install
   them.
2. **Edit the GitHub release** to add a warning at the top.
3. **Cut a new release** with the fix.
4. **Do not delete the tag.** A deleted tag means the version
   number is available again, which is confusing.

## Release cadence

There is no fixed cadence. Releases happen when:

- A meaningful feature is ready.
- An esbuild upgrade brings smaller output.
- A bug fix cannot wait.

The project does not do scheduled releases. It does not do "release
trains". A release is cut when it is ready.

## Version numbers

MinifyJS follows [Semantic Versioning](https://semver.org/):

- **MAJOR** for incompatible API changes.
- **MINOR** for new functionality in a backwards-compatible way.
- **PATCH** for bug fixes.

Before 1.0.0, the API is stable but the on-disk cache format and
the exact output bytes may change between minor versions. See
[`../reference/versioning.md`](../reference/versioning.md).

## See also

- [`build/README.md`](https://github.com/aldanedev-create/minifyjs/tree/main/build) —
  the release pipeline scripts
- [adding-platform.md](adding-platform.md) — how to add a platform
  before the next release
- [../../CHANGELOG.md](../../CHANGELOG.md) — the changelog format