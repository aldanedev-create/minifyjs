# Exit codes

The exit codes `minifyjs` returns. These are stable and part of the
public contract.

## The codes

| Code | Name | Meaning |
|---|---|---|
| 0 | OK | Success |
| 1 | Error | Processing error (bad input, I/O, esbuild error) |
| 2 | Usage | Bad command-line usage (unknown flag, missing value, conflicting flags) |
| 3 | NotFound | Input file not found |
| 4 | Internal | Unexpected internal error (a bug) |

No other exit codes are returned. In particular, MinifyJS does not
return `128 + signal` for a killed process; if it is killed by a
signal, the OS reports that directly.

## What triggers each code

### 0 — OK

- The input was minified successfully.
- The output was written (to stdout or to a file).
- Help or version output was requested and printed.

Even when `--quiet` is set, a successful run exits 0.

### 1 — Error

- The input is not valid JavaScript.
- The config file is missing, malformed, or invalid.
- An output file could not be written (permission, disk full).
- A bundle could not be resolved (missing import, circular
  dependency failure).
- An option value is valid syntactically but not semantically (e.g.
  `--target ES2015`).

### 2 — Usage

- An unknown flag was passed (`--not-a-flag`).
- A flag that requires a value was passed without one (`--target`).
- `--outdir` and `--outfile` were both set.
- Multiple input files were given in transform mode.
- `--watch` was given without an input file.
- `--splitting` was given without `--format esm` and `--bundle`.

Usage errors are the caller's fault, not the input's. A script that
retries with a fix should retry only after changing the arguments.

### 3 — NotFound

- The input file path does not exist.
- A config file passed to `--config` does not exist.

This is separated from the general "Error" code (1) so scripts can
distinguish "I typed the filename wrong" from "the file is
broken".

### 4 — Internal

- A bug in MinifyJS.
- An unexpected panic was recovered.
- An invariant was violated.

This should never happen. When it does, MinifyJS prints a stack
trace to stderr and asks the user to report a bug.

## Shell usage

### Check for success

```sh
if minifyjs app.js -o app.min.js; then
    echo "minified successfully"
else
    echo "minify failed"
    exit 1
fi
```

### Distinguish failure kinds

```sh
minifyjs app.js -o app.min.js
code=$?
case $code in
    0) echo "ok" ;;
    1) echo "processing error" >&2 ;;
    2) echo "usage error" >&2 ;;
    3) echo "file not found" >&2 ;;
    4) echo "please report this bug" >&2 ;;
esac
exit $code
```

### Retry on transient failures

```sh
for attempt in 1 2 3; do
    minifyjs app.js -o app.min.js && break
    code=$?
    if [ $code -eq 3 ]; then
        echo "input missing, giving up" >&2
        exit $code
    fi
    if [ $code -ge 4 ]; then
        echo "internal error, giving up" >&2
        exit $code
    fi
    sleep 1
done
```

### Ignore the exit code

```sh
# Run minifyjs but do not fail the script if it errors.
minifyjs app.js -o app.min.js || true
```

## CI usage

In a CI job, the exit code is what determines pass/fail. A
MinifyJS step that exits non-zero fails the job:

```yaml
- run: minifyjs static/app.js --compress --mangle -o static/app.min.js
```

If the step should not fail the job (e.g. a best-effort step), add
`|| true` or `continue-on-error: true`.

## Python usage

The Python API does not return exit codes; it raises exceptions.
The mapping is:

| Exit code | Python equivalent |
|---|---|
| 0 | No exception |
| 1 | `MinifyError` or `BundleError` |
| 2 | `ValueError` (bad args) |
| 3 | `FileNotFoundError` (bad path) |
| 4 | An internal exception (should not happen) |

The Python package's `python -m minifyjs` entry point passes the
underlying binary's exit code through unchanged, so scripts that
call it via `subprocess` see the same codes as the shell.

## Stability

Exit codes are part of MinifyJS's public contract. They will not
change without a major version bump. A script that checks `$? == 1`
will keep working across MinifyJS versions.

New exit codes may be added in a minor version if a genuinely new
category of error appears. Adding a code is a change that is
announced in `CHANGELOG.md`.

## See also

- [error-codes.md](error-codes.md) — finer-grained error
  identifiers
- [../cli/errors.md](../cli/errors.md) — how errors are formatted