# Python: `optimize()`

The full optimization pipeline. This is what most users want.

## Signature

```python
def optimize(
    source: str,
    *,
    target: Optional[str] = None,
    format: Optional[str] = None,
    sourcemap: Optional[str] = None,
    banner: Optional[str] = None,
    footer: Optional[str] = None,
    legal_comments: Optional[str] = None,
    source_name: Optional[str] = None,
) -> Result
```

## What it does

`optimize()` is exactly:

```python
minify(source, compress=True, mangle=True, **other_args)
```

It enables both the compression pass (constant folding, dead-code
elimination) and the mangling pass (short identifier names).

## The minimal call

```python
from minifyjs import optimize

result = optimize("function calculateTotal(price, tax) { "
                  "  const subtotal = price + (price * tax); "
                  "  return subtotal; }")
print(result.code)
# function calculateTotal(e,t){return e+e*t}
```

The function name is preserved (it is observable from outside), but
the parameters and locals are renamed.

## Typical reduction

For hand-written application code, `optimize()` reduces the file by
**60–75%**.

```python
source = Path("app.js").read_text()
result = optimize(source)
print(f"{result.original_bytes} -> {result.minified_bytes} bytes")
print(f"saved {result.bytes_saved} bytes ({result.ratio * 100:.1f}% of original)")
```

## Passing target

```python
result = optimize(source, target="es2015")
```

Lowers syntax newer than ES2015. Arrow functions become `function`,
`const` becomes `var`, `async` becomes a state machine.

For an ES5 target:

```python
result = optimize(source, target="es5")
```

Everything is lowered. The output runs in any JavaScript engine.

## Full example

```python
from pathlib import Path
from minifyjs import optimize

source = Path("src/app.js").read_text(encoding="utf-8")

result = optimize(
    source,
    target="es2015",
    sourcemap="external",
    banner="/* (c) 2026 Acme Corp */",
    legal_comments="inline",
    source_name="src/app.js",
)

if result.has_errors:
    for d in result.diagnostics:
        print(f"error: {d}")
    raise SystemExit(1)

Path("dist/app.min.js").write_text(result.code, encoding="utf-8")
Path("dist/app.min.js.map").write_text(result.map, encoding="utf-8")

print(f"built: {result.minified_bytes} bytes "
      f"({result.ratio * 100:.1f}% of {result.original_bytes})")
```

## Difference from `minify(compress=True, mangle=True)`

They are identical. `optimize()` exists as a named function because
it reads better at the call site and because the Python API and the
CLI should agree on their vocabulary (`minifyjs --compress --mangle`
produces the same output as `minifyjs.optimize()`).

## When to use `minify()` instead

Use `minify()` when:

- You only need whitespace and comment removal.
- You are debugging and want readable output.
- The input is already optimized and further passes would not help.

Use `optimize()` in every other case.

## Warnings and errors

Same as `minify()`:

- Errors raise `MinifyError`.
- Warnings appear on `result.diagnostics`.

```python
from minifyjs import optimize, MinifyError

try:
    result = optimize(source)
except MinifyError as e:
    print(f"failed: {e}")
    raise SystemExit(1)

for d in result.diagnostics:
    if d.severity == "warning":
        print(f"warning: {d}")
```

## Determinism

`optimize()` is deterministic. Given the same input, options, and
esbuild version, the output is byte-identical across runs and
across machines. This is enforced by
`python/tests/test_idempotence.py`.

## See also

- [minify.md](minify.md) — the underlying function
- [bundle.md](bundle.md) — for multi-file projects
- [configuration.md](configuration.md) — persisting these options