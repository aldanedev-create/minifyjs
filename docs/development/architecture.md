# Development: architecture

A contributor's tour of the codebase. For the conceptual overview,
see [../architecture.md](../architecture.md). This page is about
the files.

## Repository layout

```
minifyjs/
├── core/                    the released Go module
├── python/                  the Python package
├── integrations/            framework adapters (examples)
├── build/                   release automation
├── bench/                   benchmark suite
├── tools/                   developer utilities (separate Go module)
├── docs/                    this documentation
├── examples/                runnable examples
├── test/                    cross-cutting tests
├── config/                  config file presets
├── .github/                 CI, issue templates
├── packaging/               wheel, binary, homebrew staging
└── .devcontainer/           development container
```

## The core module

```
core/
├── go.mod                   module github.com/minifyjs/minifyjs/core
├── cmd/minifyjs/            the CLI's main package
├── api/                     the stable public Go surface
├── internal/
│   ├── engine/              the esbuild wrapper
│   ├── config/              config file loading
│   ├── cache/               optional on-disk result cache
│   ├── diagnostics/         structured errors, formatters
│   ├── source/              positions, spans, line maps
│   ├── io/                  atomic writes, stdin/stdout
│   └── version/             version string
├── testdata/fixtures/       fixture inputs and expected outputs
└── integration/             CLI subprocess tests
```

### `cmd/minifyjs/`

The CLI. It parses arguments, resolves config, reads input, calls
`api`, writes output, and translates diagnostics into exit codes.
It knows about flags and nothing else.

```
main.go       entry point; calls run()
cli.go        the run function; orchestrates everything
flags.go      flag parsing
parse.go      the argument state machine
input.go      reads stdin or a file
output.go     writes stdout or a file
stdin.go      stdin wrapper
stdout.go     stdout wrapper
watch.go      the polling watch loop
errors.go     CLI-specific error types
```

### `api/`

The public Go surface. It re-exports what external callers need and
maps its own types to the engine's.

```
minify.go       Minify() -> engine.Transform
optimize.go     Optimize() -> Minify() with compress+mangle
bundle.go       Bundle() -> engine.Build
options.go      the Options type
result.go       the Result type
errors.go       sentinel errors
diagnostics.go  re-exports the diagnostics types
```

`api` is the only package external callers import.

### `internal/engine/`

The esbuild wrapper. See [../engine/overview.md](../engine/overview.md).

```
engine.go       package doc and Version
options.go      the Options type and validation
result.go       the Result type
transform.go    wraps api.Transform
build.go        wraps api.Build
errors.go       sentinel errors
```

`engine` is the only package that imports `github.com/evanw/esbuild`.

### `internal/config/`

Config file loading.

```
config.go       the Config type
defaults.go     Defaults()
json.go         file discovery and JSON parsing
validation.go   Validate()
```

### `internal/cache/`

Optional on-disk result cache, keyed by content hash.

```
cache.go        the Cache interface
memory.go       in-memory implementation
filesystem.go   on-disk implementation
key.go          the key derivation
```

### `internal/diagnostics/`

Structured errors and their formatters.

```
diagnostic.go   the Diagnostic type
severity.go     the Severity type
formatter.go    the Formatter interface and sort helpers
terminal.go     the terminal formatter
json.go         the JSON formatter
```

### `internal/source/`

Byte offsets, line numbers, spans. Used by diagnostics and by
future source map work.

```
position.go     Position type
span.go         Span type
location.go     Location type
line_map.go     LineMap type
source.go       Source type
source_file.go  SourceFile type
```

### `internal/io/`

Filesystem and stdio helpers.

```
reader.go       ReadAll, ReadFile, ReadStdin
writer.go       WriteStdout, WriteAll
filesystem.go   Exists, IsFile, IsDir
atomic.go       WriteFileAtomic
stdin.go        package-level Stdin
stdout.go       package-level Stdout, Stderr
```

### `internal/version/`

The version string. Nothing else.

## The tools module

