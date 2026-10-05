# Python API

The `minifyjs` package exposes three functions and a small set of
types. This page is the map; the other pages in this directory are
the detail.

## Install

```console
pip install minifyjs
```

## The public surface

```python
from minifyjs import (
    # The three entry points
    minify,
    optimize,
    bundle,

    # Types
    Options,
    BundleOptions,
    Result,
    Diagnostic,
    Adapter,
    Config,

    # Config helpers
    find_config,
    load_config,

    # Exceptions
    MinifyJSError,
    MinifyError,
    BundleError,
    BinaryNotFoundError,
)
```

Every name above is part of the public API and follows semantic
versioning. Anything not in this list is internal and may change
without notice.

## Quick reference

| Function | Signature | Returns |
|---|---|---|
| `minify` | `minify(source, **opts) -> Result` | Whitespace-only minification |
| `optimize` | `optimize(source, **opts) -> Result` | Full optimization + mangling |
| `bundle` | `bundle(entries, **opts) -> Result` | Multi-file build, writes to disk |

| Type | Purpose |
|---|---|
| `Options` | Keyword arguments for `minify` and `optimize` |
| `BundleOptions` | Keyword arguments for `bundle` |
| `Result` | The value all three functions return |
| `Diagnostic` | A single message produced by the engine |
| `Adapter` | Base class for framework integrations |
| `Config` | A parsed `minifyjs.config.json` |

## Minimal example

```python
from minifyjs import minify

result = minify("function add(a, b) { return a + b; }")
print(result.code)
# function add(a,b){return a+b;}
```

## The Result object

Every function returns a `Result`:

```python
result.code              # the minified JavaScript as a string
result.map               # source map JSON, or ""
result.original_bytes    # byte length of the input
result.minified_bytes    # byte length of the output
result.diagnostics       # list of Diagnostic objects

result.ratio             # minified_bytes / original_bytes
result.bytes_saved       # original_bytes - minified_bytes
result.has_errors        # True if any diagnostic is an error
result.has_warnings      # True if any diagnostic is a warning
```

## Type checking

The package ships `py.typed`, so mypy and pyright see full type
information without any extra configuration:

```python
from minifyjs import minify, Result

result: Result = minify("const x = 1;")
reveal_type(result.code)  # str
```

## Synchronous only

The Python API is synchronous. There is no `async` variant. A
minify call takes a few milliseconds and is CPU-bound on the
subprocess side, so an async API would add complexity without
adding throughput.

To run many calls in parallel, use `concurrent.futures`:

```python
import concurrent.futures
from minifyjs import optimize

with concurrent.futures.ThreadPoolExecutor(max_workers=8) as ex:
    futures = [ex.submit(optimize, src) for src in sources]
    results = [f.result() for f in futures]
```

The engine is thread-safe. See `core/integration/api/concurrency_test.go`
for the tests that prove it.

## Where to go next

- [minify.md](minify.md) — the `minify()` function in detail
- [optimize.md](optimize.md) — the `optimize()` function
- [bundle.md](bundle.md) — the `bundle()` function
- [configuration.md](configuration.md) — reading config files
- [errors.md](errors.md) — the exception hierarchy
- [platform-support.md](platform-support.md) — which OSes have wheels