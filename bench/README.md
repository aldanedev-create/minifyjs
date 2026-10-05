# MinifyJS benchmarks

This directory measures MinifyJS against the other JavaScript
minifiers users are likely to compare it with. It exists so the
project can make claims about performance that are **reproducible**
rather than asserted.

## What is measured

Four axes:

| Axis | What it means | Script |
|---|---|---|
| **Compression** | How much smaller the output is | `measure_compression.py` |
| **Throughput** | Bytes per second through the minifier | `measure_throughput.py` |
| **Startup** | Time from process spawn to first output | `measure_startup.py` |
| **Memory** | Peak RSS while minifying | `measure_memory.py` |

Plus head-to-head comparisons against the four most common
alternatives:

| Tool | Language | Notes |
|---|---|---|
| `esbuild` | Go | Fastest minifier; MinifyJS wraps its Go API |
| `terser` | JS (Node) | Slowest of the four; best compression |
| `uglify-js` | JS (Node) | Older, still used in legacy build pipelines |
| `rjsmin` | Python | Smallest; the tool Python users reach for today |
| `minifyjs` | Go | This project |

## Running the suite

```console
cd bench
pip install -r requirements.txt
npm install
python run_all.py
```

`run_all.py` skips any tool whose CLI is not on `PATH`. If you want
to compare against all five, install:

```console
npm install -g esbuild terser uglify-js
pip install rjsmin
pip install minifyjs
```

MinifyJS is always required (the suite would be pointless otherwise).

## Fixtures

Six fixtures of increasing size:

| Fixture | Approx. size | Purpose |
|---|---|---|
| `tiny.js` | < 1 KB | Startup overhead measurement |
| `small.js` | ~5 KB | Single-file throughput |
| `medium.js` | ~50 KB | Typical application module |
| `large.js` | ~500 KB | Large bundle, stress test |
| `application.js` | ~200 KB | Realistic multi-module app |
| `real-world/` | varies | Real library files (downloaded separately) |

The `real-world/` directory is empty in the repository. Add files
there to run the suite against them. Do not commit downloaded
libraries; the directory is in `.gitignore` for that reason.

## Results

`run_all.py` writes:

- `results/latest.md` — a human-readable summary, checked in.

Put the packaged binary on `PATH` so startup and memory measurements
include MinifyJS:

```console
export PATH="$(pwd)/../python/minifyjs/bin:$PATH"
python run_all.py
```

The checked-in `latest.md` is what goes into the README. It is
regenerated on every release, not on every commit.

## Methodology

Each measurement is run in a subprocess so JIT effects from one tool
cannot affect another. Each measurement warms up with N iterations
and reports the median of M measured iterations. Defaults: N=3, M=10.

Timing is wall-clock. CPU time would be more precise on a quiet
machine but is misleading on a loaded CI runner, which is where most
of these runs happen.

Memory is peak RSS as reported by the OS. On Linux and macOS this
comes from `getrusage`; on Windows it comes from `psutil`.

## What this suite does not measure

- **Correctness.** That is what `core/integration/` is for.
- **Cold-start disk I/O.** The measurement includes it, but does
  not isolate it.
- **Parallel execution.** Every measurement runs one file at a time,
  because that is how a typical build script invokes a minifier.
- **Compression against the original source map size.** Only output
  JS bytes are compared.

## Adding a new tool

1. Write a `compare_<tool>.py` that exposes
   `def run(source: str) -> str` returning the minified output.
2. Register the tool in `run_all.py`'s `TOOLS` list.
3. Document it in the table above.

The comparison scripts deliberately do not share a base class.
Each tool has a different invocation shape (Node CLI, Python API,
native binary) and papering over that with an abstraction would
hide exactly the differences the benchmark is meant to expose.