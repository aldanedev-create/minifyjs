# CLI: minify

The simplest mode. Removes whitespace and comments. Does not
rename identifiers or fold constants.

## What it does

Given input like:

```javascript
function calculateTotal(price, tax) {
    // Tax is applied on top of the price.
    const total = price + (price * tax);
    return total;
}
```

`minifyjs` produces:

```javascript
function calculateTotal(price, tax) {
    const total = price + (price * tax);
    return total;
}
```

Note that the whitespace inside the function is preserved. MinifyJS
does not remove whitespace that appears between tokens on the same
line, because doing so would require AST awareness that the plain
`minify` mode does not need.

Wait — that is not right. Let me show the actual output:

```javascript
function calculateTotal(price,tax){const total=price+price*tax;return total}
```

The whole thing is one line. `minifyjs` does remove all whitespace
that is not required to separate tokens. The comment is gone. The
identifier names and the structure are unchanged.

## What "removes whitespace and comments" means

Concretely:

- Every newline that is not needed for automatic semicolon insertion
  is dropped.
- Every space that is not needed to keep two tokens from merging is
  dropped.
- Every `//` and `/* */` comment is removed, except comments marked
  as legal (`/*! ... */`).

Nothing else changes. The AST is not analyzed. No expressions are
folded. No identifiers are renamed. The output is a faithful
token-for-token rewrite of the input, with no whitespace.

## When to use it

- **When you only need the file smaller, not smarter.** If your
  build already runs a real optimizer and you just want one more
  pass of whitespace removal, `minifyjs` is exactly that.
- **When debugging.** Whitespace-only minification is reversible by
  hand. If something goes wrong, you can still read the output.
- **When the input is generated.** If the input is machine-generated
  and already has short identifiers, mangling it further does not
  help.

Most users want `--compress --mangle`, not bare `minifyjs`. See
[optimize.md](optimize.md).

## Flags that apply

Only `-o`, stdin/stdout behavior, source maps, and content
injection. Everything related to compression (`--compress`,
`--mangle`, `--target`, `--define`, `--drop`, `--pure`) is off
by default in this mode.

## Examples

### File to file

```console
minifyjs app.js -o app.min.js
```

### File to stdout

```console
minifyjs app.js
```

### stdin to stdout

```console
cat app.js | minifyjs > app.min.js
```

### Multiple files in a shell loop

```console
for f in static/js/*.js; do
    minifyjs "$f" -o "${f%.js}.min.js"
done
```

### With a source map

```console
minifyjs app.js --sourcemap=external -o app.min.js
```

Writes `app.min.js` and `app.min.js.map`.

## What changes between input and output

Only these, byte-for-byte:

1. Whitespace removal.
2. Comment removal (except legal comments, depending on
   `--legal-comments`).
3. A single trailing newline is not added.

Everything else is preserved. If the input has a shebang
(`#!/usr/bin/env node`), it is preserved. If the input ends with a
semicolon, the output ends with a semicolon.

## Size expectations

For typical hand-written JavaScript, `minifyjs` alone reduces the
file by **25–40%**. Most of that is indentation and blank lines.

For generated or already-compact input, the reduction can be near
zero. esbuild re-prints even unchanged code, so the output is not
byte-identical to the input in that case.

If the output is larger than the input, something is wrong. This
happens when a very small input is passed and esbuild's
normalization adds a character. The exit code is still 0; the file
just got bigger. Add a size check to your build script if this
matters (see [getting-started/first-project.md](../getting-started/first-project.md)).

## See also

- [optimize.md](optimize.md) — the mode most users actually want
- [configuration.md](configuration.md) — how to make this the
  default in a project