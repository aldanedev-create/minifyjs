# MinifyJS



 
 <p align="center">
  <img src="https://raw.githubusercontent.com/aldanedev-create/minifyjs/main/assets/Minifyjs-log.png" alt="minifyjs Logo"
   width="200"/>
</p>

A native JavaScript minifier and optimizer, driven from Python or the
command line. No Node.js, no npm, no JavaScript runtime — the engine
is a single self-contained binary.

```console
$ pip install minifyjs
$ minifyjs app.js -o app.min.js
```

```python
from minifyjs import minify

result = minify("function add(a, b) { return a + b; }")
print(result.code)
# function add(a,b){return a+b;}
```

## Why

JavaScript tooling usually means installing Node.js and a package
manager before you can minify a file. MinifyJS ships as a native
executable inside a platform-specific Python wheel: one `pip install`
and you have a working minifier, with no JavaScript ecosystem on the
machine.

## How it works

MinifyJS is a thin, opinionated wrapper around
[esbuild](https://esbuild.github.io/)'s Go API. esbuild does the
parsing, optimization, mangling, and code generation — the same
engine that powers production builds at scale — and MinifyJS adds:

- A stable Python API with typed options and structured diagnostics.
- A native CLI with sane defaults and predictable exit codes.
- A single-binary distribution model (no Node.js, no toolchain).
- Framework-agnostic integration points for build systems that need
  minification as a step rather than as a Node tool.

The Go source of esbuild is compiled into the MinifyJS binary.
Nothing shells out to Node, and nothing downloads a JavaScript
runtime at install time.

## Install

### Python

```console
pip install minifyjs
```

Requires Python 3.8+. Wheels are published for Linux (x86-64,
aarch64), macOS (x86-64, arm64), and Windows (x86-64).

### CLI only

Download the appropriate binary from the [releases page](../../releases)
and put it on your `PATH`. The Python package is not required to use
the CLI, but it is the easiest way to get the binary.

## Usage

### Command line

```console
minifyjs input.js -o output.js
minifyjs input.js --compress --mangle
cat input.js | minifyjs > output.js
minifyjs --version
minifyjs --help
```

### Python

```python
from minifyjs import minify, optimize, Options

# Whitespace/comments only
result = minify(source)

# Full optimization + identifier mangling
result = optimize(source)

# Explicit control
result = minify(source, Options(compress=True, mangle=True))
print(result.code)
print(f"{result.original_bytes} -> {result.minified_bytes} bytes")
for d in result.diagnostics:
    print(d)
```

See [`docs/python/`](docs/python/) for the full API reference.

## Documentation

- [Getting started](docs/getting-started/)
- [CLI reference](docs/cli/)
- [Python API](docs/python/)
- [Engine and esbuild integration](docs/engine/)
- [Architecture](docs/architecture.md)
- [Contributing](CONTRIBUTING.md)
- [Security policy](SECURITY.md)
- [Roadmap](ROADMAP.md)

## Status

Pre-1.0. The API is stable enough to depend on; the on-disk format
of the bundled binary may still change between minor versions.

## License

MIT. See [LICENSE](LICENSE).

## Acknowledgements

MinifyJS would not exist without
[esbuild](https://github.com/evanw/esbuild), which provides the
JavaScript parsing and optimization engine. esbuild is MIT-licensed
and is compiled into MinifyJS's native binary.
