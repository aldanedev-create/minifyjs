# Go API

The `minifyjs/core/api` package is the stable public Go surface of
MinifyJS. External Go callers import it and nothing else.

```go
import "github.com/minifyjs/minifyjs/core/api"
```

Everything under `core/internal/` is private and cannot be imported
by code outside the `core` module (the Go compiler enforces this).
The `api` package is the only supported entry point.

## The public surface

```go
// Functions
func Minify(source string, opts Options) (Result, error)
func Optimize(source string, opts Options) (Result, error)
func Bundle(opts BundleOptions) (Result, error)
func HasErrors(diags []Diagnostic) bool
func SortStable(diags []Diagnostic)

// Types
type Options struct { ... }
type BundleOptions struct { ... }
type Result struct { ... }
type Diagnostic = diagnostics.Diagnostic
type Severity = diagnostics.Severity
type Formatter = diagnostics.Formatter
type TerminalFormatter = diagnostics.TerminalFormatter
type JSONFormatter = diagnostics.JSONFormatter

// Constants
const (
    SeverityInfo    = diagnostics.SeverityInfo
    SeverityWarning = diagnostics.SeverityWarning
    SeverityError   = diagnostics.SeverityError
)

// Errors
var ErrNotImplemented = errors.New("minifyjs: not implemented")
```

Every name in this list is stable and follows semantic versioning.
Anything not in this list is internal.

## Quick reference

| Function | What it does | Filesystem |
|---|---|---|
| `Minify` | Whitespace/comment removal on a single string | No |
| `Optimize` | `Minify` with `compress=true, mangle=true` | No |
| `Bundle` | Multi-file build, writes output to disk | Yes |

## Minimal example

```go
package main

import (
    "fmt"
    "log"

    "github.com/minifyjs/minifyjs/core/api"
)

func main() {
    src := "function add(a, b) { return a + b; }"

    result, err := api.Optimize(src, api.Options{})
    if err != nil {
        log.Fatal(err)
    }
    if result.HasErrors() {
        for _, d := range result.Diagnostics {
            fmt.Printf("%s: %s\n", d.Severity, d.Message)
        }
        return
    }

    fmt.Println(result.Code)
    // function add(a,b){return a+b;}
}
```

## `Options`

```go
type Options struct {
    MinifyWhitespace  bool
    MinifyIdentifiers bool
    MinifySyntax      bool

    Target        string
    Format        string
    Sourcemap     string

    Banner        string
    Footer        string
    LegalComments string

    SourceName    string
}
```

| Field | Type | Zero value means |
|---|---|---|
| `MinifyWhitespace` | bool | Do not remove whitespace |
| `MinifyIdentifiers` | bool | Do not mangle identifiers |
| `MinifySyntax` | bool | Do not fold or eliminate |
| `Target` | string | `"esnext"` (no lowering) |
| `Format` | string | Preserve whatever the input uses |
| `Sourcemap` | string | No source map |
| `Banner` | string | No banner |
| `Footer` | string | No footer |
| `LegalComments` | string | esbuild's default (`"eof"` for transform) |
| `SourceName` | string | `"<stdin>"` in diagnostics |

The zero value of `Options` is a no-op transform: `api.Minify(src,
api.Options{})` returns a normalized but not-minified version of
the input.

## `BundleOptions`

```go
type BundleOptions struct {
    EntryPoints   []string
    OutDir        string
    OutFile       string
    AbsWorkingDir string
    Platform      string
    Bundle        bool
    Splitting     bool

    Options  // embedded
}
```

Bundle-specific fields:

| Field | Notes |
|---|---|
| `EntryPoints` | Required. At least one file path. |
| `OutDir` | Mutually exclusive with `OutFile`. |
| `OutFile` | Mutually exclusive with `OutDir`. |
| `AbsWorkingDir` | Base for relative paths. Defaults to cwd. |
| `Platform` | `"browser"`, `"node"`, or `"neutral"`. |
| `Bundle` | Almost always `true`. |
| `Splitting` | Requires `Format == "esm"`. |

The embedded `Options` supplies minification, target, format, and
sourcemap settings.

## `Result`

```go
type Result struct {
    Code          string
    Map           string
    Diagnostics   []diagnostics.Diagnostic
    OriginalBytes int
    MinifiedBytes int
}

func (r Result) HasErrors() bool
func (r Result) BytesSaved() int
```

| Field | Contents |
|---|---|
| `Code` | The minified JavaScript. |
| `Map` | The source map JSON, or empty. |
| `Diagnostics` | Errors and warnings, in order. |
| `OriginalBytes` | Length of the input in bytes. |
| `MinifiedBytes` | Length of `Code` in bytes. |

`HasErrors` is true if any diagnostic has `Severity ==
api.SeverityError`. Callers should check it before using `Code`.

`BytesSaved` is `OriginalBytes - MinifiedBytes`. Negative when the
output grew.

## `Diagnostic`

`Diagnostic` is an alias for `internal/diagnostics.Diagnostic`,
re-exported so external callers do not need to import an `internal`
package. The two names refer to the same type.

