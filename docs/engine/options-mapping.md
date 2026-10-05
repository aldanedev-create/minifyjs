# Options mapping

A field-by-field table of how MinifyJS's options map to esbuild's.

## Transform options

MinifyJS's `engine.Options` maps to `api.TransformOptions` as
follows. Every MinifyJS field on the left is set by the CLI, the
Python API, or the config loader. Every esbuild field on the right
is set by `engine.Transform`.

| MinifyJS | esbuild | Notes |
|---|---|---|
| `MinifyWhitespace` | `MinifyWhitespace` | Direct |
| `MinifyIdentifiers` | `MinifyIdentifiers` | Direct |
| `MinifySyntax` | `MinifySyntax` | Direct |
| `Target` | `Target` | Direct, passed as `api.Target(string)` |
| `Format` | `Format` | Mapped via a switch: `"esm"` → `api.FormatESModule`, `"cjs"` → `api.FormatCommonJS`, `"iife"` → `api.FormatIIFE`, `""` → `api.FormatDefault` |
| `Sourcemap` | `Sourcemap` | Mapped via a switch: `"inline"` → `api.SourceMapInline`, `"external"` → `api.SourceMapExternal`, `"both"` → `api.SourceMapBoth`, `""` → `api.SourceMapNone` |
| `Banner` | `Banner` | Direct |
| `Footer` | `Footer` | Direct |
| `LegalComments` | `LegalComments` | Mapped via a switch: `"none"` → `api.LegalCommentsNone`, `"inline"` → `api.LegalCommentsInline`, `"eof"` → `api.LegalCommentsEndOfFile`, `"external"` → `api.LegalCommentsLinked`, `""` → `api.LegalCommentsDefault` |
| `SourceName` | `Sourcefile` | Direct, but only set when non-empty. This is the name that appears in diagnostics. |

## Fields the engine sets unconditionally

These are not exposed in MinifyJS's `Options`. They are set by the
engine so their behavior is stable.

| esbuild field | Value | Why |
|---|---|---|
| `Loader` | `api.LoaderJS` | Prevents accidental TypeScript or JSX processing |
| `LogLevel` | `api.LogLevelSilent` | MinifyJS collects diagnostics itself |
| `Sourcemap` | `api.SourceMapNone` unless `Options.Sourcemap` is set | Avoids empty map generation |
| `SourcesContent` | `api.SourcesContentInclude` when sourcemap is on | Self-contained maps |
| `Charset` | `api.CharsetUTF8` | Default; explicit for stability |
| `TreeShaking` | `api.TreeShakingFalse` in Transform | Single file has nothing to shake |

## Build options

MinifyJS's `engine.BuildOptions` maps to `api.BuildOptions`:

| MinifyJS | esbuild | Notes |
|---|---|---|
| `EntryPoints` | `EntryPoints` | Direct |
| `OutDir` | `Outdir` | Direct, mutually exclusive with `OutFile` |
| `OutFile` | `Outfile` | Direct, mutually exclusive with `OutDir` |
| `AbsWorkingDir` | `AbsWorkingDir` | Direct |
| `Platform` | `Platform` | Mapped: `"browser"` → `api.PlatformBrowser`, `"node"` → `api.PlatformNode`, `"neutral"` → `api.PlatformNeutral` |
| `Bundle` | `Bundle` | Direct |
| `Splitting` | `Splitting` | Direct |
| `Options.MinifyWhitespace` | `MinifyWhitespace` | Direct |
| `Options.MinifyIdentifiers` | `MinifyIdentifiers` | Direct |
| `Options.MinifySyntax` | `MinifySyntax` | Direct |
| `Options.Target` | `Target` | Direct |
| `Options.Format` | `Format` | Same switch as Transform |
| `Options.Sourcemap` | `Sourcemap` | Same switch as Transform |
| `Options.Banner` | `Banner["js"]` | esbuild takes a map keyed by output type |
| `Options.Footer` | `Footer["js"]` | Same |
| `Options.LegalComments` | `LegalComments` | Same switch as Transform |

## Fields the Build engine sets unconditionally

| esbuild field | Value | Why |
|---|---|---|
| `Bundle` | `true` if `BuildOptions.Bundle` | The API's purpose is bundling |
| `Write` | `true` | Output goes to disk |
| `LogLevel` | `api.LogLevelSilent` | Same as Transform |
| `Color` | `api.ColorNever` | No escape codes in logs |
| `Metafile` | `false` | Not exposed |

## What is not mapped

The following esbuild fields have no MinifyJS equivalent. If you
need them, use esbuild directly.

### Transform-only options not exposed

| esbuild field | What it does |
|---|---|
| `JSXFactory` | Custom JSX pragma |
| `JSXFragment` | Custom JSX fragment |
| `JSXImportSource` | Custom JSX import source |
| `JSXDev` | Development JSX mode |
| `TsconfigRaw` | Inline tsconfig |
| `Define` | Identifier substitution (exposed as a CLI flag) |
| `Pure` | Purity annotations (exposed as a CLI flag) |
| `Drop` | Statement dropping (exposed as a CLI flag) |
| `MangleProps` | Property name mangling |
| `ReserveProps` | Property name reservation |
| `KeepNames` | Name preservation |
| `LineLimit` | Line length limit |
| `Charset` | Output charset |
| `TreeShaking` | Tree shaking control |
| `IgnoreAnnotations` | Purity comment handling |
| `Supported` | Syntax support overrides |

### Build-only options not exposed

| esbuild field | What it does |
|---|---|
| `Plugins` | The plugin system |
| `Loader` (per-extension) | Custom loaders for `.png`, `.css`, etc. |
| `AssetNames` | Output naming for assets |
| `ChunkNames` | Output naming for chunks |
| `EntryNames` | Output naming for entry points |
| `PublicPath` | URL prefix for assets |
| `Inject` | Auto-importing files |
| `External` | Mark modules as external |
| `Alias` | Path aliases |
| `Tsconfig` | Path to tsconfig |
| `MainFields` | Package.json field priority |
| `Conditions` | Package.json export conditions |
| `PreserveSymlinks` | Symlink handling |
| `NodePaths` | Node module resolution paths |
| `GlobalName` | Global variable name for IIFE output |
| `Metafile` | Build metadata output |
| `Write` | Whether to write output |
| `Watch` | Built-in watch mode |
| `Incremental` | Incremental rebuilds |

The rule: if an option changes the shape of the output in a way
that is specific to a build system or language, MinifyJS does not
expose it. MinifyJS exposes the options that make JavaScript
smaller.

## Adding a new mapped option

To expose an esbuild option through MinifyJS:

1. Add a field to `engine.Options` (or `engine.BuildOptions`).
2. Add validation for it in `Options.validate()`.
3. Add the mapping in `engine.Transform` (or `engine.Build`).
4. Add the field to `api.Options` (or `api.BundleOptions`).
5. Add the mapping in the corresponding `api` function.
6. Add a flag to `cmd/minifyjs/flags.go`.
7. Add a config key to `internal/config/config.go`.
8. Add the corresponding Python keyword argument.
9. Add tests at every layer.
10. Update this file, and `docs/reference/configuration.md`, and
    `docs/cli/reference.md` (regenerate the last one).

That is ten steps. Adding an option is not trivial, and that is
deliberate. The set of options is small so that the surface is
small.

## See also

- [overview.md](overview.md) — the engine package
- [esbuild-integration.md](esbuild-integration.md) — the deeper
  rationale
- [../reference/configuration.md](../reference/configuration.md) —
  the user-facing config reference