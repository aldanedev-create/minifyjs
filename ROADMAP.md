# Roadmap

This document describes what MinifyJS has shipped and what it plans
to ship. It is not a promise. Priorities change.

## Current status

MinifyJS is at **0.1.0**. It is usable, but pre-1.0: the on-disk
format of the cache and the exact set of flags may still change
between minor versions.

## What is done

- Native Go binary wrapping esbuild's Go API.
- CLI with full flag coverage for esbuild's transform options.
- Python package (`pip install minifyjs`) bundling the native binary
  and exposing `minify()`, `optimize()`, and `bundle()`.
- Config file discovery (`minifyjs.config.json` and variants).
- Environment variable overlay.
- On-disk content-addressed cache (opt-in).
- Watch mode (polling-based).
- Structured diagnostics with terminal and JSON formatters.
- Multi-platform wheel distribution.

## What is next

### 0.2.0 — Diagnostics and config

- `--diagnostics-json` flag for structured output consumed by the
  Python bridge.
- Config schema published at a stable URL, referenced by
  `minifyjs.config.json`'s `$schema` field.
- Source-snippet rendering in diagnostics (the "here's the line,
  here's a caret" output).
- Stable diagnostic codes documented in
  `docs/reference/error-codes.md`.

### 0.3.0 — Bundling improvements

- Bundle manifest output (`.minifyjs-manifest.json`).
- `--splitting` improvements: shared-chunk naming, entry-point
  metadata.
- Documented recipe for using MinifyJS from a Python build script
  (Django collectstatic, Flask asset pipeline, etc.).

### 0.4.0 — Watch mode improvements

- Debounce on save events (an editor that saves twice in quick
  succession currently triggers two rebuilds).
- Incremental rebuilds using esbuild's context API rather than
  re-running Transform.
- `--watch` reports compile time per file.

### 1.0.0 — Stability

- Freeze the on-disk cache format.
- Freeze the config file schema.
- Freeze the CLI flag set.
- Publish API stability guarantees for the Go and Python surfaces.

## What is deliberately not planned

These have been proposed and declined. They are listed so the
reasoning is visible.

- **A from-scratch JavaScript engine.** esbuild already exists, is
  faster than anything we would write, and is MIT-licensed. See
  `docs/architecture.md`.
- **A plugin system.** Plugins would require exposing esbuild's
  internals, which its API does not support for the transform path.
  Users who need plugin-level control should use esbuild directly.
- **TypeScript, JSX, or CSS support.** MinifyJS is a JavaScript
  minifier. Supporting other languages would change what it is.
- **A Node.js CLI wrapper.** The whole point of MinifyJS is to
  avoid requiring Node.js.
- **A GUI.** MinifyJS is a build-tool dependency. GUIs are a
  different product.

## How to influence the roadmap

Open a GitHub Discussion describing the problem you want solved.
Concrete problems get addressed; abstract feature requests usually
do not. See [CONTRIBUTING.md](CONTRIBUTING.md) for what is in scope.