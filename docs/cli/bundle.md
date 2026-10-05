# CLI: bundle

Multi-file project support. Resolves imports, tree-shakes unused
exports, and writes one or more output files.

## When you need it

The default mode (`minifyjs app.js`) operates on a single string. It
cannot resolve `import { x } from "./lib.js"`. If your project has
multiple files that import each other, you need `--bundle`.

## Basic invocation

```console
minifyjs --bundle src/main.js --outdir dist --format esm
```

This resolves `main.js`, finds everything it imports, drops unused
exports, and writes the result to `dist/`.

## Required flags

| Flag | Why |
|---|---|
| `--bundle` | Switches to bundle mode |
| an entry point | What to start from |
| `--outdir DIR` or `--outfile FILE` | Where to write output |

Without `--outdir` or `--outfile`, MinifyJS has nowhere to write the
output and exits with code 2.

## `--outdir` vs `--outfile`

| Flag | When to use it |
|---|---|
| `--outdir DIR` | Multiple entry points, or code splitting. Produces one file per entry point plus shared chunks. |
| `--outfile FILE` | One entry point, one output file. Simpler. |

You cannot set both. Doing so is a usage error.

### With `--outfile`

```console
minifyjs --bundle src/main.js --outfile dist/bundle.js --format esm
```

Produces exactly one file: `dist/bundle.js`.

### With `--outdir`

```console
minifyjs --bundle src/a.js src/b.js --outdir dist --format esm --splitting
```

Produces `dist/a.js`, `dist/b.js`, and shared chunks that both
import.

## Platform

```console
minifyjs --bundle src/main.js --outdir dist --platform node
```

| Platform | What it does |
|---|---|
| `browser` (default) | Assumes browser globals; warns on `require()` of Node builtins |
| `node` | Assumes Node globals; resolves `node:` builtins as external |
| `neutral` | Neither browser nor Node; resolves nothing implicitly |

## Code splitting

```console
minifyjs --bundle src/a.js src/b.js \
    --outdir dist \
    --format esm \
    --splitting
```

With `--splitting`, common imports between entry points are hoisted
into shared chunks. Each entry point's output `import`s from the
shared chunks.

Requires `--format esm`. Attempting to use `--splitting` with
`--format cjs` is a usage error.

## Minification in bundle mode

Bundle mode uses the same flags as transform mode:

```console
minifyjs --bundle src/main.js \
    --outfile dist/bundle.js \
    --format esm \
    --compress \
    --mangle \
    --target es2015 \
    --sourcemap=external
```

Everything from [optimize.md](optimize.md) applies, plus:

- **Tree shaking.** Unused exports are removed. A function that is
  exported but never imported is dropped from the bundle.
- **Side-effect preservation.** A module with side effects (top-level
  `console.log`, `fetch`, etc.) is kept even if none of its exports
  are used. Use `"sideEffects": false` in `package.json` to tell
  MinifyJS a module has no side effects.

## Source maps

```console
minifyjs --bundle src/main.js \
    --outfile dist/bundle.js \
    --sourcemap=external
```

Writes `dist/bundle.js` and `dist/bundle.js.map`. The map's
`sources` array points back at the original files, so the browser's
dev tools show your source, not the bundle.

With `--splitting`, one map is written per output file.

## Examples

### Simple ESM bundle

```console
minifyjs --bundle src/main.js \
    --outfile dist/bundle.js \
    --format esm \
    --compress --mangle \
    --target es2015
```

### Browser bundle with code splitting

```console
minifyjs --bundle src/entry-a.js src/entry-b.js \
    --outdir dist \
    --format esm \
    --splitting \
    --platform browser \
    --compress --mangle
```

### Node CLI tool

```console
minifyjs --bundle bin/cli.js \
    --outfile dist/cli.js \
    --platform node \
    --format cjs \
    --compress --mangle \
    --banner '#!/usr/bin/env node'
```

The `--banner` preserves the shebang through minification.

### Library with type definitions

```console
minifyjs --bundle src/index.js \
    --outfile dist/index.min.js \
    --format esm \
    --sourcemap=external \
    --compress --mangle \
    --legal-comments=external
```

`--legal-comments=external` writes licenses to
`dist/index.min.js.LEGAL.txt` instead of inlining them into the
code.

## What bundle mode is not

- **Not a full bundler.** MinifyJS does not implement loaders for
  CSS, images, fonts, or any non-JS asset. If your project has
  `import "./styles.css"`, bundle mode will fail.
- **Not a plugin host.** esbuild's plugin system is available on the
  `Build` API, which MinifyJS uses, but MinifyJS does not expose a
  way to register plugins. Projects that need plugins should use
  esbuild directly.
- **Not a TypeScript compiler.** It strips types. It does not
  type-check. See [../supported-syntax.md](../supported-syntax.md).

## When to use esbuild directly

If you need any of the following, use esbuild's CLI directly instead
of MinifyJS:

- CSS, image, or font loaders
- A plugin system
- Watch mode with HMR
- Multiple `--loader:` overrides
- `--tsconfig` handling
- Environment variable substitution

MinifyJS is a small wrapper. Its value is the Python distribution
and the no-Node.js guarantee. When those are not the constraint, use
the real tool.

## See also

- [optimize.md](optimize.md) — the flags shared with transform mode
- [configuration.md](configuration.md) — persisting bundle options
- [../python/bundle.md](../python/bundle.md) — the Python equivalent