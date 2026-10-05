# Contributing to MinifyJS

Thanks for your interest in MinifyJS. This document explains how to
set up a development environment, run the test suite, and get a
change merged.

## Code of conduct

By participating you agree to abide by the
[Code of Conduct](CODE_OF_CONDUCT.md). Please report unacceptable
behavior to the maintainers listed in that document.

## What we are looking for

MinifyJS is a small, focused project. It is a wrapper around
[esbuild](https://esbuild.github.io/)'s Go API, not a from-scratch
JavaScript toolchain. Before writing code, please open an issue to
discuss what you want to change. Features that expand the scope of
the project (a new language, a new plugin system, a new bundler) are
unlikely to be accepted; bug fixes, documentation improvements, and
small ergonomic improvements are very welcome.

## Development setup

You need:

- Go 1.22 or newer ([download](https://go.dev/dl/))
- Python 3.8 or newer ([download](https://www.python.org/downloads/))
- (Optional) `golangci-lint` and `pre-commit` if you want the same
  checks CI runs

Clone and build:

```console
git clone https://github.com/minifyjs/minifyjs
cd minifyjs
go work sync
cd core
go build ./...
go test ./...
```

Build the CLI:

```console
cd core
go build -o bin/minifyjs ./cmd/minifyjs
```

Install the Python package in editable mode:

```console
cd python
python -m pip install -e ".[dev]"
python -m pytest tests
```

## Repository layout

See [`docs/architecture.md`](docs/architecture.md) for a full
description. The short version:

- `core/` — the Go module. Everything that ships in the native
  binary lives here.
- `core/internal/engine/` — the only package that imports esbuild.
- `core/api/` — the public Go surface.
- `core/cmd/minifyjs/` — the CLI.
- `python/` — the Python package and its tests.
- `integrations/` — framework adapters (Telocel, generic examples).
- `docs/` — documentation.
- `bench/` — benchmark scripts.
- `tools/` — developer utilities.

## Making a change

1. **Open an issue first** if the change is more than a trivial fix.
2. **Fork and branch** from `main`. Use a descriptive branch name
   (`fix/stdin-crlf`, `docs/sourcemap-example`).
3. **Write tests.** Every behavioral change needs a test in the
   package that implements it. Fixes need a regression test that
   fails before the fix and passes after.
4. **Run the full suite** before pushing:
   ```console
   make test
   ```
5. **Keep the diff small.** One logical change per pull request.
6. **Update `CHANGELOG.md`** under `[Unreleased]`. If you are adding
   a flag or an API, also update the relevant file under `docs/`.

## Commit messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/)
style:

```
<type>(<scope>): <short description>

<body>

<footer>
```

Common types: `feat`, `fix`, `docs`, `test`, `chore`, `refactor`,
`perf`, `build`, `ci`.

Example:

```
fix(cli): preserve CRLF when writing output on Windows

The CLI was writing through os.Stdout in text mode, which
translated \n to \r\n. Minified output must be written
byte-for-byte, so the CLI now writes through the atomic-write
path for files and a byte-preserving writer for stdout.

Fixes #42
```

## Pull request review

A pull request is ready to merge when:

- CI is green on all platforms.
- At least one maintainer has approved.
- The `CHANGELOG.md` entry is present.
- New flags or API surface are documented.

We aim to review pull requests within one week. If a week passes
without a response, feel free to leave a comment on the pull request
poking the maintainers.

## Reporting bugs

Open an issue using the **Bug report** template. Include:

- MinifyJS version (`minifyjs --version`).
- Operating system and architecture.
- The smallest input that reproduces the problem.
- The actual output and the expected output.

If the bug is a crash or a hang, include a stack trace if one is
available (`MINIFYJS_VERBOSE=1 minifyjs ...` prints extra
information).

## Security issues

Do **not** open a public issue for security problems. See
[SECURITY.md](SECURITY.md).

## License

By contributing you agree that your contributions are licensed under
the MIT License (see [LICENSE](LICENSE)).