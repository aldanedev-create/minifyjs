# Testing

How tests are organized, how to run them, and how to write new ones.

## Four layers

| Layer | Where | What it verifies |
|---|---|---|
| **Unit** | `core/internal/*_test.go`, `python/tests/test_*.py` | Individual functions in isolation |
| **Integration** | `core/integration/` | The compiled CLI as a subprocess |
| **End-to-end** | `test/e2e/` | Full workflows (build, serve, verify) |
| **Compatibility** | `test/compatibility/` | ES version support from es2015 to es2025 |

Plus three cross-cutting mechanisms:

- **Idempotence tests.** Minifying already-minified code produces
  identical output.
- **Cross-language agreement.** The Go CLI and the Python API
  produce byte-identical output for the same input and options.
- **Golden files.** Certain outputs are frozen and compared
  against; regenerated with `UPDATE_GOLDEN=1`.

## Running tests

### All tests

```console
make test
```

Runs the Go test suite and the Python test suite.

### Go only

```console
cd core
go test ./...
```

### Go with race detection

```console
cd core
go test -race ./...
```

Slower but catches data races. Run this before a release.

### Python only

```console
cd python
python -m pytest tests -v
```

### Python tests that need the binary

```console
# Build the binary first.
cd core
go build -o ../python/minifyjs/bin/minifyjs ./cmd/minifyjs
cd ..

# Then run all Python tests.
cd python
python -m pytest tests -v
```

Tests marked `@pytest.mark.binary` skip when the binary is missing,
so a fresh clone can run the suite without building.

### Integration tests

```console
cd core
go test ./integration/... -v
```

First run compiles the CLI once, which takes 30–60 seconds.
Subsequent runs are fast.

### Specific test

```console
cd core
go test ./internal/engine/ -run TestTransformMinifyWhitespace -v
```

```console
cd python
python -m pytest tests/test_minifier.py::test_basic_minify -v
```

## Writing Go tests

Tests live in the same package as the code they test, in a file
named `<file>_test.go`.

```go
package engine

import "testing"

func TestTransformMinifyWhitespace(t *testing.T) {
    src := "function add(a, b) {\n    return a + b;\n}\n"
    res, err := Transform(src, Options{MinifyWhitespace: true})
    if err != nil {
        t.Fatal(err)
    }
    if res.HasErrors() {
        t.Fatalf("diagnostics: %v", res.Diagnostics)
    }
    if len(res.Code) >= len(src) {
        t.Fatalf("output not smaller: %d >= %d", len(res.Code), len(src))
    }
}
```

Rules:

- **Use `t.Fatal` for setup failures** (the test cannot continue).
- **Use `t.Error` for assertions** (the test can report more
  failures before ending).
- **Do not write a test that depends on another test.** Each test
  is independent.
- **Do not use `time.Sleep` for synchronization.** Use channels or
  `t.Cleanup`.

### Table-driven tests

For tests with many cases of the same shape:

```go
func TestTargetValidation(t *testing.T) {
    cases := []struct {
        name  string
        opts  Options
        valid bool
    }{
        {"esnext", Options{Target: "esnext"}, true},
        {"ES2015", Options{Target: "ES2015"}, false},
        {"compound", Options{Target: "es2020,chrome90"}, true},
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            err := c.opts.validate()
            if (err == nil) != c.valid {
                t.Fatalf("validate() = %v, want valid=%v", err, c.valid)
            }
        })
    }
}
```

### Test helpers

Helpers take `*testing.T` and call `t.Helper()` so failure line
numbers point at the caller:

```go
func writeFile(t *testing.T, path, content string) {
    t.Helper()
    if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
        t.Fatal(err)
    }
}
```

## Writing Python tests

Tests live in `python/tests/`, one file per module or concept.

```python
from minifyjs import minify

def test_basic_minify():
    src = "function add(a, b) {\n    return a + b;\n}\n"
    result = minify(src)
    assert result.code == "function add(a,b){return a+b;}"
    assert result.minified_bytes < result.original_bytes
```

