# Configuration reference

Every key in `minifyjs.config.json`, with its type, default, and
behavior.

## File names

MinifyJS searches for these file names, in this order:

1. `minifyjs.config.json`
2. `.minifyjsrc.json`
3. `.minifyjsrc`

The first one found wins. The search starts at the current directory
and walks up to the filesystem root.

## Precedence

Settings are resolved in this order, highest priority first:

1. **Command-line flags** (`--target es2015`)
2. **Environment variables** (`MINIFYJS_TARGET=es2015`)
3. **Config file** (the discovered file)
4. **Built-in defaults** (documented below)

## Path resolution

Paths in the config file are resolved relative to the config file's
own directory, not to the current working directory. This lets a
config checked into a project work regardless of where MinifyJS is
invoked from.

## The full schema

```json
{
  "$schema": "./docs/reference/config-schema.json",

  "inputs": [],
  "output": "",

  "minify": {
    "whitespace": true,
    "identifiers": true,
    "syntax": true
  },

  "target": "esnext",
  "format": "",
  "sourcemap": "",

  "banner": "",
  "footer": "",
  "legalComments": "eof",

  "bundle": {
    "enabled": false,
    "entryPoints": [],
    "platform": "browser",
    "splitting": false
  },

  "cache": {
    "enabled": false,
    "dir": ""
  },

  "quiet": false,
  "verbose": false
}
```

Every key is optional. A missing key uses its default value.

## Key reference

### `inputs`

- **Type:** array of strings
- **Default:** `[]` (read from stdin)
- **CLI equivalent:** positional arguments

The input files to minify. If empty, MinifyJS reads from stdin.

When `bundle.enabled` is `true`, this is ignored in favor of
`bundle.entryPoints`.

### `output`

- **Type:** string
- **Default:** `""` (write to stdout)
- **CLI equivalent:** `-o`, `--output`

Where to write the minified JavaScript. Empty means stdout.

When `bundle.enabled` is `true`, this is the output directory (like
`--outdir`).

### `minify.whitespace`

- **Type:** boolean
- **Default:** `true`
- **CLI equivalent:** part of `--compress`

Remove whitespace that is not required to keep tokens separate.

### `minify.identifiers`

- **Type:** boolean
- **Default:** `true`
- **CLI equivalent:** `--mangle`

Rename local variables and function parameters to short names.
Names observable from outside the file are preserved.

### `minify.syntax`

- **Type:** boolean
- **Default:** `true`
- **CLI equivalent:** part of `--compress`

Enable constant folding, dead-code elimination, and expression
simplification.

### `target`

- **Type:** string
- **Default:** `"esnext"`
- **CLI equivalent:** `--target`

The ECMAScript version for the output. Valid values:

- `es5`, `es2015`, `es2016`, ..., `es2024`, `esnext`
- Browser + version: `chrome90`, `firefox88`, `safari14`, `edge90`
- Node: `node18`, `node20`
- Compound: `es2020,chrome90,firefox88`

### `format`

- **Type:** string
- **Default:** `""` (preserve)
- **CLI equivalent:** `--format`

The module wrapper for the output.

| Value | Effect |
|---|---|
| `""` | Preserve whatever the input uses |
| `"esm"` | ES module syntax |
| `"cjs"` | CommonJS |
| `"iife"` | Immediately-invoked function expression |

### `sourcemap`

- **Type:** string
- **Default:** `""` (none)
- **CLI equivalent:** `--sourcemap`

Source map mode.

| Value | Effect |
|---|---|
| `""` | No source map |
| `"inline"` | Map embedded in the output |
| `"external"` | Map written to `<output>.map` |
| `"both"` | Both of the above |

### `banner`

- **Type:** string
- **Default:** `""`
- **CLI equivalent:** `--banner`

Text prepended verbatim to the output. Not parsed or minified.

### `footer`

- **Type:** string
- **Default:** `""`
- **CLI equivalent:** `--footer`

Text appended verbatim to the output.

### `legalComments`

- **Type:** string
- **Default:** `"eof"`
- **CLI equivalent:** `--legal-comments`

What happens to `/*! ... */` and `//! ... */` comments.

