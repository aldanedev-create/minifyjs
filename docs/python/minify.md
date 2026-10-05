# Python: `minify()`

The simplest entry point. Removes whitespace and comments. Does not
rename identifiers or fold constants.

## Signature

```python
def minify(
    source: str,
    *,
    compress: bool = False,
    mangle: bool = False,
    target: Optional[str] = None,
    format: Optional[str] = None,
    sourcemap: Optional[str] = None,
    banner: Optional[str] = None,
    footer: Optional[str] = None,
    legal_comments: Optional[str] = None,
    source_name: Optional[str] = None,
) -> Result
```

Every argument is keyword-only. The source is positional.

## The minimal call

```python
from minifyjs import minify

result = minify("function add(a, b) { return a + b; }")
print(result.code)
# function add(a,b){return a+b;}
```

With no arguments, `minify()` removes whitespace and comments and
nothing else.

## Enabling compression

```python
result = minify(source, compress=True)
```

`compress=True` turns on constant folding, dead-code elimination,
and expression simplification.

## Enabling mangling

```python
result = minify(source, mangle=True)
```

`mangle=True` renames local variables and function parameters.
Names that are observable from outside the file (top-level
declarations, exports) are preserved.

## Both at once

```python
result = minify(source, compress=True, mangle=True)
```

This is equivalent to calling `optimize(source)`. See
[optimize.md](optimize.md).

## Picking a target

```python
result = minify(source, target="es2015")
```

Valid values are the same as the CLI: `es5`, `es2015` through
`es2024`, `esnext`, browser+version pairs (`chrome90`), Node
versions (`node20`), and comma-separated compounds.

The default is `None`, which means "do not lower anything".

## Source maps

```python
result = minify(source, sourcemap="inline")
print("sourceMappingURL=data:application/json" in result.code)
# True

result = minify(source, sourcemap="external")
print(result.map)
# {"version":3,"sources":["<stdin>"],"mappings":"..."}
```

| Value | Result |
|---|---|
| `None` (default) | No map |
| `"inline"` | Map is embedded in `result.code`; `result.map` is empty |
| `"external"` | Map is in `result.map`; `result.code` ends with a `sourceMappingURL` comment |
| `"both"` | Map is both inline and in `result.map` |

## Content injection

```python
result = minify(
    source,
    banner="/* (c) 2026 Acme Corp */",
    footer="// built by CI",
)
```

Both are added verbatim. They are not parsed, minified, or
otherwise touched.

## Legal comments

```python
result = minify(source, legal_comments="inline")
```

| Value | Behavior |
|---|---|
| `None` (default) | esbuild's default, which is `"eof"` |
| `"none"` | Remove `/*! ... */` comments |
| `"inline"` | Keep them where they are |
| `"eof"` | Move them to the end of the file |
| `"external"` | Write them to a separate file (bundle mode only) |

## Source name in diagnostics

```python
result = minify(source, source_name="src/app.js")
```

Sets the filename shown in diagnostics. Useful when the source was
read from a file and the caller wants error messages to reference
the file:

```python
from pathlib import Path
from minifyjs import minify, MinifyError

path = Path("src/app.js")
try:
    result = minify(path.read_text(), source_name=str(path))
except MinifyError as e:
    print(e)   # src/app.js:12:5: error: ...
```

The default is `None`, which produces `<stdin>` in diagnostics.

## What the result contains

```python
result.code              # minified JavaScript
result.map               # source map JSON, or ""
result.original_bytes    # len(source.encode('utf-8'))
result.minified_bytes    # len(result.code.encode('utf-8'))
result.diagnostics       # list[Diagnostic]

result.ratio             # minified_bytes / original_bytes
result.bytes_saved       # original_bytes - minified_bytes
result.has_errors        # bool
result.has_warnings      # bool
```

## Handling errors

`minify()` raises `MinifyError` when the input cannot be parsed:

```python
from minifyjs import minify, MinifyError

try:
    result = minify("function () { } }")
except MinifyError as e:
    print(f"failed: {e}")
```

Warnings do not raise. They appear on `result.diagnostics`:

```python
for d in result.diagnostics:
    if d.severity == "warning":
        print(f"warning at {d.file}:{d.line}:{d.column}: {d.message}")
```

## Idempotence

Minifying an already-minified file is a no-op:

```python
r1 = minify(source)
r2 = minify(r1.code)
assert r1.code == r2.code
```

This is guaranteed by the tests in
`python/tests/test_idempotence.py`. If it ever fails, it is a bug.

## Concurrency

`minify()` is thread-safe. Call it from multiple threads or
processes without synchronization:

```python
import concurrent.futures
from minifyjs import minify

with concurrent.futures.ThreadPoolExecutor(max_workers=8) as ex:
    futures = [ex.submit(minify, s) for s in sources]
    results = [f.result() for f in futures]
```

## What `minify()` does not do

- **File I/O.** It takes a string and returns a `Result`. Reading
  and writing files is the caller's job.
- **Bundling.** For multi-file projects, use `bundle()`.
- **Type checking.** It strips TypeScript types. It does not verify
  them. See [../supported-syntax.md](../supported-syntax.md).

## See also

- [optimize.md](optimize.md) — `minify()` with `compress=True,
  mangle=True`
- [bundle.md](bundle.md) — multi-file builds
- [errors.md](errors.md) — the exception hierarchy