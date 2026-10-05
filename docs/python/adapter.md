# Python: the `Adapter` class

`Adapter` is the framework-agnostic base class MinifyJS ships for
build-system integrations. It knows how to walk a directory, minify
the JavaScript files it finds, and write the results back out.
Anything more specific — knowing where Django's `STATIC_ROOT` is,
or how a Flask app's `static_folder` maps to a URL — belongs in a
subclass, not in MinifyJS.

This page is the guide for using `Adapter` in real projects. It
starts with the one-line case, moves through directory builds, and
finishes with writing your own adapter for a framework MinifyJS has
never heard of.

## Why an adapter exists

The natural alternative would be to ship `MinifyJS-for-Django`,
`MinifyJS-for-Flask`, `MinifyJS-for-FastAPI`, and so on. That
approach has two problems:

1. **Frameworks change.** A Django adapter written today will break
   in a Django release two years from now. MinifyJS would have to
   track that release.
2. **The list is infinite.** Django, Flask, FastAPI, Pyramid,
   Bottle, Tornado, aiohttp, Starlette, Litestar, a custom build
   system, a CI pipeline, and a hundred others.

`Adapter` flips the dependency. MinifyJS provides a small, stable
base class. Each framework's integration is written by the people
who use that framework, in that framework's own repository, in
about twenty lines. MinifyJS never learns that Django exists.

## The three method calls that matter

`Adapter` exposes three things that a subclass will actually use:

| Method | What it does |
|---|---|
| `minify_file(path, output=None)` | Minify one file, write the result, return a `Result` |
| `minify_directory(dir, recursive=True)` | Minify every matching file under a directory, return a list of `Result` |
| `output_name(path)` | Return the output path for a given input path |

Plus two hooks that a subclass overrides:

| Hook | Default | When to override |
|---|---|---|
| `extension` | `".js"` | To process a different file extension |
| `output_suffix` | `".min"` | To change the `.min.js` convention |
| `should_process(path)` | True for `.js` files that are not already minified | To skip specific files |

That is the entire integration surface. Everything else is a
subclass's business.

## The one-line case

For a project that just wants "minify every `.js` under `static/`
into `static/*.min.js`":

```python
from minifyjs import Adapter

Adapter().minify_directory("static")
```

That is the whole program. It walks `static/`, minifies every
`.js` file, and writes `name.min.js` next to each one. Files that
already end in `.min.js` are skipped.

The default `Adapter` handles the following:

- Skips `node_modules/`, `.git/`, and `__pycache__/` when recursing.
- Skips files whose names already end in `.min.js`.
- Skips non-`.js` files.
- Uses `compress=True, mangle=True` as the default options.
- Returns a list of `Result` objects, one per file processed.

## Custom options

The default options are `Options(compress=True, mangle=True)`. To
change them, pass an `Options` instance:

```python
from minifyjs import Adapter, Options

adapter = Adapter(options=Options(
    compress=True,
    mangle=False,           # keep names readable for debugging
    target="es2015",
    sourcemap="external",
    banner="/* (c) 2026 */",
))
adapter.minify_directory("static/js")
```

`Options` is the same dataclass used by `minify()` and `optimize()`.
See [../python/minify.md](minify.md) for the full field list.

## Where the output goes

By default, output is written next to the input with a `.min`
inserted before the extension:

```
static/js/app.js       →  static/js/app.min.js
static/js/utils.js     →  static/js/utils.min.js
static/js/lib/api.js   →  static/js/lib/api.min.js
```

To change where files go, override `output_name`:

```python
from pathlib import Path
from minifyjs import Adapter

class DistAdapter(Adapter):
    def __init__(self, source: Path, dist: Path) -> None:
        super().__init__()
        self.source = source
        self.dist = dist

    def output_name(self, path: Path) -> Path:
        rel = path.relative_to(self.source)
        return self.dist / rel

adapter = DistAdapter(source=Path("src"), dist=Path("dist"))
adapter.minify_directory("src")
```

Now `src/app.js` produces `dist/app.js` (same name, different
directory).

## Where the input is walked from

`minify_directory(root)` walks everything under `root`. When `root`
is a specific directory like `static/js/`, only that directory is
walked. When `root` is a project directory like `.`, everything
under it is walked, but the default `Adapter`'s noise-directory
filter keeps `node_modules/`, `.git/`, and `__pycache__/` out of the
walk.

To walk a specific set of directories, call `minify_directory` once
per directory:

```python
adapter = Adapter()
for d in ["static/js", "static/vendor", "src/scripts"]:
    adapter.minify_directory(d)
```

## Big projects

For a project with many files, or with a build system that already
knows what needs rebuilding, the useful patterns are:

### Pattern 1 — build script with a summary

