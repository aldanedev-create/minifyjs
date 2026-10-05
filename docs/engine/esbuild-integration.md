# esbuild integration

How MinifyJS uses esbuild, and what that means for users.

## The dependency

MinifyJS depends on exactly one Go module:

```
github.com/evanw/esbuild v0.28.2
```

It is pinned in `core/go.mod`. It is not a build-time tool. It is
compiled into the MinifyJS binary alongside MinifyJS's own code.

There is no subprocess. There is no `npx esbuild`. There is no
`node_modules/`. The esbuild source is linked into the MinifyJS
binary as a Go library, exactly like any other Go dependency.

## What esbuild does in MinifyJS

esbuild does all of the following:

- **Lexing.** Turning source text into tokens.
- **Parsing.** Turning tokens into an AST.
- **Scope analysis.** Tracking which identifier refers to which
  binding.
- **Optimization.** Constant folding, dead-code elimination,
  expression simplification.
- **Mangling.** Renaming local variables to short names.
- **Code generation.** Emitting JavaScript text from the AST.
- **Source map generation.** Producing V3 maps.
- **Tree shaking.** Dropping unused exports in bundle mode.
- **Module resolution.** Finding imports on disk in bundle mode.
- **Syntax lowering.** Rewriting newer syntax for older targets.

That is the entire compiler. MinifyJS adds a wrapper around it.

## What MinifyJS does

- **Argument parsing.** The CLI's flags are MinifyJS's.
- **Config file discovery.** esbuild does not have a config file
  format. MinifyJS's config loading is its own.
- **Diagnostic formatting.** esbuild's messages are converted to
  MinifyJS's `Diagnostic` type, then rendered as terminal text or
  JSON.
- **Python packaging.** esbuild's Go library is not on PyPI. The
  MinifyJS release pipeline compiles the binary and ships it inside
  a platform-specific wheel.
- **Convenience wrappers.** `optimize()` is `minify(compress=True,
  mangle=True)`. `Adapter` walks directories. `bundle()` writes
  files.

The wrapper is about 1000 lines of Go and about 1500 lines of
Python. esbuild is several hundred thousand lines.

## Why not use esbuild's CLI directly

You can. `esbuild --minify app.js --outfile=app.min.js` works. If
that command does what you want, use it.

The reasons to use MinifyJS instead:

1. **No Node.js.** esbuild's CLI requires Node.js, because its
   wrapper is a Node script. esbuild's Go library does not. MinifyJS
   exposes the Go library through a native binary.

2. **A Python API.** esbuild does not have one. MinifyJS does.

3. **A config file.** esbuild has no config file format. MinifyJS
   has `minifyjs.config.json`.

4. **Sane defaults.** esbuild's defaults are reasonable for
   bundling; MinifyJS's defaults are reasonable for minifying one
   file.

5. **Diagnostic stability.** esbuild's error messages change between
   versions. MinifyJS's diagnostic format is stable.

None of those are huge. Together they are the reason MinifyJS
exists.

## Which esbuild APIs MinifyJS uses

Two:

- `api.Transform` — single-file transformation. Used by
  `engine.Transform`, which is called by `minify()` and
  `optimize()`.
- `api.Build` — multi-file build. Used by `engine.Build`, which is
  called by `bundle()`.

Two APIs MinifyJS does not use:

- `api.Context` — esbuild's watch/rebuild API. MinifyJS's
  `--watch` uses polling on top of `Transform` instead.
- The plugin API. Plugins are only available on `api.Build` and
  MinifyJS does not expose them.

## The dependency surface

MinifyJS uses about **10 options** out of the roughly **40** that
`api.TransformOptions` accepts. It uses about **15 options** out of
the roughly **80** that `api.BuildOptions` accepts.

That means MinifyJS is not a full wrapper. It is a wrapper for the
options that most single-file minification workflows need.

Options that MinifyJS does not expose, and why:

| Option | Why not |
|---|---|
| `JSXFactory`, `JSXFragment` | Rarely needed; defaults are fine |
| `TsconfigRaw` | Full TypeScript projects need `tsc`, not a minifier |
| `Define` | Exposed via `--define` in the CLI, not in the API |
| `Pure` | Exposed via `--pure` in the CLI |
| `Drop` | Exposed via `--drop` in the CLI |
| `MangleProps` | Would break code that depends on property names |
| `KeepNames` | Would defeat the purpose of mangling |
| `Charset` | Always UTF-8 |
| `LineLimit` | Output is single-line by design |
| `TreeShaking` | Always on in bundle mode, off in transform mode |

If you need any of those, use esbuild directly.

## Upgrading esbuild

Upgrading esbuild is a minor-version change for MinifyJS. The steps:

1. Bump the version in `core/go.mod`.
2. Run `go mod tidy`.
3. Run the full test suite. `go test ./...` in `core/`.
4. Run the Python test suite. `pytest` in `python/`.
5. Run the benchmark suite. `python bench/run_all.py` in `bench/`.
6. Compare the new results against the previous `latest.md`.
7. Update `CHANGELOG.md` with a note about the version bump.

If the benchmark shows that output got smaller, that is expected
and welcome. If output got *larger*, investigate before merging.

If the test suite fails, the failure is either:

- A real behavior change in esbuild. Update the tests if the new
  behavior is correct, or file a bug with esbuild if it is not.
- A MinifyJS bug that the new esbuild version exposed. Fix it.

## Pinning vs. floating

MinifyJS pins a specific esbuild version. It does not use a range.

This means a MinifyJS release always produces the same output on
every machine. A user who upgrades MinifyJS from 0.1.0 to 0.1.1
might get different output if the underlying esbuild version
changed, but two installs of 0.1.0 always agree.

Reproducible builds depend on this. So do cache keys, which hash
both the input and the options but not the engine version (the
assumption being that the engine does not change under a fixed
version).

## Licensing

esbuild is MIT-licensed. MinifyJS is MIT-licensed. The full esbuild
license is reproduced in the repository's `LICENSE` file under the
"Acknowledgements" section.

Compiling esbuild into the MinifyJS binary is permitted by MIT and
is standard practice.

## Reporting bugs

If a bug is in esbuild's parsing, optimization, or code generation,
it will affect anyone using esbuild, not just MinifyJS users.
Report it to the esbuild project:

<https://github.com/evanw/esbuild/issues>

If a bug is in MinifyJS's option mapping, error reporting, or
wrapper behavior, report it to MinifyJS. The `minifyjs --version`
output names the esbuild version, which is usually the first thing
a maintainer will ask for.

## See also

- [overview.md](overview.md) — the engine package
- [options-mapping.md](options-mapping.md) — how each option maps
- [../architecture.md](../architecture.md) — the full pipeline