# CLI: stdout

How MinifyJS writes output, and what to expect when reading it.

## The exact bytes

When writing to stdout, MinifyJS writes **exactly the minified
bytes** and nothing else:

- No trailing newline.
- No leading newline.
- No banner line.
- No progress marker.
- No BOM.

This is deliberate. Standard tools that read a minifier's output
for comparison or further processing need byte-exact output.
Anything extra would need to be stripped, and every consumer would
strip it slightly differently.

If you want a trailing newline in the output file (the POSIX
convention), add it yourself:

```console
minifyjs app.js | tee app.min.js > /dev/null
```

Or, more commonly:

```console
minifyjs app.js -o app.min.js
```

When writing to a file with `-o`, MinifyJS writes the same exact
bytes. No newline is added.

## stdout vs stderr

| Stream | What goes there |
|---|---|
| **stdout** | Exactly the minified JavaScript, or nothing on error |
| **stderr** | Progress messages, warnings, errors |

This means this pipeline works:

```console
minifyjs app.js --compress -o /dev/null 2>/dev/null
```

The progress messages are on stderr, the (empty because of
`/dev/null`) output is on stdout. They never mix.

And this works:

```console
minifyjs app.js --compress 2>/dev/null > app.min.js
```

The progress messages are hidden, the minified code goes to the file.

## Progress messages

By default, MinifyJS prints a summary to stderr when writing to a
file:

```
app.js -> app.min.js: 1234 -> 456 bytes (63.0% smaller)
```

The summary is on stderr, not stdout. When writing to stdout (no `-o`
flag), MinifyJS prints nothing to stderr on success, because there
is nothing useful to say.

To suppress the summary:

```console
minifyjs app.js -o app.min.js --quiet
```

## Verbose mode

```console
minifyjs app.js -o app.min.js --verbose
```

Prints additional information:

```
minifyjs 0.1.0 (esbuild 0.28.2)
input:      app.js
output:     app.min.js
input bytes:  1234
output bytes: 456
saved:        778 (63.0%)
```

All on stderr.

## Writing to a file

When `-o` is given, MinifyJS writes the output **atomically**: it
writes to a sibling temp file and renames it into place. This means:

- The output file is either the old version or the complete new
  version. It is never a partial write.
- A process reading the output file concurrently with MinifyJS sees
  one version or the other, not a mix.

This matters for watch mode and for CI pipelines where another
process may be watching the output directory.

## File permissions

The output file is created with mode `0644` (owner read/write,
group read, world read). If the file already exists, its permissions
are preserved by the atomic-rename path.

On Windows, the mode is ignored; the file inherits the default ACL
of its parent directory.

## Encoding

Output is UTF-8. No BOM is written. If the input was Latin-1, the
output is still UTF-8, because esbuild normalizes to UTF-8 during
parsing.

## What if the output is larger than the input?

It can happen, usually for very small inputs. esbuild's normalized
re-print adds a character here and there. The exit code is still 0.

To detect this in a build script:

```sh
before=$(wc -c < app.js)
minifyjs app.js -o app.min.js
after=$(wc -c < app.min.js)
if [ "$after" -ge "$before" ]; then
    echo "warning: output grew ($before -> $after)" >&2
fi
```

## Redirecting

All standard shell redirections work:

```console
# To a file (the shell truncates first; MinifyJS then writes to it)
minifyjs app.js > app.min.js

# To another command
minifyjs app.js | gzip > app.min.js.gz

# To a specific fd
minifyjs app.js > /dev/fd/3
```

The shell-redirect form (`> app.min.js`) does not use MinifyJS's
atomic-write path. If you want atomicity, use `-o`:

```console
minifyjs app.js -o app.min.js
```

## See also

- [stdin.md](stdin.md) — the input side
- [configuration.md](configuration.md) — controlling the output path
  and quiet/verbose from a config file