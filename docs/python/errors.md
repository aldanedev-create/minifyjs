# Python: errors

The exception hierarchy and how to handle errors from the Python API.

## The hierarchy

```
MinifyJSError                  (base)
├── MinifyError                (processing failed)
│   └── BundleError            (bundling failed)
└── BinaryNotFoundError        (native binary missing)
```

Every exception raised by MinifyJS is a subclass of
`MinifyJSError`. Catching that one exception catches everything:

```python
from minifyjs import MinifyJSError

try:
    result = optimize(source)
except MinifyJSError as e:
    print(f"failed: {e}")
```

Catch a specific type when you want to handle it differently:

```python
from minifyjs import optimize, BundleError, MinifyError, BinaryNotFoundError

try:
    result = optimize(source)
except BinaryNotFoundError:
    print("minifyjs is not installed correctly")
    raise SystemExit(1)
except BundleError as e:
    print(f"bundle failed: {e}")
    raise SystemExit(1)
except MinifyError as e:
    print(f"minify failed: {e}")
    raise SystemExit(1)
```

## `MinifyError`

Raised when the engine cannot process the input.

Common causes:

- The input is not valid JavaScript.
- A required target version is malformed.
- An option has an unrecognized value.

The exception message contains a full diagnostic line when the
location is known:

```python
from minifyjs import optimize, MinifyError

try:
    optimize("function () { } }")
except MinifyError as e:
    print(e)
    # <stdin>:1:17: error: Expected identifier but found "}"
```

To attach a file name, pass `source_name`:

```python
try:
    optimize(source, source_name="src/app.js")
except MinifyError as e:
    print(e)
    # src/app.js:1:17: error: Expected identifier but found "}"
```

## `BundleError`

A subclass of `MinifyError`. Raised when a bundle fails.

Common causes:

- A file in the import graph does not exist.
- A module has a syntax error.
- `outdir` and `outfile` were both passed (this raises `ValueError`
  before the subprocess runs, but the import graph can still fail).
- A module tries to import a Node built-in when
  `platform="browser"`.

```python
from minifyjs import bundle, BundleError

try:
    bundle(["src/main.js"], outfile="dist/bundle.js")
except BundleError as e:
    print(f"build failed: {e}")
```

## `BinaryNotFoundError`

Raised when the bundled native binary cannot be found. This should
not happen after a normal `pip install minifyjs`.

Common causes:

- A source install where `go build` was not run.
- A wheel installed on an unsupported platform (the binary is
  present but not for the current architecture).
- A file permission problem (the binary exists but is not
  executable).

```python
from minifyjs import find_binary, BinaryNotFoundError

try:
    path = find_binary()
    print(f"using {path}")
except BinaryNotFoundError as e:
    print(f"binary not found: {e}")
    raise SystemExit(1)
```

## Warnings vs. errors

Warnings do **not** raise. They appear on `result.diagnostics`:

```python
from minifyjs import optimize

result = optimize(source)
for d in result.diagnostics:
    if d.severity == "error":
        print(f"ERROR: {d}")
    elif d.severity == "warning":
        print(f"WARN:  {d}")
    elif d.severity == "info":
        print(f"INFO:  {d}")
```

`result.has_errors` is a shortcut:

```python
if result.has_errors:
    raise SystemExit(1)
```

## The `Diagnostic` dataclass

```python
@dataclass
class Diagnostic:
    severity: str    # "info" | "warning" | "error"
    message: str
    code: str = ""
    file: str = ""
    line: int = 0
    col: int = 0

    @property
    def is_error(self) -> bool: ...
    @property
    def has_location(self) -> bool: ...
```

`str(diagnostic)` renders the same line the CLI would print:

```python
d = Diagnostic(
    severity="error",
    message="Expected identifier but found \"}\"",
    file="src/app.js",
    line=1,
    col=17,
)
print(str(d))
# src/app.js:1:17: error: Expected identifier but found "}"
```

## Non-MinifyJS exceptions

Three standard Python exceptions can escape from the API:

| Exception | Cause |
|---|---|
| `ValueError` | Bad arguments to `bundle()` (both `outdir` and `outfile`, or neither) |
| `FileNotFoundError` | `Adapter.minify_file()` on a missing file |
| `OSError` | Disk full, permission denied, etc., during `Adapter` I/O |

These are not wrapped. They mean what they say.

## Logging errors

If you want errors on stderr but want to keep going:

```python
import sys
from minifyjs import optimize, MinifyError

results = []
errors = 0

for path in files:
    try:
        result = optimize(path.read_text(), source_name=str(path))
        results.append((path, result))
    except MinifyError as e:
        print(f"skipping {path}: {e}", file=sys.stderr)
        errors += 1

print(f"minified {len(results)} files, {errors} failed")
```

## Testing error handling

In tests, use `pytest.raises`:

```python
import pytest
from minifyjs import optimize, MinifyError

def test_syntax_error_raises():
    with pytest.raises(MinifyError):
        optimize("function () { } }")
```

For diagnostics-only failures (warnings, not errors), do not use
`pytest.raises`; check `result.diagnostics` directly.

## See also

- [api.md](api.md) — the full public surface
- [../cli/errors.md](../cli/errors.md) — the CLI's error format
- [../reference/error-codes.md](../reference/error-codes.md) —
  stable error codes for scripting