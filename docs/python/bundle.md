# Python: `bundle()`

Multi-file project support. Resolves imports, tree-shakes unused
exports, and writes one or more output files.

## Signature

```python
def bundle(
    entry_points: List[str],
    *,
    outdir: Optional[str] = None,
    outfile: Optional[str] = None,
    working_dir: Optional[str] = None,
    platform: str = "browser",
    splitting: bool = False,
    compress: bool = True,
    mangle: bool = True,
    target: Optional[str] = None,
    format: Optional[str] = "esm",
    sourcemap: Optional[str] = None,
    banner: Optional[str] = None,
    footer: Optional[str] = None,
    legal_comments: Optional[str] = None,
) -> Result
```

Exactly one of `outdir` or `outfile` must be given. Passing neither
raises `ValueError`. Passing both raises `ValueError`.

## Minimal call

```python
from minifyjs import bundle

result = bundle(
    ["src/main.js"],
    outfile="dist/bundle.js",
    format="esm",
)
```

This resolves `src/main.js`, follows its imports, drops unused
exports, and writes `dist/bundle.js`.

## Difference from `minify()` and `optimize()`

| | `minify` / `optimize` | `bundle` |
|---|---|---|
| Input | A string | File paths |
| Output | A string in `result.code` | Files on disk |
| Imports | Not resolved | Resolved and inlined |
| Tree-shaking | No | Yes |
| Network/filesystem | No | Reads source files |

`bundle()` is the only Python API function that touches the
filesystem.

## Return value

`bundle()` returns a `Result`, but `result.code` is usually empty.
The output went to disk, not to the result.

`result.diagnostics` is populated. Check it for errors:

```python
if result.has_errors:
    for d in result.diagnostics:
        print(f"error: {d}")
    raise SystemExit(1)
```

## Options

### `entry_points`

Required. A list of file paths to start bundling from.

```python
result = bundle(
    ["src/main.js", "src/admin.js"],
    outdir="dist",
    format="esm",
)
```

Multiple entry points produce multiple output files.

### `outdir` vs `outfile`

| | When to use |
|---|---|
| `outfile` | One entry point, one output file. |
| `outdir` | Multiple entry points, or code splitting. |

```python
# One file
bundle(["src/main.js"], outfile="dist/bundle.js")

# Multiple files
bundle(["src/a.js", "src/b.js"], outdir="dist", splitting=True)
```

### `format`

Controls the module wrapper.

| Value | What happens |
|---|---|
| `"esm"` (default) | ES module syntax preserved |
| `"cjs"` | CommonJS output (`module.exports`, `require`) |
| `"iife"` | Immediately-invoked function expression |

### `platform`

| Value | Assumes |
|---|---|
| `"browser"` (default) | Browser globals; warns on Node builtins |
| `"node"` | Node globals; `node:` builtins are external |
| `"neutral"` | Neither; nothing is implicit |

### `splitting`

```python
bundle(
    ["src/a.js", "src/b.js"],
    outdir="dist",
    format="esm",
    splitting=True,
)
```

Hoists shared imports into common chunks. Requires `format="esm"`.

### `compress` and `mangle`

Both default to `True` in `bundle()` (unlike `minify()`, where they
default to `False`). Bundles are almost always production artifacts,
so aggressive optimization is the sensible default.

To disable:

```python
bundle(["src/main.js"], outfile="dist/bundle.js", compress=False)
```

### `sourcemap`

Same as `minify()`. In bundle mode, one map is written per output
file.

```python
bundle(["src/main.js"], outfile="dist/bundle.js", sourcemap="external")
# Writes dist/bundle.js and dist/bundle.js.map
```

### `banner` and `footer`

Added verbatim to every output file.

```python
bundle(
    ["bin/cli.js"],
    outfile="dist/cli.js",
    platform="node",
    format="cjs",
    banner="#!/usr/bin/env node",
)
```

The banner is useful for shebangs, which would otherwise be
stripped by the parser.

## Errors

`bundle()` raises `BundleError` (a subclass of `MinifyError`) when
the build fails:

