# CLI overview

The `minifyjs` command is the primary entry point for shell users.
This page is the map; the other pages in this directory are the
territory.

## Invocation forms

```console
minifyjs [INPUT.js] [flags]        # read a file, write to stdout or -o
cat INPUT.js | minifyjs [flags]    # read from stdin
minifyjs --bundle ENTRY [flags]    # resolve imports, write to a directory
```

All three forms share the same flag set. A flag that does not apply
to the current mode is ignored rather than rejected.

## The three modes

| Mode | Triggered by | What it does |
|---|---|---|
| **Transform** | default | Reads one input, produces one output |
| **Bundle** | `--bundle` | Resolves imports, produces one or more outputs |
| **Watch** | `--watch` | Runs transform mode in a loop |

Transform mode is what 95% of users will use. Bundle mode is for
projects with ES module imports. Watch mode is for local
development.

## Flag groups

The flags fall into five groups. Each group has its own page.

| Group | Flags | Page |
|---|---|---|
| **Input/output** | `-o/--output`, input path or stdin | [stdin.md](stdin.md), [stdout.md](stdout.md) |
| **Minification** | `--compress`, `--mangle`, `--no-compress`, `--no-mangle`, `--no-minify` | [minify.md](minify.md), [optimize.md](optimize.md) |
| **Target and format** | `--target`, `--format` | [optimize.md](optimize.md) |
| **Content injection** | `--banner`, `--footer`, `--legal-comments`, `--define`, `--drop`, `--pure` | [optimize.md](optimize.md) |
| **Source maps** | `--sourcemap` | [sourcemaps.md](sourcemaps.md) |
| **Configuration** | `--config`, `--no-config` | [configuration.md](configuration.md) |
| **Bundling** | `--bundle`, `--outdir`, `--outfile`, `--platform`, `--splitting` | [bundle.md](bundle.md) |
| **Caching** | `--cache`, `--no-cache`, `--cache-dir` | [configuration.md](configuration.md) |
| **Output verbosity** | `-q/--quiet`, `--verbose`, `-h/--help`, `-V/--version` | below |
| **Mode** | `--watch` | below |

## Meta flags

| Flag | What it does |
|---|---|
| `-h`, `--help` | Print usage and exit 0 |
| `-V`, `--version` | Print version and exit 0 |
| `-q`, `--quiet` | Suppress non-error output on stderr |
| `--verbose` | Print extra detail on stderr |
| `--watch` | Re-run on input changes |

## Exit codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 1 | Processing error (bad input, I/O, esbuild error) |
| 2 | Bad usage (unknown flag, missing value) |
| 3 | Input file not found |
| 4 | Internal error |

Full list: [../reference/exit-codes.md](../reference/exit-codes.md).

## A minimal invocation

```console
minifyjs app.js -o app.min.js
```

Reads `app.js`, removes whitespace and comments, writes the result
to `app.min.js`. It will use a config file if one is present. That
is the entire default behavior.

## What gets printed where

**stdout** receives exactly the minified bytes. There is no banner
line, no trailing newline, no decoration of any kind. This matters
when the output is piped or captured.

**stderr** receives progress information ("wrote 187 bytes to
app.min.js"), diagnostics, and errors. It is safe to redirect
stdout to a file while leaving stderr on the terminal.

## Config file discovery

MinifyJS looks for a config file in this order:

1. `minifyjs.config.json`
2. `.minifyjsrc.json`
3. `.minifyjsrc`

Starting from the current directory, walking up to the filesystem
root. The first file found wins. See
[configuration.md](configuration.md) for the full details.

## Where to go next

- [minify.md](minify.md) — the simplest mode
- [optimize.md](optimize.md) — the full optimization pipeline
- [bundle.md](bundle.md) — multi-file projects
- [reference.md](reference.md) — every flag, generated from `--help`