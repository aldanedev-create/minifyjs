# MinifyJS developer tools

Small internal utilities that the maintainers run during
development, testing, and releasing. None of these ship to end users
and none of them are part of the `minifyjs` binary.

## The tools

| Tool | Language | What it does |
|---|---|---|
| `fixture_runner` | Go | Walks `core/testdata/fixtures/`, runs each input through the engine, and reports pass/fail per fixture |
| `syntax_matrix` | Go | Walks `core/testdata/fixtures/`, reads each `README.md`, and produces a markdown table of which ES versions support which features |
| `benchmark_report` | Go | Reads `bench/results/history/*.json` and prints the delta between the two most recent runs, flagging regressions |
| `release/check_release.py` | Python | Pre-release checklist: version consistency, changelog entry, tag format, no uncommitted files |
| `release/verify_artifacts.py` | Python | Verifies that a staged release contains every platform wheel declared in `build/platforms.json` |
| `release/verify_checksums.py` | Python | Verifies a `checksums.txt` against the files it names |
| `docs/generate_cli_reference.py` | Python | Parses `core/cmd/minifyjs/cli.go` and produces `docs/cli/reference.md` |
| `docs/generate_syntax_matrix.py` | Python | Parses `syntax_matrix` output and writes `docs/supported-syntax.md` |

## Building and running

The Go tools are built via the workspace:

```console
# From the repository root.
go build ./tools/...
go run ./tools/fixture_runner
go run ./tools/syntax_matrix
go run ./tools/benchmark_report
```

The Python tools are run directly:

```console
python tools/release/check_release.py --version 1.2.3
python tools/docs/generate_cli_reference.py --output docs/cli/reference.md
```

## Why a separate Go module

If `tools/` shared `core/`'s `go.mod`, every dependency a tool pulled
in would also be a dependency of the released binary. A dev tool that
needs a YAML parser, for instance, would add that dependency to the
wheel. Separate modules prevent that class of mistake by construction.

## Adding a new tool

1. Create `tools/<name>/`.
2. Add a `main.go` (or `__main__.py`) that is runnable from the
   repository root.
3. Add it to the table above.
4. If it is a Go tool, no further wiring is needed — `go build ./...`
   picks it up automatically.

Tools are meant to be small. A tool that grows past a few hundred
lines should become a package under `core/internal/` (if it is used
by the engine) or a separate repository (if it is not).