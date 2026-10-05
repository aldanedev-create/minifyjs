# CLI: configuration

How to make MinifyJS remember your settings, so you do not have to
type the same flags on every invocation.

## Config file discovery

MinifyJS looks for a config file in this order:

1. `minifyjs.config.json`
2. `.minifyjsrc.json`
3. `.minifyjsrc`

Starting from the current directory, walking up to the filesystem
root. The first file found wins. Only one config file is loaded per
invocation.

This means a project layout like:

```
/home/you/project/
├── minifyjs.config.json
├── src/
│   └── app.js
└── static/
    └── styles.css
```

will find the config no matter where inside `project/` you run
`minifyjs` from.

## Overriding discovery

To use a specific file:

```console
minifyjs --config path/to/config.json app.js -o app.min.js
```

To ignore any config file:

```console
minifyjs --no-config app.js -o app.min.js
```

`--no-config` takes precedence over `--config`.

## Precedence

Settings are resolved in this order, highest priority first:

1. **Command-line flags.** `--target es2015` wins over everything.
2. **Environment variables.** `MINIFYJS_TARGET=es2015` wins over the
   config file.
3. **Config file.** Whatever the discovered (or `--config`-specified)
   file says.
4. **Built-in defaults.** The values MinifyJS uses if nothing else
   is set.

## The config file format

A JSON object with a fixed set of keys. Every key is optional.

```json
{
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

The reference for every key is in
[../reference/configuration.md](../reference/configuration.md).

## Path resolution

Paths in the config file are resolved **relative to the config
file's own directory**, not to the current working directory.

```json
// /home/you/project/minifyjs.config.json
{
  "output": "dist/app.min.js"
}
```

Running `minifyjs src/app.js` from `/home/you/project/` writes to
`/home/you/project/dist/app.min.js`.

Running the same command from `/home/you/project/src/` **also**
writes to `/home/you/project/dist/app.min.js`, because the config's
directory is the base. The `src/` prefix in the first run and the
lack of it in the second are both interpreted the same way.

This matters for monorepos and for any workflow that changes
directories between runs.

## Environment variables

Every config key has an environment variable equivalent:

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
insensitive) as true.

Example:

```console
MINIFYJS_TARGET=es2015 MINIFYJS_QUIET=1 minifyjs app.js -o app.min.js
```

## Caching

```console
minifyjs app.js --cache -o app.min.js
```

Enables the on-disk cache. The cache key is a SHA-256 of the input
and every option that affects output. A cache hit skips the esbuild
call entirely.

The default cache directory is:

- **Linux/macOS:** `$XDG_CACHE_HOME/minifyjs` or `~/.cache/minifyjs`
- **Windows:** `%LOCALAPPDATA%\minifyjs\Cache`

Override with `--cache-dir DIR` or `MINIFYJS_CACHE_DIR`.

To disable:

```console
minifyjs app.js --no-cache -o app.min.js
```

The cache is off by default. Turn it on for build pipelines where
you run MinifyJS repeatedly on mostly-unchanged files.

## Where to put the config file

**Project root.** The conventional place. Committing the file to
the repository ensures every developer and every CI run uses the
same settings.

**Home directory.** If you want the same settings for every project
you run MinifyJS in, put a config in `~/.config/` or equivalent. Do
not do this unless you mean it; a home-directory config affects
every project.

**Do not** put a config file in a directory that gets committed by
mistake, like a build output directory. MinifyJS will pick it up
from there and use it for unrelated runs.

## Adding a schema for editor support

The config file can reference a JSON schema for editor
autocomplete and validation:

```json
{
  "$schema": "./docs/reference/config-schema.json",
  "target": "es2015"
}
```

The schema is published at `docs/reference/config-schema.json` in
this repository. Any editor that supports JSON Schema (VS Code,
JetBrains, Vim with the right plugin) will use it.

## Recipes

### Same options for a whole team

Commit a `minifyjs.config.json` at the project root:

```json
{
  "minify": { "whitespace": true, "identifiers": true, "syntax": true },
  "target": "es2015",
  "sourcemap": "external"
}
```

Then `minifyjs src/app.js -o dist/app.min.js` uses those settings
everywhere.

### Different options for dev and production

Use two config files and pick one with `--config`:

```console
# Development: readable output, inline maps.
minifyjs --config minifyjs.dev.json app.js -o app.min.js

# Production: mangled, no maps, banner.
minifyjs --config minifyjs.prod.json app.js -o app.min.js
```

Or use environment variables to switch at runtime:

```sh
# In dev
export MINIFYJS_SOURCEMAP=inline

# In prod (in a different shell)
export MINIFYJS_TARGET=es2015
export MINIFYJS_NO_MINIFY=       # unset
```

### Turn off minification temporarily

```console
minifyjs app.js --no-minify -o app.debug.js
```

Useful when a minified file fails and you want to see the same
output without optimization.

## See also

- [../reference/configuration.md](../reference/configuration.md) —
  every config key, in detail
- [../reference/config-schema.json](../reference/config-schema.json) —
  the JSON schema
- [errors.md](errors.md) — what happens when a config file is invalid