```python
#!/usr/bin/env python3
"""Build step for a project with a static/js directory."""

from __future__ import annotations

import sys
from pathlib import Path

from minifyjs import Adapter, MinifyJSError, Options


PROJECT = Path(__file__).parent
SOURCE = PROJECT / "static" / "js"
DIST = PROJECT / "static" / "dist"


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
    pct = (saved / total_in * 100) if total_in else 0

    print(f"minified {len(results)} file(s)")
    print(f"  before: {total_in} bytes")
    print(f"  after:  {total_out} bytes")
    print(f"  saved:  {saved} bytes ({pct:.1f}%)")

    # Fail the build if the output grew. This catches the case
    # where a `--banner` or a very small file makes the output
    # larger, which would be a silent regression in a CI script.
    if total_out >= total_in:
        print("error: minified output is not smaller", file=sys.stderr)
        return 1

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
```

Run it as `python build.py` or wire it into a Makefile / CI step.

### Pattern 2 — incremental build

`Adapter` does not track modification times. If you only want to
re-minify files that changed since the last run, wrap it:

```python
import time
from pathlib import Path
from minifyjs import Adapter

CACHE = Path(".build-cache")


def is_stale(source: Path, output: Path) -> bool:
    if not output.exists():
        return True
    return source.stat().st_mtime > output.stat().st_mtime


class IncrementalAdapter(Adapter):
    def minify_file(self, path, output=None, options=None):
        path = Path(path)
        out = Path(output) if output is not None else self.output_name(path)
        if not is_stale(path, out):
            return None  # signal "skipped"
        return super().minify_file(path, output=out, options=options)


# Callers check for None:
results = [r for r in IncrementalAdapter().minify_directory("static") if r is not None]
```

The `None` return is a signal, not an error. Existing callers that
expect a `Result` will not work with this subclass; that is the
tradeoff for the speedup.

### Pattern 3 — build script that also bundles

For a project that has a mix of independent scripts and one main
entry point, run `Adapter` for the scripts and `bundle()` for the
entry point:

```python
from pathlib import Path

from minifyjs import Adapter, Options, bundle, BundleError


def build() -> int:
    # 1. Minify every standalone script in place.
    adapter = Adapter(options=Options(compress=True, mangle=True))
    adapter.minify_directory("static/js")

    # 2. Bundle the main entry point into dist/.
    try:
        result = bundle(
            entry_points=["src/main.js"],
            outfile="dist/bundle.js",
            format="esm",
            target="es2015",
            sourcemap="external",
        )
    except BundleError as e:
        print(f"bundle failed: {e}")
        return 1

    return 0 if not result.has_errors else 1


if __name__ == "__main__":
    raise SystemExit(build())
```

### Pattern 4 — build a `dist/` directory from scratch

For a project that wants a clean `dist/` on every build, subclass
`Adapter` to control where output lands:

```python
import shutil
from pathlib import Path

from minifyjs import Adapter


class DistAdapter(Adapter):
    """Writes minified JS to dist/, copying everything else verbatim."""

    def __init__(self, source: Path, dist: Path) -> None:
        super().__init__()
        self.source = Path(source)
        self.dist = Path(dist)

    def output_name(self, path: Path) -> Path:
        rel = Path(path).relative_to(self.source)
        return self.dist / rel

    def build(self) -> int:
        if self.dist.exists():
            shutil.rmtree(self.dist)
        self.dist.mkdir(parents=True)

        count = 0
        for src in self._iter_files(self.source, recursive=True):
            dest = self.output_name(src)
            dest.parent.mkdir(parents=True, exist_ok=True)
            if self.should_process(src):
                self.minify_file(src, output=dest)
                count += 1
            else:
                dest.write_bytes(src.read_bytes())
        return count


n = DistAdapter(source=Path("src"), dist=Path("dist")).build()
print(f"built {n} JS file(s) into dist/")
```

## Writing your own adapter

Every adapter subclasses `Adapter` and overrides at most a few
methods. Here is the checklist.

### Step 1 — Decide what to process

The default processes `.js` files that do not already end in
`.min.js`. If your project uses a different convention, override
`extension` or `should_process`:

```python
class MjsAdapter(Adapter):
    extension = ".mjs"

    def should_process(self, path):
        return super().should_process(path) and not path.name.startswith("_")
```

### Step 2 — Decide where output goes

The default writes next to the input. If your framework expects
build artifacts somewhere else, override `output_name`:

```python
class DjangoAdapter(Adapter):
    def __init__(self, static_root, static_dirs):
        super().__init__()
        self.static_root = Path(static_root)
        self.static_dirs = [Path(p) for p in static_dirs]

    def output_name(self, path):
        path = Path(path)
        for src in self.static_dirs:
            try:
                rel = path.resolve().relative_to(src.resolve())
            except ValueError:
                continue
            return self.static_root / rel.parent / f"{rel.stem}.min{rel.suffix}"
        return self.static_root / f"{path.stem}.min{path.suffix}"
```

### Step 3 — Decide what options to use

Pass an `Options` instance to `Adapter.__init__`:

```python
class MyAdapter(Adapter):
    def __init__(self, **kwargs):
        super().__init__(options=Options(
            compress=True,
            mangle=True,
            target="es2015",
            **kwargs,
        ))
```

Callers can override per-call by passing `options=` to `minify_file`
or `minify_directory`.