Rules:

- **Use plain `assert`.** pytest rewrites them for good failure
  messages.
- **Use `pytest.raises` for exceptions.** Not a try/except.
- **Use `tmp_path` for file tests.** Not `tempfile.mkdtemp`.
- **Use `monkeypatch` for environment manipulation.** Not `os.environ`.

### Marking tests that need the binary

```python
import pytest
from .conftest import BINARY_PATH

pytestmark = pytest.mark.binary

@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")

def test_something_that_runs_the_cli():
    ...
```

Every test file that calls the actual binary uses this pattern.

### Parametrized tests

```python
import pytest

@pytest.mark.parametrize("target,expect_arrow", [
    ("es5", False),
    ("es2015", True),
    ("esnext", True),
])
def test_arrow_lowering(target, expect_arrow):
    result = minify("const f = (a) => a + 1;", target=target)
    assert ("=>" in result.code) == expect_arrow
```

## The idempotence guarantee

Running MinifyJS on its own output produces identical output:

```python
r1 = minify(source)
r2 = minify(r1.code)
assert r1.code == r2.code
```

This is tested for many inputs in
`python/tests/test_idempotence.py` and
`core/integration/minify/idempotence_test.go`. A failure indicates
a real bug: the engine is doing something different on the second
pass.

## Cross-language agreement

The Go CLI and the Python API must produce byte-identical output
for the same input and options:

```python
def test_python_api_matches_cli():
    src = "function add(a, b) { return a + b; }"
    api_result = minify(src, compress=True, mangle=True)

    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs", "--compress", "--mangle"],
        input=src.encode(),
        stdout=subprocess.PIPE,
    )
    assert proc.stdout.decode() == api_result.code
```

Tested in `python/tests/test_integration.py`. A failure means the
Python package is passing different arguments than the CLI, which
is a bug in `_protocol.py`.

## Golden files

For outputs that are correct but tedious to spell out inline:

```python
from helpers.golden import AssertGolden

def test_minify_output_snapshot(t):
    result = minify(complex_input)
    AssertGolden(t, "minify-complex", result.code.encode())
```

Golden files live in `testdata/golden/`. To regenerate:

```console
UPDATE_GOLDEN=1 go test ./...
```

Regenerating produces a diff in the pull request. Reviewers approve
the diff explicitly. That is the point: the golden system makes
every output change visible and reviewed.

## Fuzzing

Go's native fuzzing is used for the lexer, parser, and printer
(future work; the current engine delegates those to esbuild).

```go
func FuzzLexer(f *testing.F) {
    f.Add([]byte("const x = 1;"))
    f.Fuzz(func(t *testing.T, data []byte) {
        // Feed data to the engine; it must not panic.
    })
}
```

Run with:

```console
go test -fuzz=FuzzLexer -fuzztime=30s ./...
```

## Coverage

```console
cd core
go test -cover ./...

cd ../python
python -m pytest --cov=minifyjs tests
```

Coverage is a tool, not a target. The goal is that every public
function is exercised and every error path is reachable by a test.
A line of code that exists only to be covered is a line that should
not exist.

## What tests do not cover

- **esbuild's parser and optimizer.** Those have their own test
  suite. MinifyJS tests that the wrapper passes the right options
  through, not that esbuild handles every edge case.
- **The exact bytes of esbuild's output across versions.** A new
  esbuild version might produce different (usually smaller) output.
  Tests that would break on such a change are written against
  properties ("output is smaller", "output does not contain
  `=>` after targeting es5"), not against exact bytes.
- **The Python interpreter's behavior.** MinifyJS assumes Python
  works.

## See also

- [setup.md](setup.md) — the development environment
- [benchmarking.md](benchmarking.md) — the benchmark suite
- [../../CONTRIBUTING.md](../../CONTRIBUTING.md) — the contribution
  workflow