```python
from minifyjs import bundle, BundleError

try:
    bundle(["src/main.js"], outfile="dist/bundle.js")
except BundleError as e:
    print(f"build failed: {e}")
    raise SystemExit(1)
```

Common causes:

- A file in the import graph does not exist.
- A file has a syntax error.
- A module tries to import a Node built-in when `platform="browser"`.
- `outdir` and `outfile` were both passed.
- Neither `outdir` nor `outfile` was passed.

## Working directory

`working_dir` is the directory all relative paths in `entry_points`
are resolved against.

```python
bundle(
    ["main.js"],
    outfile="dist/bundle.js",
    working_dir="/path/to/project",
)
```

Default is the current working directory.

## Tree shaking

Unused exports are removed automatically. To keep them, either
import them somewhere, or mark the module as having side effects in
`package.json`:

```json
{
  "name": "my-package",
  "version": "1.0.0",
  "sideEffects": ["*.css", "./src/polyfills.js"]
}
```

The `sideEffects` field is read from the package that contains the
module being tree-shaken. Modules with side effects are kept even
if none of their exports are used.

## Full example

```python
#!/usr/bin/env python3
"""Build script for a small project."""

from pathlib import Path
from minifyjs import bundle, BundleError

PROJECT = Path(__file__).parent
DIST = PROJECT / "dist"


def main() -> int:
    DIST.mkdir(exist_ok=True)

    try:
        result = bundle(
            entry_points=[str(PROJECT / "src/main.js")],
            outfile=str(DIST / "bundle.js"),
            format="esm",
            target="es2015",
            compress=True,
            mangle=True,
            sourcemap="external",
            banner="/* built by CI */",
        )
    except BundleError as e:
        print(f"build failed: {e}")
        return 1

    if result.has_errors:
        for d in result.diagnostics:
            print(f"error: {d}")
        return 1

    size = (DIST / "bundle.js").stat().st_size
    print(f"wrote {DIST / 'bundle.js'} ({size} bytes)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
```

## When to use esbuild directly instead

`bundle()` does not support:

- CSS, image, font, or any non-JS loader
- esbuild plugins
- A `tsconfig.json` reference
- Custom `--loader:` overrides
- Environment variable substitution beyond `--define`

If you need any of those, use esbuild's own CLI. `bundle()` exists
for the common case: a JavaScript or TypeScript project with no
special asset handling.

## See also

- [optimize.md](optimize.md) — single-file optimization
- [../cli/bundle.md](../cli/bundle.md) — the CLI equivalent
## Production build controls

Bundle calls honor `mangle=True` and use esbuild's standard tree shaking and dependency resolution. Packages and dynamic imports are bundled unless explicitly externalized.

```python
result = bundle(
    ["src/main.js"], working_dir="/absolute/project/path", outdir="dist",
    format="esm", splitting=True, compress=True, mangle=True,
    target="es2020", entry_names="[name]-[hash]",
    chunk_names="chunks/[name]-[hash]", asset_names="assets/[name]-[hash]",
    metafile="reports/build.json", define={"DEBUG": "false"},
    drop=["debugger"], charset="utf8",
)
for output in result.output_files:
    print(output["path"], output["bytes"])
print(result.minified_bytes)
print(result.metafile["inputs"])
```

Additional options: `external=["dependency"]`, `packages="bundle"` (default) or `"external"`, `tree_shaking=None` (esbuild default), `True` or `False`, and `charset="utf8"` or `"ascii"`. Relative entries, outputs and metafiles resolve against `working_dir`, defaulting to the current directory. Python bundle calls ignore discovered CLI configuration files to keep API settings reproducible.

`result.code` contains the JavaScript for `outfile` builds and is empty for `outdir`. `output_files` and parsed `metafile` are returned in both modes. `minified_bytes` counts all emitted files, including maps. `original_bytes` remains zero because bundles include dependencies beyond their entries.

Migration: ESM no longer automatically preserves package imports; re-exports no longer disable tree shaking; dynamic imports are bundled without splitting. Use explicit `external` patterns or `packages="external"` where preservation is intended, and retest application behavior before deployment.
