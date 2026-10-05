# Getting started: the Python API

Everything you can do from Python.

## The three entry points

```python
from minifyjs import minify, optimize, bundle
```

| Function | What it does | Analogous CLI |
|---|---|---|
| `minify(source, **opts)` | Removes whitespace and comments | `minifyjs app.js` |
| `optimize(source, **opts)` | Also folds constants and mangles names | `minifyjs app.js --compress --mangle` |
| `bundle(entries, **opts)` | Resolves imports, tree-shakes, writes files | `minifyjs --bundle ...` |

## Minify a string

```python
from minifyjs import minify

result = minify("function add(a, b) { return a + b; }")
print(result.code)
# function add(a,b){return a+b;}

print(f"{result.original_bytes} -> {result.minified_bytes} bytes")
# 40 -> 29 bytes
```

## Optimize a string

```python
from minifyjs import optimize

result = optimize("function calculateTotal(price, tax) { "
                  "  const subtotal = price + (price * tax); "
                  "  return subtotal; }")
print(result.code)
# function calculateTotal(e,t){return e+e*t}
```

`optimize()` is shorthand for `minify(..., compress=True, mangle=True)`.

## Passing options

```python
from minifyjs import minify

result = minify(
    source,
    compress=True,
    mangle=True,
    target="es2015",
    sourcemap="external",
    banner="/* (c) 2026 */",
)

print(result.code)   # the minified JavaScript
print(result.map)    # the source map JSON (if sourcemap was requested)
```

Every keyword argument in `minify()` corresponds to a CLI flag. See
the [full option list](../python/minify.md).

## Handling the result

`minify()` and `optimize()` return a `Result`:

```python
from minifyjs import optimize

result = optimize(source)

print(result.code)              # the minified JavaScript as a string
print(result.map)               # the source map JSON, or ""
print(result.original_bytes)    # length of the input in bytes
print(result.minified_bytes)    # length of the output in bytes
print(result.bytes_saved)       # original_bytes - minified_bytes
print(result.ratio)             # minified_bytes / original_bytes
print(result.has_errors)        # True if any diagnostic is an error
print(result.has_warnings)      # True if any diagnostic is a warning

for d in result.diagnostics:
    print(f"{d.severity}: {d.message}")
```

## Handling errors

`minify()` and `optimize()` raise `MinifyError` when the input cannot
be processed:

```python
from minifyjs import minify, MinifyError

try:
    result = minify("function () { } }")
except MinifyError as e:
    print(f"failed: {e}")
```

The exception message includes the file, line, and column when
available.

For bundling, the exception type is `BundleError`, which is a
subclass of `MinifyError`.

## Bundling

`bundle()` writes output to disk rather than returning a string,
because a bundle can be multiple files:

```python
from minifyjs import bundle

result = bundle(
    ["src/main.js"],
    outfile="dist/bundle.js",
    format="esm",
    target="es2015",
    sourcemap="external",
)

if result.has_errors:
    for d in result.diagnostics:
        print(f"error: {d}")
```

Exactly one of `outdir` or `outfile` must be given.

## Reading from a file

MinifyJS does not do file I/O for the single-file API. Read the file
yourself:

```python
from pathlib import Path
from minifyjs import optimize

source = Path("app.js").read_text(encoding="utf-8")
result = optimize(source)
Path("app.min.js").write_text(result.code, encoding="utf-8")
```

For a whole directory, use the `Adapter` base class:

```python
from minifyjs import Adapter

adapter = Adapter()
results = adapter.minify_directory("static/js")
print(f"minified {len(results)} files")
```

See [integrations/generic-python.md](../integrations/generic-python.md)
for writing your own adapter.

## Using a config file

```python
from minifyjs import load_config, optimize

cfg = load_config("minifyjs.config.json")
result = optimize(source, target=cfg.target, sourcemap=cfg.sourcemap)
```

Or discover the config automatically:

```python
from minifyjs import find_config, load_config

path = find_config()  # walks up from cwd
if path:
    cfg = load_config(path)
```

## Type hints

The package ships `py.typed`, so mypy and pyright see the full type
information without any extra configuration:

```python
from minifyjs import minify, Result

result: Result = minify("const x = 1;")
assert isinstance(result.code, str)
```

## What's next

- [Python API reference](../python/api.md) — every function
- [Integrations](../integrations/) — using MinifyJS with a framework
- [Platform support](../python/platform-support.md) — which OSes
  have prebuilt wheels