| Value | Effect |
|---|---|
| `"none"` | Removed |
| `"inline"` | Kept in place |
| `"eof"` | Moved to end of file |
| `"external"` | Written to a separate `.LEGAL.txt` file |

### `bundle.enabled`

- **Type:** boolean
- **Default:** `false`
- **CLI equivalent:** `--bundle`

Switch to bundle mode. Requires `bundle.entryPoints` (or `inputs`)
and either `output` or `bundle.outFile`.

### `bundle.entryPoints`

- **Type:** array of strings
- **Default:** `[]`
- **CLI equivalent:** positional arguments with `--bundle`

The files to bundle. At least one is required when
`bundle.enabled` is `true`.

### `bundle.platform`

- **Type:** string
- **Default:** `"browser"`
- **CLI equivalent:** `--platform`

| Value | Effect |
|---|---|
| `"browser"` | Browser globals; warns on Node builtins |
| `"node"` | Node globals |
| `"neutral"` | Neither |

### `bundle.splitting`

- **Type:** boolean
- **Default:** `false`
- **CLI equivalent:** `--splitting`

Enable code splitting. Requires `format: "esm"`.

### `cache.enabled`

- **Type:** boolean
- **Default:** `false`
- **CLI equivalent:** `--cache`

Enable the on-disk result cache.

### `cache.dir`

- **Type:** string
- **Default:** `""` (OS default)
- **CLI equivalent:** `--cache-dir`

Override the cache directory. Empty means the OS default:
`~/.cache/minifyjs` on Linux/macOS, `%LOCALAPPDATA%\minifyjs\Cache`
on Windows.

### `quiet`

- **Type:** boolean
- **Default:** `false`
- **CLI equivalent:** `-q`, `--quiet`

Suppress non-error output on stderr.

### `verbose`

- **Type:** boolean
- **Default:** `false`
- **CLI equivalent:** `--verbose`

Print extra detail on stderr. Mutually exclusive with `quiet`;
specifying both is a validation error.

## Environment variables

Every config key with a scalar value has an environment variable
equivalent:

| Config key | Environment variable |
|---|---|
| `target` | `MINIFYJS_TARGET` |
| `format` | `MINIFYJS_FORMAT` |
| `sourcemap` | `MINIFYJS_SOURCEMAP` |
| `cache.enabled` | `MINIFYJS_CACHE` |
| `cache.dir` | `MINIFYJS_CACHE_DIR` |
| `quiet` | `MINIFYJS_QUIET` |
| `verbose` | `MINIFYJS_VERBOSE` |
| (turns off all minify passes) | `MINIFYJS_NO_MINIFY` |

Boolean variables accept `1`, `true`, `yes`, or `on` (case
insensitive) as true. Any other value is false.

## Examples

### Minimal — just minify with defaults

```json
{}
```

An empty config file is valid. Every key uses its default. Since
the built-in defaults already enable full minification, an empty
config means "minify with `--compress --mangle`".

### Production

```json
{
  "minify": { "whitespace": true, "identifiers": true, "syntax": true },
  "target": "es2015",
  "drop": ["console", "debugger"],
  "sourcemap": "",
  "quiet": true,
  "cache": { "enabled": true }
}
```

### Development

```json
{
  "minify": { "whitespace": true, "identifiers": false, "syntax": true },
  "target": "esnext",
  "sourcemap": "inline",
  "legalComments": "inline",
  "verbose": true
}
```

### Library

```json
{
  "minify": { "whitespace": true, "identifiers": true, "syntax": true },
  "target": "es2015",
  "format": "esm",
  "banner": "/* lib v1.0.0 | MIT */",
  "legalComments": "inline",
  "sourcemap": "external"
}
```

### Bundle

```json
{
  "bundle": {
    "enabled": true,
    "entryPoints": ["src/main.js"],
    "platform": "browser",
    "splitting": true
  },
  "output": "dist",
  "format": "esm",
  "target": "es2015"
}
```

## See also

- [`config-schema.json`](config-schema.json) — the JSON Schema for
  editor autocomplete
- [../cli/configuration.md](../cli/configuration.md) — the CLI's
  view
- [../python/configuration.md](../python/configuration.md) — the
  Python view