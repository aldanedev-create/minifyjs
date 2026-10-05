# Python Troubleshooting

## `BinaryNotFoundError`

The Python package needs the native `minifyjs` binary. Install a wheel for
one of the supported platforms, or build the binary from a source checkout:

```console
cd core
go build -o ../python/minifyjs/bin/minifyjs ./cmd/minifyjs
cd ..
python -m pip install -e python
```

On Windows, build `minifyjs.exe`. Check the detected platform with:

```python
from minifyjs._platform import platform_tag
print(platform_tag())
```

## `PermissionError` when starting the binary

The binary must be executable on Unix systems. A wheel normally preserves
this bit; repair a source checkout with:

```console
chmod +x python/minifyjs/bin/minifyjs
```

The wheel builder also preserves mode `0755` for packaged binaries.

## Alpine Linux or musl

The published Linux wheels use the `manylinux2014` glibc baseline and do
not install directly on Alpine/musl. Build from source on Alpine, or use a
platform-specific musllinux wheel when one is published.

## Invalid JavaScript

`minify()` and `optimize()` raise `MinifyError` when the native engine cannot
parse the input. The exception contains the diagnostic location and message:

```python
from minifyjs import MinifyError, minify

try:
    minify("const = 1;")
except MinifyError as exc:
    print(exc)
```

Use `source_name="app.js"` to put a useful filename in diagnostics.

## Source maps

For a transform, `sourcemap="inline"` adds a data URL to `Result.code`.
`sourcemap="external"` returns compact JSON in `Result.map`:

```python
result = minify("const x = 1;", sourcemap="external")
print(result.map)
```

For bundles, external maps are written beside the output file and are also
returned in `Result.map`:

```python
result = bundle(
    ["src/main.js"],
    outfile="dist/bundle.js",
    format="esm",
    sourcemap="external",
)
# dist/bundle.js.map and result.map are available after the call.
```

## Python environments

Use a virtual environment for development and release checks:

```console
python -m venv .venv
. .venv/bin/activate
python -m pip install -e 'python[dev]'
```

Python 3.14 is included in the package metadata and type-check target. Run
the full test suite with a Python 3.14 interpreter when it is available.

## Bundle failures

`bundle()` requires exactly one of `outdir` or `outfile`. A missing entry
file raises `BundleError`. Code splitting requires `format="esm"` and an
output directory.
