# Getting started: the CLI

Everything you can do from the `minifyjs` command.

## Basic invocation

```console
minifyjs INPUT.js -o OUTPUT.js
```

Reads `INPUT.js`, writes `OUTPUT.js`. If `-o` is omitted, the output
goes to stdout.

## Reading from stdin

If no input file is given, MinifyJS reads from stdin:

```console
cat app.js | minifyjs > app.min.js
```

The literal path `-` also means stdin:

```console
minifyjs - < app.js > app.min.js
```

## The three minify levels

| Flags | What it does | Typical reduction |
|---|---|---|
| (none) | Removes whitespace and comments | 30–50% |
| `--compress` | Also folds constants, drops dead code | +10–20% |
| `--compress --mangle` | Also renames local variables | +10–15% |

Most users want the last one. `--compress --mangle` is what
`optimize()` does in the Python API.

## Overriding the defaults

`--no-compress` and `--no-mangle` turn off the corresponding pass.
They exist so a config file that enables a pass can be overridden
on the command line:

```console
# Config file says mangle=true, but I want readable output this time.
minifyjs app.js --no-mangle -o app.debug.js
```

`--no-minify` turns off every pass. The output is esbuild's
normalized version of the input — same code, reformatted.

## Picking a target

```console
minifyjs app.js --target es2015 -o app.min.js
```

Valid targets: `es5`, `es2015` through `es2024`, `esnext`, and
compound forms like `es2020,chrome90`.

Default is `esnext`, which means "do not lower anything". If you
need older-browser support, you have to ask for it.

## Source maps

```console
minifyjs app.js --sourcemap=external -o app.min.js
```

Three modes:

| Mode | What happens |
|---|---|
| `--sourcemap=inline` | The map is embedded in the output as a data URL |
| `--sourcemap=external` | The map is written to `app.min.js.map` |
| `--sourcemap=both` | Both of the above |

The bare flag `--sourcemap` (no value) means `external`.

## Output format

```console
minifyjs app.js --format esm -o app.min.js
```

Valid formats: `iife`, `cjs`, `esm`. The default (empty) means
"preserve whatever the input uses".

## Content injection

Prepend a banner and append a footer:

```console
minifyjs app.js --banner "/* (c) 2026 */" --footer "// end" -o app.min.js
```

These are added verbatim. They are not parsed, minified, or
otherwise touched. They are useful for license headers and
build-stamp comments.

## Legal comments

```console
minifyjs app.js --legal-comments=eof -o app.min.js
```

Controls what happens to `/*! ... */` and `//! ... */` comments
(the ones with the extra `!` that conventionally mean "preserve
this"). Valid modes:

| Mode | Behavior |
|---|---|
| `none` | Remove all legal comments |
| `inline` | Keep them where they are |
| `eof` (default) | Move them all to the end of the file |
| `external` | Write them to a separate file |

## Output to a file vs stdout

```console
# To a file
minifyjs app.js -o app.min.js

# To stdout
minifyjs app.js

# To stdout, piped
minifyjs app.js | gzip > app.min.js.gz
```

Output to stdout contains exactly the minified bytes, with no added
trailing newline. This matters when the output is compared or
captured.

## Config file discovery

MinifyJS looks for `minifyjs.config.json` (or `.minifyjsrc.json`, or
`.minifyjsrc`) starting from the current directory and walking up to
the filesystem root. The first one found is used.

To use a specific file:

```console
minifyjs --config path/to/config.json app.js -o app.min.js
```

To ignore any config file:

```console
minifyjs --no-config app.js -o app.min.js
```

## Bundling

```console
minifyjs --bundle src/main.js --outdir dist --format esm
```

The `--bundle` flag switches to build mode, which resolves imports
and produces one output file per entry point. See
[cli/bundle.md](../cli/bundle.md).

## Watch mode

```console
minifyjs app.js -o app.min.js --watch
```

Re-runs whenever the input changes. Uses polling (every 250 ms), so
it works on network mounts and inside Docker volumes. Press Ctrl-C
to stop.

## Quiet and verbose

```console
minifyjs app.js -o app.min.js --quiet    # suppress progress
minifyjs app.js -o app.min.js --verbose  # extra detail
```

`--quiet` still prints errors. `--verbose` prints per-file
statistics and esbuild version info.

## Exit codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 1 | Processing error (bad input, I/O error, esbuild error) |
| 2 | Bad command-line usage (unknown flag, missing value) |
| 3 | Input file not found |
| 4 | Internal error (a bug) |

Full list: [reference/exit-codes.md](../reference/exit-codes.md).

## Common patterns

### Minify every `.js` file in a directory

```sh
for f in static/js/*.js; do
    minifyjs "$f" --compress --mangle -o "${f%.js}.min.js"
done
```

### Minify in a Makefile

```makefile
%.min.js: %.js
	minifyjs $< --compress --mangle -o $@
```

### Minify in a GitHub Actions workflow

```yaml
- uses: actions/setup-python@v5
  with:
    python-version: "3.11"
- run: pip install minifyjs
- run: minifyjs static/js/app.js --compress --mangle -o static/js/app.min.js
```

### Check the size before and after

```sh
before=$(wc -c < app.js)
minifyjs app.js --compress --mangle -o app.min.js
after=$(wc -c < app.min.js)
echo "reduced from $before to $after bytes ($((100 - after * 100 / before))% smaller)"
```