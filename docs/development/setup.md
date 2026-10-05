# Development setup

How to get a working MinifyJS development environment.

## Requirements

- **Go 1.22 or newer** — [download](https://go.dev/dl/)
- **Python 3.8 or newer** — [download](https://www.python.org/downloads/)
- **Git** — to clone the repository
- **(Optional) Node.js 18+** — only needed to run the benchmark suite
  against the JavaScript-based tools
- **(Optional) `golangci-lint`** — for the linting step
- **(Optional) `pre-commit`** — for the pre-commit hooks

You do not need Node.js to build or test MinifyJS. Node is only for
the benchmark comparisons.

## Clone and build

```console
git clone https://github.com/minifyjs/minifyjs
cd minifyjs

# Sync the Go workspace.
go work sync

# Build the CLI.
cd core
go build -o bin/minifyjs ./cmd/minifyjs
cd ..

# Install the Python package in editable mode.
cd python
python -m pip install -e ".[dev]"
cd ..
```

Verify:

```console
core/bin/minifyjs --version
# minifyjs 0.1.0 (esbuild 0.28.2)

python -c "from minifyjs import minify; print(minify('const x = 1;').code)"
# const x=1;
```

On Windows, the binary is `core\bin\minifyjs.exe`.

## Editor setup

The repository has an `.editorconfig` that keeps indentation,
line endings, and trailing whitespace consistent. Most editors
respect it automatically. If yours does not, install the
EditorConfig plugin.

For VS Code, the recommended extensions are listed in
`.vscode/extensions.json` (not checked in; each contributor chooses
their own).

## The Go workspace

The repository has two Go modules:

- `core/` — the engine, CLI, and API. This is what ships.
- `tools/` — developer utilities. This is not shipped.

`go.work` ties them into one workspace, so a single `go build ./...`
or `go test ./...` at the repository root builds and tests both.

The split exists so that a dependency a tool pulls in cannot leak
into the released binary's dependency graph.

## Common commands

From the repository root:

| Command | What it does |
|---|---|
| `go build ./...` | Build every Go package in both modules |
| `go test ./...` | Run every Go test |
| `go vet ./...` | Run vet on both modules |
| `make build` | Build the CLI into `core/bin/minifyjs` |
| `make test` | Run the full test suite (Go + Python) |
| `make lint` | Run linters for both languages |
| `make fmt` | Format everything |
| `make bench` | Run the benchmark suite |
| `make docs` | Regenerate `docs/cli/reference.md` and `docs/supported-syntax.md` |
| `make clean` | Remove build artifacts |

## Running the Python tests

```console
cd python
python -m pytest tests -v
```

Tests marked `@pytest.mark.binary` require the compiled binary at
`python/minifyjs/bin/minifyjs`. Build it first:

```console
cd core
go build -o ../python/minifyjs/bin/minifyjs ./cmd/minifyjs
```

On Windows:

```powershell
cd core
go build -o ..\python\minifyjs\bin\minifyjs.exe .\cmd\minifyjs
```

If the binary is missing, those tests skip rather than fail. That
is intentional: a fresh clone can run the test suite without
building anything.

## Running the integration tests

```console
cd core
go test ./integration/... -v
```

The integration tests compile the CLI once and run it as a
subprocess. The first run takes 30–60 seconds because of the
compile step; subsequent runs are fast.

## Pre-commit hooks

Install them once per clone:

```console
pip install pre-commit
pre-commit install
```

The hooks run formatting, linting, and a few file checks on
staged changes. They are the same checks CI runs, so a clean
pre-commit run is a good sign that CI will pass.

To run them manually on all files:

```console
pre-commit run --all-files
```

## Environment variables

A few development-time environment variables:

| Variable | Effect |
|---|---|
| `MINIFYJS_VERBOSE=1` | The CLI prints extra detail on stderr |
| `UPDATE_GOLDEN=1` | Tests that compare against golden files overwrite them instead |
| `SKIP_BINARY_TESTS=1` | Python tests that require the binary skip unconditionally |

None of these affect production behavior.

## Troubleshooting

### `go work sync` fails

Check that `go.work` exists at the repository root and lists both
`./core` and `./tools`. If one of the modules is missing a `go.mod`,
`go work sync` will fail with a clear error.

### `go build ./...` fails in `core/`

The most common cause is an esbuild version mismatch. Run:

```console
cd core
go mod tidy
go build ./...
```

If the error mentions a specific missing package, check
`core/go.mod`.

### Python tests fail with `BinaryNotFoundError`

The tests need the compiled binary. Build it:

```console
cd core
go build -o ../python/minifyjs/bin/minifyjs ./cmd/minifyjs
```

Or skip those tests:

```console
cd python
SKIP_BINARY_TESTS=1 python -m pytest tests
```

### Windows: `go build` creates `minifyjs.exe` but tests look for `minifyjs`

The Python package's platform detection handles this
automatically via `_platform.binary_filename()`. If you see a
`BinaryNotFoundError` on Windows, verify that
`python/minifyjs/bin/minifyjs.exe` exists:

```powershell
Test-Path python\minifyjs\bin\minifyjs.exe
```

If it does not, the Go build wrote to a different location.

### Benchmark suite skips every tool

The benchmark suite skips tools that are not on `PATH`. To install
them locally:

```console
cd bench
npm install
```

This installs esbuild, terser, and uglify-js into
`bench/node_modules/.bin/`. The `bench/run_all.py` script adds that
directory to `PATH` before running.

## Where to go next

- [architecture.md](architecture.md) — the code layout in detail
- [testing.md](testing.md) — how to write and run tests
- [benchmarking.md](benchmarking.md) — how to run benchmarks
- [../../CONTRIBUTING.md](../../CONTRIBUTING.md) — the contribution
  workflow