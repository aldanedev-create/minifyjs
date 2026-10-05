# Error codes

Stable identifiers for the errors MinifyJS produces. A script that
needs to react to a specific error should match on the code, not on
the message text.

## Currently

MinifyJS does not yet assign stable codes to every diagnostic. The
`code` field on a diagnostic is populated when the underlying
esbuild error has an ID, and empty otherwise.

**Stable today:**

- Diagnostics produced by esbuild carry esbuild's ID (e.g.
  `"parse-error"`).
- MinifyJS's own validation errors do not carry codes yet.

**Planned:**

A stable `mnj.*` code namespace is planned for a future release.
When it lands, every MinifyJS diagnostic will carry a code, and the
codes will be documented on this page.

## Current diagnostic shape

```python
from minifyjs import Diagnostic

d: Diagnostic
d.severity       # "info" | "warning" | "error"
d.code           # esbuild's ID, or ""
d.message        # human-readable description
d.file           # source file path, or ""
d.line           # 1-based line number, or 0
d.col            # 1-based column number, or 0
```

## Matching on error type

Because stable codes are not yet available, match on exception type:

```python
from minifyjs import MinifyError, BundleError, BinaryNotFoundError

try:
    result = optimize(source)
except BinaryNotFoundError:
    # The native binary is missing. Reinstall.
    ...
except BundleError:
    # Multi-file build failed (missing import, syntax error).
    ...
except MinifyError:
    # Single-file processing failed (syntax error, bad option).
    ...
```

## The planned `mnj.*` namespace

When stable codes land, they will follow this pattern:

| Prefix | Meaning | Example |
|---|---|---|
| `mnj.lex.*` | Lexer errors | `mnj.lex.unterminated-string` |
| `mnj.parse.*` | Parser errors | `mnj.parse.unexpected-token` |
| `mnj.opt.*` | Optimizer errors | `mnj.opt.unsupported-syntax` |
| `mnj.target.*` | Target lowering errors | `mnj.target.feature-not-available` |
| `mnj.io.*` | I/O errors | `mnj.io.file-not-found` |
| `mnj.config.*` | Configuration errors | `mnj.config.invalid-value` |
| `mnj.cli.*` | Command-line errors | `mnj.cli.unknown-flag` |

Because the underlying engine is esbuild, most `mnj.lex.*` and
`mnj.parse.*` errors will also carry esbuild's ID in the `code`
field's suffix. Consumers should prefer the `mnj.*` form when both
are present, because the `mnj.*` form is stable and the esbuild ID
is not.

## Matching on message text

As a last resort, match on the message:

```python
if "expected identifier" in str(e).lower():
    ...
```

This is fragile: esbuild's messages change between versions. It is
documented here only as a fallback for code that must support
versions before the stable namespace lands.

## Exit codes vs. error codes

Exit codes are stable today (see
[exit-codes.md](exit-codes.md)). Error codes are the finer-grained
identifier that a script can match on when it needs to distinguish
"parse error at line 5" from "parse error at line 12".

For most scripts, the exit code is enough:

| Exit code | Meaning |
|---|---|
| 0 | Success |
| 1 | Processing error |
| 2 | Bad usage |
| 3 | Input file not found |
| 4 | Internal error |

Match on the exit code first. Match on the error code only when the
exit code is too coarse.

## See also

- [exit-codes.md](exit-codes.md) — the exit codes, stable today
- [../cli/errors.md](../cli/errors.md) — how the CLI formats errors
- [../python/errors.md](../python/errors.md) — the Python exception
  hierarchy