### Step 4 — Optionally add framework-specific methods

The base class exposes `minify_file` and `minify_directory`.
Framework adapters usually add a `build()` or `run()` that returns a
summary the framework's CLI can print:

```python
class MyAdapter(Adapter):
    def build(self) -> int:
        results = self.minify_directory(self.source_dir)
        print(f"minified {len(results)} files")
        return 0
```

## A complete worked example

Here is an adapter for a hypothetical build system called "Brisk"
that expects JavaScript sources in `brisk/src/` and build output in
`brisk/build/`, with `.brisk.js` as the output suffix.

```python
from __future__ import annotations

from pathlib import Path

from minifyjs import Adapter, Options, Result


class BriskAdapter(Adapter):
    """Adapter for the Brisk build system.

    Brisk expects:
      - sources in <project>/brisk/src/**/*.js
      - output in <project>/brisk/build/**, same relative path
      - output files named <name>.brisk.js
    """

    output_suffix = ".brisk"

    def __init__(self, project: Path, options: Options | None = None) -> None:
        super().__init__(options=options or Options(
            compress=True,
            mangle=True,
            target="es2015",
        ))
        self.project = Path(project)
        self.source_dir = self.project / "brisk" / "src"
        self.build_dir = self.project / "brisk" / "build"

    def output_name(self, path: Path) -> Path:
        rel = Path(path).relative_to(self.source_dir)
        return self.build_dir / rel.parent / f"{rel.stem}{self.output_suffix}{rel.suffix}"

    def build(self) -> int:
        """Build every source file. Returns an exit code."""
        if not self.source_dir.is_dir():
            print(f"error: {self.source_dir} does not exist")
            return 1

        self.build_dir.mkdir(parents=True, exist_ok=True)

        results: list[Result] = []
        for source in self._iter_files(self.source_dir, recursive=True):
            if not self.should_process(source):
                continue
            dest = self.output_name(source)
            dest.parent.mkdir(parents=True, exist_ok=True)
            results.append(self.minify_file(source, output=dest))

        total_in = sum(r.original_bytes for r in results)
        total_out = sum(r.minified_bytes for r in results)
        print(f"brisk: built {len(results)} file(s) "
              f"({total_in} -> {total_out} bytes)")
        return 0


if __name__ == "__main__":
    import sys
    project = Path(sys.argv[1]) if len(sys.argv) > 1 else Path(".")
    raise SystemExit(BriskAdapter(project).build())
```

Things worth noticing:

1. **The subclass overrides `output_name`, not `minify_file`.**
   Changing where output goes does not require reimplementing the
   minification call.
2. **The `build()` method is a convenience, not a contract.** The
   base class does not call it. A framework's CLI would.
3. **The subclass uses `self._iter_files`, `self.should_process`,
   and `self.minify_file`.** All three are inherited. The subclass
   only adds the Brisk-specific knowledge (where sources live, what
   the output is called).
4. **The options default is on the subclass, not the base.**
   Different frameworks want different defaults; the base class
   just provides the mechanism.

## What `Adapter` does not do

- **Modification-time tracking.** Use
  `path.stat().st_mtime` yourself, or a real build system.
- **Watch mode.** `Adapter` runs once per call. The CLI has
  `--watch`; the Python adapter does not.
- **Parallelism.** `Adapter.minify_directory` calls the engine
  serially. For parallel builds, use
  `concurrent.futures.ThreadPoolExecutor` on top:

  ```python
  import concurrent.futures
  from pathlib import Path
  from minifyjs import Adapter

  adapter = Adapter()
  files = list(adapter._iter_files(Path("static"), recursive=True))

  with concurrent.futures.ThreadPoolExecutor(max_workers=8) as ex:
      results = list(ex.map(adapter.minify_file, files))
  ```

- **Dependency awareness.** If file A imports file B, `Adapter`
  does not know or care. That is what `bundle()` is for.
- **CSS, images, fonts, or any non-JS asset.** `Adapter` is a
  JavaScript minifier. `DistAdapter`-style subclasses usually copy
  non-JS files verbatim, but that is a subclass's decision.

## When not to use `Adapter`

`Adapter` is for the common case of "there is a folder of `.js`
files, minify them". It is not for every build integration. Use
`bundle()` instead when:

- Files import each other.
- There is one entry point and many modules.
- Tree shaking matters.
- Code splitting matters.

Use the CLI directly when:

- The build is a shell script.
- There is exactly one file.
- The output path is fixed and does not need decision logic.

Use neither when:

- The project uses webpack, Rollup, or Vite for bundling. Those
  tools already have minification built in. Running MinifyJS on top
  of them produces worse results, not better.
- The project needs a plugin system. Use esbuild directly.

## See also

- [minify.md](minify.md) — the underlying single-file function
- [bundle.md](bundle.md) — for multi-file builds
- [configuration.md](configuration.md) — reading options from a
  config file
- [../integrations/generic-python.md](../integrations/generic-python.md) —
  short recipes for framework adapters