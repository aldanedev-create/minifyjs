# minifyjs

Python bindings for [MinifyJS](https://github.com/aldanedev-create/minifyjs),
a native JavaScript minifier and optimizer. No Node.js required.

## Install

```console
pip install minifyjs
```

## Use

```python
from minifyjs import minify

result = minify("function add(a, b) { return a + b; }")
print(result.code)
# function add(a,b){return a+b;}
```

Full optimization and identifier mangling:

```python
from minifyjs import optimize

result = optimize("const x = 1 + 2 + 3;")
print(result.code)   # const x=6;
```

Explicit control:

```python
from minifyjs import minify, Options

result = minify(source, Options(
    compress=True,
    mangle=True,
    target="es2015",
    sourcemap="inline",
))
```

Bundling:

```python
from minifyjs import bundle

result = bundle("src/main.js", outdir="dist", format="esm")
```

## Adapter for frameworks

MinifyJS ships a generic adapter (`minifyjs.adapter`) that any
framework or build system can subclass. See
[`docs/integrations/generic-python.md`](https://github.com/aldanedev-create/minifyjs/blob/main/docs/integrations/generic-python.md)
for how to write an adapter for Django, Flask, FastAPI, or a custom
build pipeline.

## Command line

The native `minifyjs` binary is installed alongside the Python
package:

```console
minifyjs app.js -o app.min.js
minifyjs app.js --compress --mangle
cat app.js | minifyjs
```

Or, without the binary on `PATH`:

```console
python -m minifyjs --version
```

## Documentation

See the main repository's [`docs/python/`](https://github.com/aldanedev-create/minifyjs/tree/main/docs/python).

## License

MIT.