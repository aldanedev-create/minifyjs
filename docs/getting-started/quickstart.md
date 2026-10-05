# Quickstart

Five minutes from zero to a minified file.

## Install

```console
pip install minifyjs
```

## Minify a file

```console
minifyjs app.js -o app.min.js
```

That is the whole command. It reads `app.js`, removes whitespace and
comments, and writes the result to `app.min.js`.

## Optimize a file

For a smaller output, add `--compress` and `--mangle`:

```console
minifyjs app.js --compress --mangle -o app.min.js
```

- `--compress` enables constant folding, dead-code elimination, and
  expression simplification.
- `--mangle` renames local variables and function parameters to
  short names.

Together they typically reduce a file by 60–75%.

## Pipe through stdin

MinifyJS reads from stdin when no file is given:

```console
cat app.js | minifyjs --compress > app.min.js
```

## Use it from Python

```python
from minifyjs import optimize

result = optimize("function add(a, b) { return a + b; }")
print(result.code)
# function add(a,b){return a+b;}
```

Or for a file:

```python
from minifyjs import optimize
from pathlib import Path

source = Path("app.js").read_text()
result = optimize(source)
Path("app.min.js").write_text(result.code)

print(f"{result.original_bytes} -> {result.minified_bytes} bytes")
print(f"saved {result.bytes_saved} bytes ({result.ratio * 100:.1f}% of original)")
```

## Target an older browser

By default, MinifyJS does not lower syntax. If your code uses arrow
functions and you need it to run in ES5 environments:

```console
minifyjs app.js --target es5 -o app.min.js
```

The output uses `function` instead of `=>`, `var` instead of `let`,
and so on.

## Generate a source map

```console
minifyjs app.js --sourcemap=external -o app.min.js
```

This writes two files: `app.min.js` and `app.min.js.map`. The
browser's developer tools will load the map automatically when you
inspect the minified file.

For a map that is embedded in the output (no separate file), use
`--sourcemap=inline`.

## Bundle a project

For a multi-file project:

```console
minifyjs --bundle src/main.js --outdir dist --format esm
```

This resolves imports, removes unused exports, and writes the
result to `dist/`. See [cli/bundle.md](../cli/bundle.md) for the
full options.

## Use a config file

Create `minifyjs.config.json` in your project root:

```json
{
  "target": "es2015",
  "minify": {
    "whitespace": true,
    "identifiers": true,
    "syntax": true
  },
  "sourcemap": "external"
}
```

MinifyJS finds this automatically. From then on, `minifyjs app.js
-o app.min.js` uses those settings.

## What's next

- [CLI reference](../cli/overview.md) — every flag and every behavior
- [Python API](../python/api.md) — the full programmatic surface
- [Integrations](../integrations/) — using MinifyJS with Flask,
  Django, or FastAPI
- [Architecture](../architecture.md) — how MinifyJS is built and why