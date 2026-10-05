# Changelog

All notable changes to MinifyJS are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial project structure: Go engine wrapper around esbuild's Go
  API, native CLI, public Go API, and Python package.
- `minifyjs` CLI with the following flags: `-o/--output`,
  `--compress`, `--no-compress`, `--mangle`, `--no-mangle`,
  `--no-minify`, `--target`, `--format`, `--sourcemap`, `--banner`,
  `--footer`, `--legal-comments`, `--define`, `--drop`, `--pure`,
  `--bundle`, `--outdir`, `--platform`, `--splitting`, `--cache`,
  `--no-cache`, `--cache-dir`, `--config`, `--no-config`, `--watch`,
  `-q/--quiet`, `--verbose`, `-h/--help`, `-V/--version`.
- Config file discovery (`minifyjs.config.json`, `.minifyjsrc.json`,
  `.minifyjsrc`) with walk-up search from the current directory.
- Environment variable overlay (`MINIFYJS_*`).
- On-disk cache keyed by content hash, opt-in via `--cache`.
- Watch mode using polling (works on network mounts and inside
  Docker volumes).
- Python package (`pip install minifyjs`) that bundles the native
  binary and exposes `minify()`, `optimize()`, and `bundle()`.
- Structured diagnostics with terminal and JSON formatters.
- Exit codes: 0 success, 1 processing error, 2 bad usage, 3 file not
  found, 4 internal error.

### Changed

- The engine is implemented as a wrapper around esbuild's Go API
  rather than as a from-scratch JavaScript parser and optimizer.
  See `docs/architecture.md` for the rationale.

### Removed

- Earlier prototype implementations of a hand-written lexer,
  parser, AST, and printer. esbuild now owns those stages.



## [0.1.2] - 2026-10-05

- fix metadata build and cli entry point user can now used
minifyjs input.js -o output.min.js instead of python -m minifyjs input.js -o output.min.js


## [0.1.1] - 2026-10-05

- changes to readme.md and docs and version 0.1.1

## [0.1.0] - 2026-10-05

- Initial 0.1.0 release of the native CLI, Go API, Python package,
  bundler, diagnostics, source maps, and cache support.

## [0.0.0] - 2026-01-01

- Empty placeholder release to establish the changelog.

[Unreleased]: https://github.com/aldanedev-create/minifyjs/compare/v0.1.0...HEAD