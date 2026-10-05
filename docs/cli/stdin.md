# CLI: stdin

How MinifyJS reads input, and what to expect when reading from a
pipe.

## Reading from stdin

Two forms, both equivalent:

```console
cat app.js | minifyjs -o app.min.js
minifyjs - < app.js > app.min.js
```

When no input file is given, MinifyJS reads from stdin. When the
input path is `-`, MinifyJS also reads from stdin. Everything else
is treated as a file path.

## Byte-for-byte input

MinifyJS reads stdin as a byte stream. It does not decode to text,
does not strip a BOM automatically, and does not translate line
endings.

This matters for two reasons:

1. **A UTF-8 BOM at the start of stdin is passed through to the
   parser.** esbuild rejects it with a clear error. If your source
   files have BOMs, strip them first or feed MinifyJS a file path
   instead of a pipe.

2. **CRLF line endings are preserved in the input.** esbuild handles
   them. The output uses LF only.

On Windows, PowerShell's default pipe behavior may translate
`\r\n` to `\n` before MinifyJS sees it. If byte-for-byte input
matters, use a file path rather than a pipe.

## Empty stdin

```console
echo -n "" | minifyjs
```

Produces no output and exits 0. An empty input is a valid input.

## Whitespace-only stdin

```console
printf "   \n\n   " | minifyjs
```

Produces no output and exits 0. esbuild parses an empty program and
MinifyJS writes zero bytes.

## Comments-only stdin

```console
printf "// just a comment\n" | minifyjs
```

Produces no output and exits 0. All comments are removed, leaving
an empty program.

## stdin and shell pipelines

MinifyJS plays well with pipelines:

```console
# Minify, then gzip.
cat app.js | minifyjs --compress | gzip > app.min.js.gz

# Minify a file from a URL.
curl -s https://example.com/app.js | minifyjs --compress -o app.min.js

# Minify the output of a templating engine.
template render app.js.j2 | minifyjs --compress -o app.min.js
```

In every case, stderr goes to the terminal and stdout goes to the
pipe. Progress messages never end up in the piped output.

## What about binary stdin?

If you pipe non-JavaScript bytes into MinifyJS, esbuild reports a
syntax error and MinifyJS exits 1. There is no mode that accepts
binary input.

## stdin with `--watch`

Watch mode requires a file path. `minifyjs --watch` with stdin is a
usage error:

```
minifyjs: --watch requires at least one input file
```

To watch a file, pass it as a path:

```console
minifyjs app.js --watch -o app.min.js
```

## stdin with `--bundle`

Bundle mode requires file paths, because resolving imports needs a
directory to resolve against. `minifyjs --bundle` with stdin is a
usage error.

## Signals and interrupted pipes

If MinifyJS is killed with SIGINT (Ctrl-C) while reading stdin, it
exits with the conventional signal exit code. The partial output is
discarded; the output file (if `-o` was given) is not written.

If the process on the other end of the pipe dies before writing
everything, MinifyJS sees a short read and reports an I/O error,
exit code 1.

## See also

- [stdout.md](stdout.md) — how MinifyJS writes output
- [configuration.md](configuration.md) — reading config from a
  different directory when the input comes from stdin