# Engine overview

The `core/internal/engine` package is the boundary between MinifyJS
and esbuild. It is the only package in the repository that imports
`github.com/evanw/esbuild`. Everything else — the CLI, the Python
bridge, the public Go API — goes through it.

This page explains what the package does and why it exists.

## What the engine package is

Three things:

1. **A small, stable options type.** `engine.Options` has ten
   fields. esbuild's `api.TransformOptions` has about forty. The
   wrapper hides the rest.

2. **A small, stable result type.** `engine.Result` has five
   fields. It is independent of esbuild's `api.TransformResult`, so
   the rest of MinifyJS does not have to import esbuild.

3. **The only place that calls esbuild.** Both
   `api.Transform` and `api.Build` are invoked from this package
   and nowhere else.

## The package's files

```
core/internal/engine/
├── engine.go       package documentation and version
├── options.go      the Options type and its validation
├── result.go       the Result type and derived properties
├── transform.go    wraps esbuild's Transform API
├── build.go        wraps esbuild's Build API
└── errors.go       sentinel errors
```

Plus tests:

```
├── engine_test.go
├── transform_test.go
└── build_test.go
```

## Why a wrapper exists

Three reasons, in order of importance.

### 1. Stability

esbuild's options struct is large and evolves. Fields get added,
renamed, and deprecated. If every MinifyJS package imported esbuild
directly, every esbuild release would require changes throughout
the codebase.

With the wrapper, an esbuild change requires editing exactly one
file. The `engine.Options` type is stable; its mapping to esbuild's
type is not.

### 2. Determinism

esbuild has many options with reasonable defaults for its own use
cases. Some of those defaults are not what MinifyJS wants:

- esbuild's `Transform` defaults to `loader: "js"`, but it also
  accepts TypeScript and JSX. MinifyJS passes `LoaderJS` explicitly
  so `.ts` input fails with a clear error unless the caller opts in.

- esbuild's `LogLevel` defaults to `warning`. MinifyJS sets it to
  `silent` and collects diagnostics itself, so the CLI and the
  Python bridge can present them uniformly.

- esbuild's `LegalComments` defaults depend on the mode. MinifyJS
  makes the default explicit.

Pinning these in one place means they cannot drift.

### 3. Diagnostics

esbuild reports diagnostics as `api.Message` structs with a nested
`Location` type. MinifyJS converts them to its own
`diagnostics.Diagnostic` at the boundary, so the rest of the
codebase never imports esbuild for diagnostics.

## The transform path

`Transform` is the entry point for the common case: a caller has
one JavaScript string and wants a smaller JavaScript string out.

```go
func Transform(source string, opts Options) (Result, error)
```

Flow:

1. `opts.validate()` checks that `Target`, `Format`, `Sourcemap`,
   and `LegalComments` are valid values.
2. The options are mapped to `api.TransformOptions`.
3. `api.Transform` is called.
4. The result is converted to a `Result`, with esbuild's errors and
   warnings folded into a single ordered diagnostic list.

`Transform` never touches the filesystem.

## The build path

`Build` is the entry point for multi-file projects:

```go
func Build(opts BuildOptions) (Result, error)
```

Flow:

1. Structural validation: at least one entry point, exactly one of
   `OutDir` or `OutFile`, `Splitting` requires ESM and bundling.
2. The options are mapped to `api.BuildOptions`.
3. `api.Build` is called.
4. The result is converted to a `Result`.

`Build` writes to the filesystem. It is the only engine function
that does.

## The Options type

```go
type Options struct {
    MinifyWhitespace  bool
    MinifyIdentifiers bool
    MinifySyntax      bool
    Target            string
    Format            string
    Sourcemap         string
    Banner            string
    Footer            string
    LegalComments     string
    SourceName        string
}
```

Ten fields. Every one maps to exactly one esbuild option. See
[options-mapping.md](options-mapping.md) for the details.

## The Result type

```go
type Result struct {
    Code          string
    Map           string
    Diagnostics   []diagnostics.Diagnostic
    OriginalBytes int
    MinifiedBytes int
}
```

Plus three derived methods:

- `HasErrors() bool`
- `HasWarnings() bool`
- `BytesSaved() int`
- `Ratio() float64`

The fields are plain data. The methods are conveniences so callers
do not have to iterate the diagnostics themselves for the common
cases.

## The esbuild version

`engine.Version` is a string that names the esbuild version this
build of MinifyJS wraps.

```go
var Version = "0.28.2"
```

It is set at link time by the release pipeline, so it always
matches the version pinned in `core/go.mod`. The CLI's `--version`
output reads it:

```
minifyjs 0.1.0 (esbuild 0.28.2)
```

This matters for bug reports. A minification bug could be in
MinifyJS's option mapping or in esbuild itself. The version string
disambiguates.

## What the engine does not do

- **Read config files.** That is `internal/config`.
- **Write to disk.** Except for `Build`, which writes build output.
  The single-file `Transform` never touches the filesystem.
- **Parse command-line flags.** That is `cmd/minifyjs`.
- **Validate input JavaScript.** esbuild does that, and the engine
  surfaces the result.
- **Cache results.** That is `internal/cache`, which is not wired
  into the engine yet.

## Reading the code

If you want to understand what MinifyJS does to your JavaScript,
read these three files in this order:

1. `core/internal/engine/options.go` — what options exist
2. `core/internal/engine/transform.go` — how they map to esbuild
3. `core/internal/engine/result.go` — what comes back

That is the whole story. Everything else is plumbing around those
three.

## See also

- [esbuild-integration.md](esbuild-integration.md) — the deeper
  details of the mapping
- [options-mapping.md](options-mapping.md) — the field-by-field
  table
- [../architecture.md](../architecture.md) — the full pipeline