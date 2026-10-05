# Architecture

MinifyJS is small. This page explains every part of it.

## The one-paragraph version

MinifyJS is a native binary that wraps
[esbuild](https://esbuild.github.io/)'s Go API. It exposes that API
through a command-line tool and a Python package. Both entry points
call the same Go code, which calls the same esbuild code, which does
all of the actual parsing, optimization, and code generation.
MinifyJS's job is the wrapper: a stable Python surface, sane CLI
defaults, config file discovery, and a distribution model that does
not require Node.js.

## The pipeline

```
              ┌──────────────────┐
              │   User input     │
              │  (.js file, or   │
              │   stdin, or a    │
              │   Python string) │
              └────────┬─────────┘
                       │
        ┌──────────────┴──────────────┐
        │                             │
        ▼                             ▼
   ┌─────────┐                  ┌──────────┐
   │   CLI   │                  │  Python  │
   │ (Go)    │                  │  API     │
   └────┬────┘                  └────┬─────┘
        │                             │
        │  (both call the same        │
        │   Go API in process,        │
        │   not through a socket)     │
        │                             │
        └──────────────┬──────────────┘
                       │
                       ▼
              ┌─────────────────┐
              │  core/api       │  ← stable public surface
              └────────┬────────┘
                       │
                       ▼
              ┌─────────────────┐
              │  core/internal/ │
              │     engine      │  ← the only place that imports esbuild
              └────────┬────────┘
                       │
                       ▼
              ┌─────────────────┐
              │    esbuild      │  ← parses, optimizes, generates
              │  (Go library)   │
              └─────────────────┘
```

## Why wrap esbuild

MinifyJS is not a from-scratch JavaScript minifier. It is a wrapper
around one. That decision is deliberate, and it is the single most
important architectural choice in the project.

The alternative — writing a JavaScript parser, AST, scope analyzer,
optimizer, mangler, and code generator from scratch — is a two-year
effort for a team, and produces a tool that is slower and less
correct than esbuild on day one. esbuild's source is MIT-licensed,
its Go API is stable, and its performance is best-in-class.

MinifyJS's contribution is not the engine. It is the **Python
distribution**, the **CLI ergonomics**, and the **no-Node.js
guarantee**. Those are real problems for real users, and they are
solvable in a small codebase.

See [design-principles.md](design-principles.md) for the full
reasoning.

## The Go module structure

```
core/                          the released Go module
├── cmd/minifyjs/              the CLI's main package
├── api/                       the stable public surface
│   ├── minify.go
│   ├── optimize.go
│   ├── bundle.go
│   ├── options.go
│   ├── result.go
│   └── diagnostics.go
└── internal/
    ├── engine/                the wrapper around esbuild
    ├── config/                config file discovery and validation
    ├── cache/                 optional on-disk result cache
    ├── diagnostics/           structured errors and formatters
    ├── source/                byte offsets, line/column bookkeeping
    ├── io/                    atomic writes, stdin/stdout helpers
    └── version/               the version string
```

Two rules govern the layout:

1. **`internal/` is private.** No package under `internal/` may be
   imported by code outside the `core` module. The Go compiler
   enforces this.

2. **`api/` is the only public surface.** External Go callers import
   `github.com/minifyjs/minifyjs/core/api`. They never reach into
   `internal/`. This means every `internal/` package is free to
   change shape without breaking callers.

## The Python module structure

```
python/minifyjs/
├── __init__.py                re-exports the public surface
├── minifier.py                minify()
├── optimizer.py               optimize()
├── bundler.py                 bundle()
├── adapter.py                 the Adapter base class
├── options.py                 Options and BundleOptions
├── result.py                  Result
├── errors.py                  the exception hierarchy
├── diagnostics.py             Diagnostic
├── config.py                  config file loading
├── sourcemap.py               V3 source map helpers
├── cli.py                     the python -m minifyjs entry point
├── _binary.py                 locates the bundled native binary
├── _runner.py                 runs the native binary as a subprocess
├── _protocol.py               builds CLI arguments from Options
├── _platform.py               OS and architecture detection
├── _paths.py                  resolves the binary's on-disk location
├── _errors.py                 internal error message builders
└── _version.py                the version string
```

The Python package never re-implements minification. It builds the
right CLI arguments, invokes the bundled binary as a subprocess over
stdin/stdout, and parses the result.

## The distribution model

End users never install Go, clone the repository, or compile
anything. The Go engine is compiled once by the release pipeline and
shipped inside platform-specific Python wheels:

```
Go source
   │  go build -trimpath -ldflags ...
   ▼
native binary (minifyjs or minifyjs.exe)
   │  build_wheels.py
   ▼
Python wheel (one per OS/architecture)
   │  twine upload
   ▼
PyPI
   │  pip install minifyjs
   ▼
user's machine
```

The result: a user runs `pip install minifyjs` and gets a working
minifier with no toolchain. See [getting-started/installation.md](getting-started/installation.md).

## What MinifyJS does not do

- **It does not type-check TypeScript.** esbuild strips types; it
  does not verify them. If you pass a `.ts` file, MinifyJS will
  happily produce output with type errors in it. Run `tsc --noEmit`
  separately.
- **It does not bundle in the default (single-file) mode.** `bundle()`
  is a separate API. `minify()` and `optimize()` operate on one
  string.
- **It does not minify CSS, HTML, or JSON.** The engine is a
  JavaScript minifier. esbuild can handle CSS but MinifyJS does not
  expose that.
- **It does not offer plugin hooks.** esbuild's plugin system is
  available on its `Build` API but not its `Transform` API, and
  MinifyJS deliberately exposes only the latter for single-file work.

## The engine boundary

The single most important interface in the codebase is
`core/internal/engine`. It is the only package that imports
`github.com/evanw/esbuild`. Every other package goes through it.

This means:

- If esbuild's Go API changes shape, exactly one package needs to be
  updated.
- If a caller wants to know "what esbuild options does MinifyJS
  actually set?", the answer is in exactly one file
  (`internal/engine/options.go`).
- If MinifyJS ever needs to support a second engine (it will not),
  the interface is already in place.

The engine package is documented in [engine/overview.md](engine/overview.md).

## Testing

Four layers:

| Layer | Where | What it verifies |
|---|---|---|
| Unit tests | `core/internal/*_test.go`, `python/tests/test_*.py` | Individual functions |
| Integration tests | `core/integration/` | The compiled CLI as a subprocess |
| End-to-end tests | `test/e2e/` | Full workflows (build, serve, verify) |
| Compatibility tests | `test/compatibility/` | ES version support from es2015 to es2025 |

Plus two cross-cutting mechanisms:

- **Idempotence tests** at every level: minifying already-minified
  code produces identical output.
- **Cross-language agreement tests**: the Go CLI and the Python API
  produce byte-identical output for the same input and options.

See [development/testing.md](development/testing.md) for the full
picture.

## What is deliberately absent

- **A plugin system.** MinifyJS's surface is small on purpose.
  Plugins would require exposing esbuild internals that it does not
  expose for the `Transform` path.
- **A web dashboard.** MinifyJS is a build step, not a service.
  Users who need visual output should use their CI's logging or a
  dedicated tool.
- **A REPL.** The Python `-c` one-liner is the REPL.
- **A watch mode beyond a simple polling loop.** Polling works on
  every filesystem, including network mounts and Docker volumes.
  inotify and friends do not.