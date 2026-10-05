# Generic Python integration

This page is the starting point for anyone whose framework is not
Flask, Django, or FastAPI. It shows the pattern those three
integrations follow, so you can write your own in about twenty
lines.

## The pattern

Every integration does three things:

1. **Decide where the sources live.** A directory the framework
   already knows about, or one you configure.
2. **Decide where the built output goes.** Usually the same
   directory with a `.min` suffix, sometimes a separate `dist/`.
3. **Run MinifyJS once.** At app startup, during a build step, or
   from a CLI command.

The `Adapter` base class covers everything else. See
[../python/adapter.md](../python/adapter.md) for the full API.

## The one-line case

For a project with a directory of JavaScript files and no special
requirements:

```python
from minifyjs import Adapter

Adapter().minify_directory("static/js")
```

That walks `static/js/`, minifies every `.js` file that does not
already end in `.min.js`, and writes the result next to the input.
No configuration, no subclass.

## Adding options

The default options are `compress=True, mangle=True`. To change
them, pass an `Options`:

```python
from minifyjs import Adapter, Options

adapter = Adapter(options=Options(
    compress=True,
    mangle=True,
    target="es2015",
    sourcemap="external",
))
adapter.minify_directory("static/js")
```

## Changing where output goes

Override `output_name`:

```python
from pathlib import Path
from minifyjs import Adapter

class DistAdapter(Adapter):
    def __init__(self, source: Path, dist: Path) -> None:
        super().__init__()
        self.source = Path(source)
        self.dist = Path(dist)

    def output_name(self, path: Path) -> Path:
        rel = Path(path).relative_to(self.source)
        return self.dist / rel


DistAdapter(source=Path("src"), dist=Path("dist")).minify_directory("src")
```

Now `src/app.js` produces `dist/app.js`.

## Changing what gets processed

Override `extension` to process a different file type, or
`should_process` to filter:

```python
from pathlib import Path
from minifyjs import Adapter

class MjsAdapter(Adapter):
    extension = ".mjs"

    def should_process(self, path: Path) -> bool:
        # Skip test files.
        if path.name.endswith(".test.mjs"):
            return False
        return super().should_process(path)
```

## Build scripts

Wrap the adapter in a script that produces a summary:

```python
#!/usr/bin/env python3
"""Build step for a project with static/js."""

from __future__ import annotations

import sys
from pathlib import Path

from minifyjs import Adapter, MinifyJSError, Options


SOURCE = Path(__file__).parent / "static" / "js"


def main() -> int:
    adapter = Adapter(options=Options(
        compress=True,
        mangle=True,
        target="es2015",
        sourcemap="external",
    ))

    try:
        results = adapter.minify_directory(SOURCE)
    except MinifyJSError as e:
        print(f"build failed: {e}", file=sys.stderr)
        return 1

    total_in = sum(r.original_bytes for r in results)
    total_out = sum(r.minified_bytes for r in results)
    saved = total_in - total_out

    print(f"minified {len(results)} file(s)")
    print(f"  before: {total_in} bytes")
    print(f"  after:  {total_out} bytes")
    print(f"  saved:  {saved} bytes ({saved * 100 / total_in:.1f}%)")

    if total_out >= total_in:
        print("error: output is not smaller", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
```

This is the pattern to copy. It is exactly what the framework
integrations do, minus the framework.

## Reading options from a config file

If your project has a `minifyjs.config.json`, thread it through:

```python
from pathlib import Path
from minifyjs import Adapter, Options, find_config, load_config

config_path = find_config()
if config_path is None:
    opts = Options(compress=True, mangle=True)
else:
    cfg = load_config(config_path)
    opts = Options(
        compress=cfg.compress,
        mangle=cfg.mangle,
        target=cfg.target,
        sourcemap=cfg.sourcemap,
    )

adapter = Adapter(options=opts)
adapter.minify_directory("static/js")
```

## Parallel builds

`Adapter.minify_directory` calls MinifyJS serially. For a project
with hundreds of files, parallelize:

```python
import concurrent.futures
from pathlib import Path

from minifyjs import Adapter

adapter = Adapter()
files = [
    p for p in Path("static/js").rglob("*.js")
    if adapter.should_process(p)
]

with concurrent.futures.ThreadPoolExecutor(max_workers=8) as ex:
    results = list(ex.map(adapter.minify_file, files))
```

The engine is thread-safe. `ThreadPoolExecutor` is enough; you do
not need `ProcessPoolExecutor`.

## Incremental builds

Track modification times yourself:

```python
from pathlib import Path
from minifyjs import Adapter

class IncrementalAdapter(Adapter):
    def minify_file(self, path, output=None, options=None):
        path = Path(path)
        out = Path(output) if output is not None else self.output_name(path)
        if out.exists() and path.stat().st_mtime <= out.stat().st_mtime:
            return None
        return super().minify_file(path, output=out, options=options)

adapter = IncrementalAdapter()
for r in adapter.minify_directory("static/js"):
    if r is None:
        continue
    print(f"built {r.minified_bytes} bytes")
```

Callers must handle the `None` return. The rest of the API is
unchanged.

## When to use `bundle()` instead

If your project has modules that import each other, use
`bundle()` rather than `Adapter`:

```python
from minifyjs import bundle

result = bundle(
    entry_points=["src/main.js"],
    outfile="dist/bundle.js",
    format="esm",
    target="es2015",
)
```

`Adapter` minifies files independently. It does not resolve
imports. `bundle()` does.

A common pattern: use `Adapter` for the standalone scripts in
`static/js/`, and `bundle()` for the one main entry point in
`src/main.js`.

## Testing

Test your adapter against a temp directory:

```python
import tempfile
from pathlib import Path

from minifyjs import Adapter

def test_adapter_writes_minified_file():
    with tempfile.TemporaryDirectory() as tmp:
        tmp = Path(tmp)
        (tmp / "app.js").write_text("function add(a, b) { return a + b; }")
        Adapter().minify_directory(tmp)
        out = tmp / "app.min.js"
        assert out.is_file()
        assert "\n" not in out.read_text()
```

## See also

- [../python/adapter.md](../python/adapter.md) — the `Adapter` API
  in full
- [flask.md](flask.md) — the Flask-specific recipe
- [django.md](django.md) — the Django-specific recipe
- [fastapi.md](fastapi.md) — the FastAPI-specific recipe