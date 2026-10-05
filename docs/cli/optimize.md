# CLI: optimize

The full optimization pipeline. This is what most users want.

## The two flags that matter

```console
minifyjs app.js --compress --mangle -o app.min.js
```

| Flag | What it adds |
|---|---|
| `--compress` | Constant folding, dead-code elimination, boolean simplification, expression simplification |
| `--mangle` | Renames local variables, function parameters, and closure captures to short names |

Together they typically reduce a file by **60–75%**.

## What each pass does

### `--compress`

- **Constant folding.** `1 + 2 + 3` becomes `6`. `"a" + "b"` becomes
  `"ab"`. `!true` becomes `false`.
- **Dead-code elimination.** `if (false) { ... }` is removed. Code
  after an unconditional `return` is removed.
- **Boolean simplification.** `x === true` becomes `x`. `!!x`
  becomes `x` when `x` is already a boolean.
- **Expression simplification.** `(a, b)` becomes `a, b` at
  statement level. `void 0` becomes `undefined`.
- **Conditional optimization.** `if (a) b; else c;` becomes
  `a ? b : c;` when both branches are expressions.

### `--mangle`

- **Local renaming.** `function calculateTotal(price, tax)` becomes
  `function calculateTotal(e,t)`. The function name itself is
  preserved because it is observable from outside the file.
- **Closure capture renaming.** A variable captured by a nested
  function is renamed consistently in both the outer and inner
  scope.
- **Parameter renaming.** Function parameters become short names,
  except when the function uses `arguments` (which would break).
- **Property names are NOT renamed.** `obj.userId` stays
  `obj.userId`. Mangling property names requires knowing every
  place the property is accessed, which the single-file transform
  cannot know.

## Turning off passes

`--no-compress` and `--no-mangle` exist so a config file's settings
can be overridden:

```json
// minifyjs.config.json
{
  "minify": { "whitespace": true, "identifiers": true, "syntax": true }
}
```

```console
# Config says mangle, but I want readable output this time.
minifyjs app.js --no-mangle -o app.debug.js
```

`--no-minify` turns everything off. The output is esbuild's
normalized version of the input, without any optimization.

## Target environments

```console
minifyjs app.js --target es2015 -o app.min.js
```

Valid values:

- **ECMAScript years:** `es5`, `es2015`, `es2016`, ..., `es2024`,
  `esnext`
- **Browser + version:** `chrome90`, `firefox88`, `safari14`,
  `edge90`
- **Node:** `node18`, `node20`
- **Compound:** `es2020,chrome90,firefox88`

The default is `esnext`, which means "do not lower anything". If you
need older-browser support, you must ask for it.

The lowering is done by esbuild. For example, `--target es5` turns
`const f = (a) => a + 1;` into:

```javascript
var f = function(a) { return a + 1; };
```

## Output format

```console
minifyjs app.js --format esm -o app.min.js
```

| Value | What the output looks like |
|---|---|
| (empty, default) | Preserve whatever the input uses |
| `esm` | ES module syntax (`import`/`export` preserved) |
| `cjs` | CommonJS (`module.exports`, `require`) |
| `iife` | Immediately-invoked function expression |

The `iife` format wraps the output in `(() => { ... })()`. Useful
when the code will be loaded with a plain `<script>` tag and needs
its own scope.

## Content injection

### `--banner` and `--footer`

```console
minifyjs app.js \
    --banner "/* (c) 2026 Acme Corp */" \
    --footer "// built at $(date)" \
    -o app.min.js
```

Both are added verbatim. They are not parsed, minified, or touched
in any way. Use them for license headers and build stamps.

### `--legal-comments`

```console
minifyjs app.js --legal-comments=eof -o app.min.js
```

Controls `/*! ... */` and `//! ... */` comments — the ones with the
extra `!`.

| Mode | Behavior |
|---|---|
| `none` | Removed |
| `inline` | Kept where they are |
| `eof` (default) | Moved to the end of the file |
| `external` | Written to a separate `.LEGAL.txt` file |

### `--define`

```console
minifyjs app.js --define DEBUG=false --define VERSION='"1.0.0"' -o app.min.js
```

Replaces the identifier `DEBUG` with `false` everywhere. This is
how build-time constants work.

The value must be a valid JavaScript expression, quoted as a shell
string. `--define VERSION='"1.0.0"'` produces the string `"1.0.0"` in
the output. `--define PORT=8080` produces the number `8080`.

### `--drop`

```console
minifyjs app.js --drop console --drop debugger -o app.min.js
```

Removes the named statement kinds. Valid values: `console`,
`debugger`.

`--drop console` removes every `console.log(...)`, `console.warn(...)`,
etc. This is the safe way to strip logging from production code.

### `--pure`

```console
minifyjs app.js --pure console.log -o app.min.js
```

Marks a function call as side-effect-free, so the optimizer can
remove it if its result is unused. This is what enables tree-shaking
of specific function calls.

`--pure console.log` means "calling `console.log` has no observable
effect, so if the return value is unused, remove the call". This is
a lie in the strict sense (calling `console.log` does write to
stdout), but it is what the user wants for production builds.

## Combining flags

Flags can appear before or after the input path:

```console
# These are the same.
minifyjs app.js --compress --mangle -o app.min.js
minifyjs --compress --mangle app.js -o app.min.js
```

Values can use `=` or a space:

```console
minifyjs app.js --target=es2015 -o app.min.js
minifyjs app.js --target es2015 -o app.min.js
```

Multiple values for the same flag accumulate (for `--define`,
`--drop`, `--pure`):

```console
minifyjs app.js \
    --define A=1 \
    --define B=2 \
    --pure console.log \
    --pure console.warn \
    -o app.min.js
```

## Examples

### Production build

```console
minifyjs app.js \
    --compress \
    --mangle \
    --target es2015 \
    --drop console \
    --drop debugger \
    --banner "/* (c) 2026 Acme Corp */" \
    --sourcemap=external \
    -o app.min.js
```

### Development build

```console
minifyjs app.js \
    --compress \
    --no-mangle \
    --sourcemap=inline \
    -o app.min.js
```

Mangling off so stack traces are readable. Source map inline so
there is no separate file to ship.

### Library build

```console
minifyjs lib.js \
    --compress \
    --mangle \
    --format esm \
    --target es2015 \
    --banner "/* lib v1.0.0 | MIT */" \
    --legal-comments=inline \
    -o lib.min.js
```

Legal comments inline because libraries often ship license headers
that must survive minification.

## See also

- [minify.md](minify.md) — the whitespace-only mode
- [bundle.md](bundle.md) — multi-file projects
- [configuration.md](configuration.md) — persisting these flags