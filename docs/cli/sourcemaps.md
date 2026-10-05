# CLI: source maps

Source maps let browser dev tools show you the original source when
you are inspecting minified code. They are a debugging aid and are
not required for the code to run.

## Three modes

| Flag | What happens |
|---|---|
| `--sourcemap=inline` | The map is embedded in the output as a data URL |
| `--sourcemap=external` | The map is written next to the output file |
| `--sourcemap=both` | Both of the above |

The bare flag `--sourcemap` (no `=value`) means `external`. This is
the form most people type.

## Inline mode

```console
minifyjs app.js --sourcemap=inline -o app.min.js
```

The output file ends with:

```javascript
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpb...
```

The base64 blob is the entire source map. No separate file is
created.

**When to use it:** local development, quick debugging, shipping a
demo where you do not want to manage a second file.

**When not to use it:** production. The base64 blob typically
doubles the size of the output. And browsers only load the inline
map when dev tools are open, so production users never benefit.

## External mode

```console
minifyjs app.js --sourcemap=external -o app.min.js
```

Writes two files:

- `app.min.js` — the minified code
- `app.min.js.map` — the source map (JSON)

And the minified code ends with:

```javascript
//# sourceMappingURL=app.min.js.map
```

**When to use it:** production, or any time the output will be
shipped alongside the original source.

## Both mode

```console
minifyjs app.js --sourcemap=both -o app.min.js
```

Writes all three:

- `app.min.js` — the minified code, with an inline map at the end
- `app.min.js.map` — the same map as a separate file
- The inline comment refers to itself; the external file is extra

This is rarely useful. It exists for toolchains that expect both.

## File naming

In external mode, the map file is always named `<output>.map`. If
the output is `dist/app.min.js`, the map is
`dist/app.min.js.map`.

There is no flag to rename the map file. If you need a different
name, rename it after the fact and update the
`sourceMappingURL` comment in the output file.

## What the map contains

A V3 source map with:

- `version: 3`
- `sources`: the original file paths
- `mappings`: the VLQ-encoded position mappings
- `names`: identifiers whose names were changed
- `sourcesContent`: by default, omitted. Pass `--sources-content` to
  embed the original source in the map.

Without `sourcesContent`, the browser must fetch the original source
files to display them. With it, the map is self-contained but
larger.

## `sourcesContent` and privacy

If your original source contains secrets (API keys in comments,
internal URLs), do not use `--sources-content`. The map will contain
the full source, and the map will be served to anyone who can load
the JS.

The safe default is to omit `sourcesContent`. Add it only when you
control who can access the map files.

## Source maps in bundle mode

```console
minifyjs --bundle src/main.js --outfile dist/bundle.js --sourcemap=external
```

Works the same way, but the map's `sources` array lists every
original file that contributed to the bundle. Debugging jumps
straight to the right file, not to the bundle line.

With `--splitting`, one map is written per output file, and shared
chunks get their own maps.

## Multiple inputs

Source maps are per-output. If you run MinifyJS on 100 files
individually, you get 100 maps. There is no "combined" mode for
that; use `--bundle` if you want a single map covering everything.

## How browsers find the map

The minified file ends with a `//# sourceMappingURL=` comment. When
the browser's dev tools are open and you inspect the file, it loads
the map (inline, or from the URL in the comment). It then resolves
each runtime location back to a source location using the mappings.

Browsers only do this when dev tools are open. Users do not pay for
the map at runtime.

## Source map validation

To verify a generated map:

```console
# Any V3 map is JSON.
cat app.min.js.map | python -m json.tool | head
```

You should see `"version": 3` and a `"mappings"` field.

For interactive inspection, load the minified file in a browser with
dev tools open and set a breakpoint. If the breakpoint lands in the
original source (not the minified file), the map is working.

## Compressing source maps

Source maps compress well because the mappings are highly
repetitive:

```console
gzip -9 app.min.js.map
# Typically 5-10x smaller
```

Most web servers serve `.map` files with gzip enabled by default.
If yours does not, enable it.

## See also

- [../python/api.md](../python/api.md) — the Python API's
  `Result.map` field
- [../reference/configuration.md](../reference/configuration.md) —
  the `sourcemap` config key