```go
type Diagnostic struct {
    Severity Severity
    Code     string
    Message  string
    File     string
    Line     int
    Col      int
}

func (d Diagnostic) IsError() bool
func (d Diagnostic) HasLocation() bool
```

| Field | Contents |
|---|---|
| `Severity` | `SeverityInfo`, `SeverityWarning`, or `SeverityError` |
| `Code` | A stable error identifier, or empty |
| `Message` | Human-readable description |
| `File` | Source file path, or empty |
| `Line` | 1-based line, or 0 |
| `Col` | 1-based column, or 0 |

`IsError` is `Severity == SeverityError`. `HasLocation` is `Line >
0 && Col > 0`.

## `Severity`

```go
type Severity int

const (
    SeverityInfo Severity = iota
    SeverityWarning
    SeverityError
)

func (s Severity) String() string
```

`String` returns `"info"`, `"warning"`, or `"error"`.

## Errors

`api` returns Go errors for two categories:

1. **Option validation errors.** Returned by `Minify`, `Optimize`,
   and `Bundle` when the options are internally inconsistent (an
   unknown target, an invalid format). These are the caller's fault.
2. **Engine errors.** Returned when the underlying engine cannot
   start (essentially never in practice).

Processing failures (bad JavaScript, missing import) are **not**
returned as `error`. They appear in `Result.Diagnostics`, and
`Result.HasErrors()` is true. This split matters:

- `err != nil` means "the call could not be made".
- `result.HasErrors()` means "the call was made and the input could
  not be processed".

Callers should check both.

```go
result, err := api.Minify(src, opts)
if err != nil {
    // Wrong options; caller's fault.
    return err
}
if result.HasErrors() {
    // The input could not be minified.
    for _, d := range result.Diagnostics {
        log.Println(d)
    }
    return fmt.Errorf("minify failed")
}
```

## Formatters

The terminal and JSON formatters from `internal/diagnostics` are
re-exported:

```go
tf := api.TerminalFormatter{Prefix: "minifyjs: "}
tf.Format(os.Stderr, result.Diagnostics)

jf := api.JSONFormatter{Pretty: true}
jf.Format(os.Stdout, result.Diagnostics)
```

Both implement `api.Formatter`:

```go
type Formatter interface {
    Format(w io.Writer, diags []Diagnostic) error
}
```

`SortStable(diags)` sorts in place by file, then line, then column,
then severity. Both formatters call it internally; callers who want
deterministic output from their own rendering can call it directly.

`HasErrors(diags)` is the same check `Result.HasErrors` uses, but
callable on any slice.

## Bundling example

```go
package main

import (
    "fmt"
    "log"

    "github.com/minifyjs/minifyjs/core/api"
)

func main() {
    result, err := api.Bundle(api.BundleOptions{
        EntryPoints: []string{"src/main.js"},
        OutFile:     "dist/bundle.js",
        Bundle:      true,
        Platform:    "browser",
        Options: api.Options{
            MinifyWhitespace:  true,
            MinifyIdentifiers: true,
            MinifySyntax:      true,
            Target:            "es2015",
            Format:            "esm",
            Sourcemap:         "external",
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    if result.HasErrors() {
        for _, d := range result.Diagnostics {
            fmt.Printf("%s\n", d)
        }
        return
    }
    fmt.Printf("wrote %d bytes\n", result.MinifiedBytes)
}
```

## Stability

Every exported name in `api` follows semantic versioning. A major
version bump is required to remove or rename any of them.

The underlying implementation (`internal/engine`) is not stable.
Callers must not depend on its behavior beyond what `api`
documents.

`api.Diagnostic` and `api.Severity` are type aliases (using `=`),
not distinct types. This means a `diagnostics.Diagnostic` value and
an `api.Diagnostic` value are interchangeable without conversion.
It also means a `Severity` value produced by the engine compares
equal to `api.SeverityError` without conversion.

## What is not in `api`

- **The `Adapter` equivalent.** There is no Go equivalent of the
  Python `Adapter` class. Go callers who want directory-walking
  behavior write it themselves in about thirty lines.
- **Config file loading.** `internal/config` is private. Go
  callers who want to read `minifyjs.config.json` from Go must
  either parse it themselves or shell out to the CLI.
- **Cache management.** `internal/cache` is private.
- **The CLI's flag parsing.** `internal/cmd/minifyjs` is private.

If you need any of these from Go, the answer today is "use the
CLI" or "write it yourself". They may be promoted to `api` if
there is real demand.

## Importing from another module

MinifyJS's Go module is `github.com/minifyjs/minifyjs/core`. It is
importable like any other Go module:

```console
go get github.com/minifyjs/minifyjs/core@v0.1.0
```

```go
import "github.com/minifyjs/minifyjs/core/api"
```

The module's minimum Go version is declared in `core/go.mod`. A
module that imports `api` inherits the same minimum.

The dependency on esbuild is transitive: importing `api` pulls in
`github.com/evanw/esbuild` as an indirect dependency. That is
expected and cannot be avoided; `api` is a wrapper around it.

## See also

- [../python/api.md](../python/api.md) — the Python equivalent
- [../engine/overview.md](../engine/overview.md) — the engine
  package that `api` wraps
- [../architecture.md](../architecture.md) — the full pipeline