# CLI: errors

What MinifyJS does when something goes wrong, and how to read its
error messages.

## Exit codes

| Code | Meaning | Typical cause |
|---|---|---|
| 0 | Success | Nothing went wrong |
| 1 | Processing error | Bad JavaScript, I/O error, esbuild error |
| 2 | Bad usage | Unknown flag, missing value, conflicting flags |
| 3 | Input file not found | Path does not exist |
| 4 | Internal error | A bug in MinifyJS |

Exit codes are stable. A CI script that checks `$? == 1` will keep
working across MinifyJS versions.

## Where errors go

All errors go to **stderr**. stdout receives only the minified
output, and only on success. On error, stdout is empty.

This means the exit code is the only signal in a pipeline:

```sh
minifyjs app.js > app.min.js
if [ $? -ne 0 ]; then
    echo "minify failed, app.min.js was not modified" >&2
    exit 1
fi
```

## Reading an error message

Errors have this shape:

```
minifyjs: FILE:LINE:COL: severity: message
```

Example:

```
minifyjs: src/app.js:12:5: error: Expected identifier but found "="
```

| Part | Meaning |
|---|---|
| `minifyjs: ` | The prefix that identifies MinifyJS output |
| `src/app.js` | The file the error is in |
| `12:5` | Line 12, column 5 (1-based) |
| `error:` | Severity (error, warning, or info) |
| `Expected identifier...` | The message from esbuild |

For stdin, the file is `<stdin>`:

```
minifyjs: <stdin>:3:1: error: Unexpected end of file
```

## Common errors

### Syntax error

```
minifyjs: app.js:12:5: error: Expected identifier but found "="
```

The input is not valid JavaScript. Read the line and column, fix the
source, rerun. MinifyJS does not attempt to recover from syntax
errors; it exits with code 1.

### Unterminated string or template

```
minifyjs: app.js:8:15: error: Unterminated string literal
```

The string is missing a closing quote. This can also happen if the
string contains an unescaped newline.

### Unterminated comment

```
minifyjs: app.js:20:1: error: Unterminated comment
```

A `/*` without a matching `*/`.

### Duplicate declaration

```
minifyjs: app.js:5:5: error: The symbol "x" has already been declared
```

`let x; let x;` or `const x = 1; const x = 2;`. JavaScript does not
allow redeclaration of `let` and `const` bindings in the same scope.

### Reserved word as identifier

```
minifyjs: app.js:3:7: error: Expected identifier but found "if"
```

Using a reserved word where an identifier is expected. This often
follows a missing semicolon or a stray operator.

### Invalid target

```
minifyjs: invalid target "ES2015"
```

Target strings are case sensitive. Use `es2015`, not `ES2015`.

### Unknown flag

```
minifyjs: unknown flag "--targett"
```

Exit code 2. The flag name is misspelled. Check the spelling or run
`minifyjs --help`.

### Missing flag value

```
minifyjs: flag --target requires a value
```

Exit code 2. A flag that takes a value was passed without one.

### Config file not found

```
minifyjs: config: open /path/to/config.json: no such file or directory
```

Exit code 1. `--config` pointed at a path that does not exist.

### Config file is invalid JSON

```
minifyjs: config: parse: invalid character '}' looking for beginning of value
```

Exit code 1. The config file has a JSON syntax error.

### Config file has an invalid value

```
minifyjs: invalid format "bogus" (want iife, cjs, or esm)
```

Exit code 1. The config file's `format` key has a value that is not
one of the allowed options.

### Input file not found

```
minifyjs: app.js: no such file
```

Exit code 3. The input path does not exist.

## Warnings

Warnings do not stop the build. They print on stderr and the exit
code is still 0.

Example:

```
minifyjs: app.js:1:1: warning: This feature is not available in the target environment
```

Warnings are usually about features that will not work in an older
target, or about code that the optimizer cannot simplify.

To suppress warnings without suppressing errors, use `--quiet`. To
see all warnings, use `--verbose`.

## When the output is not smaller

MinifyJS does not warn when the output is the same size as the
input, or larger. This can happen with tiny inputs. If your build
needs to guarantee a reduction, add your own check:

```sh
before=$(wc -c < app.js)
minifyjs app.js -o app.min.js || exit 1
after=$(wc -c < app.min.js)
if [ "$after" -ge "$before" ]; then
    echo "error: output is not smaller ($before -> $after)" >&2
    exit 1
fi
```

## Internal errors

Exit code 4 means MinifyJS has a bug. This should never happen.
When it does, the error message includes a stack trace:

```
minifyjs: internal error: runtime error: index out of range [5] with length 3
goroutine 1 [running]:
...
```

Please report it with the input that triggered it. See
[../../CONTRIBUTING.md](../../CONTRIBUTING.md) for how to file a bug.

## See also

- [../reference/exit-codes.md](../reference/exit-codes.md) — every
  code, in detail
- [../reference/error-codes.md](../reference/error-codes.md) —
  stable error IDs for scripting