```
tools/
├── go.mod                   module github.com/minifyjs/minifyjs/tools
├── fixture_runner/          runs fixtures, reports pass/fail
├── syntax_matrix/           produces the syntax table
├── benchmark_report/        compares two history runs
├── release/                 release checklists (Python)
└── docs/                    docs generators (Python)
```

`tools` is a separate Go module so its dependencies cannot leak
into the released binary.

## The Python package

```
python/minifyjs/
├── __init__.py              public surface re-exports
├── minifier.py              minify()
├── optimizer.py             optimize()
├── bundler.py               bundle()
├── adapter.py               the Adapter base class
├── options.py               Options, BundleOptions
├── result.py                Result
├── errors.py                exception hierarchy
├── diagnostics.py           Diagnostic
├── config.py                config file loading
├── sourcemap.py             V3 source map helpers
├── cli.py                   python -m minifyjs entry point
├── _binary.py               locates the bundled binary
├── _runner.py               runs the binary as a subprocess
├── _protocol.py             builds CLI arguments
├── _platform.py             OS/arch detection
├── _paths.py                binary path resolution
├── _errors.py               internal message builders
└── _version.py              version string
```

Files whose names start with `_` are internal. Files without the
underscore are public and their contents are re-exported from
`__init__.py`.

## The Python bridge

The Python package does not re-implement minification. It builds
the right CLI arguments, runs the bundled binary as a subprocess,
and parses the result.

```
Python              Native binary
   |                     |
   |  subprocess.run([binary, *args], input=source, ...)
   |-------------------->|
   |                     |  parse, optimize, generate
   |                     |
   |<--------------------|
   |  stdout: minified code
   |  stderr: diagnostics
```

The `_protocol.py` file owns the argument construction. The
`_runner.py` file owns the subprocess call. Everything else is
plumbing.

## Testing layout

```
core/internal/*_test.go         unit tests, in-package
core/integration/cli/           CLI subprocess tests
core/integration/minify/        fixture-driven minify tests
core/integration/bundle/        bundle-mode tests
core/integration/api/           api package end-to-end
core/fuzz/                      Go native fuzz tests (future)
python/tests/                   Python unit and integration tests
integrations/generic/tests/     adapter example tests
test/e2e/                       cross-cutting scenarios (future)
test/compatibility/             ES version matrix (future)
```

The Go tests live next to the code they test. The Python tests are
in a single `tests/` directory because the package is small enough
that per-module test files would outnumber the modules.

## Dependency direction

```
      api           cmd
       |             |
       +------+------+
              |
           engine
              |
           esbuild
```

`api` and `cmd` are siblings. Both call `engine`. Neither calls the
other. `engine` calls esbuild and nothing else.

`internal/config`, `internal/cache`, `internal/diagnostics`,
`internal/source`, `internal/io`, and `internal/version` are leaf
packages. They do not import anything in this repository except
each other (and rarely). `diagnostics` and `source` are used by
`engine`; `config`, `cache`, and `io` are used by `cmd`.

Nothing in `internal/` imports `api`. That is the dependency
inversion: `api` may change freely without affecting `internal/`.

## The most important boundary

`core/internal/engine` is the only package that imports esbuild.
If you find yourself wanting to import esbuild from another
package, the change belongs in `engine` first.

This rule is enforced by review. A future CI check could enforce it
by grepping for `github.com/evanw/esbuild` outside `internal/engine`.

## Adding a package

New packages go in `core/internal/`. The rules:

1. The package must have a doc comment explaining what it does and
   who calls it.
2. The package must not import esbuild unless it is `engine`.
3. The package must not import `api`.
4. The package must have tests.

If the package is a utility that ships to end users, it belongs in
`api` and `internal/`. If it is a developer utility, it belongs in
`tools/`.

## See also

- [../architecture.md](../architecture.md) — the conceptual
  overview
- [../engine/overview.md](../engine/overview.md) — the engine
  package
- [testing.md](testing.md) — how to write tests
- [../../CONTRIBUTING.md](../../CONTRIBUTING.md) — the contribution
  workflow