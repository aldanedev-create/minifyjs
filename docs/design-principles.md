# Design principles

Nine rules that govern every decision in MinifyJS. They are listed
in order of priority: when two principles conflict, the earlier one
wins.

## 1. Do one thing

MinifyJS minifies JavaScript. It does not type-check it, bundle CSS,
lint code, transpile JSX, or manage dependencies. Each of those is a
different tool.

The corollary: if a feature request would require MinifyJS to know
about a language other than JavaScript, a framework, or a build
system, it is out of scope.

## 2. Do not require Node.js

Every other JavaScript minifier requires Node.js. This is the single
biggest source of friction for the users MinifyJS targets:
Python developers, backend teams, and CI pipelines that do not
otherwise need a JavaScript toolchain.

The engine is a Go binary. The Python package bundles it. `pip install`
is the whole story.

## 3. Be a wrapper, not an engine

MinifyJS wraps esbuild. It does not implement a JavaScript parser,
optimizer, or mangler. Those are two-year efforts that would produce
a worse tool than esbuild on day one.

The corollary: when a user asks for a feature esbuild does not
provide, the answer is "use esbuild directly", not "let us build
it".

## 4. Small surface, stable surface

The public Go API is `Minify`, `Optimize`, `Bundle`. The public
Python API is `minify`, `optimize`, `bundle`, plus `Options`,
`Result`, and `Adapter`. That is the entire surface.

Everything else is internal and may change between minor versions.

## 5. Fail loudly

When input is invalid, MinifyJS returns an error with a file, line,
and column. It never silently produces wrong output. This is a
deliberate choice: a minifier that silently corrupts code is worse
than one that refuses to run.

This is also why MinifyJS does not strip TypeScript types. esbuild
will strip them without type-checking, which means a `.ts` file with
type errors produces a `.js` file with the same errors, silently. The
docs say so explicitly, and the CLI refuses `.ts` files unless
`--transpile-ts` is passed.

## 6. No configuration by default

`minifyjs app.js -o app.min.js` works without a config file, without
environment variables, and without any prior setup. Every default
is what a typical user wants.

Configuration exists for the users who need it. It does not get in
the way of the users who do not.

## 7. Deterministic output

Given the same input, the same options, and the same esbuild version,
MinifyJS produces byte-identical output. This is enforced by tests.

Determinism matters because build tools are cached, mirrored, and
compared. A build that produces different bytes on different runs is
a build that produces spurious diffs, invalidates caches, and makes
reproducible builds impossible.

## 8. Ship the binary, not the toolchain

End users do not install Go, do not clone the repository, and do not
run `go build`. They run `pip install minifyjs` and get a native
binary for their platform.

The toolchain is a maintainer concern, not a user concern. This
means the release pipeline is more complex than it would be
otherwise, and that is the correct trade.

## 9. Documentation is part of the product

A feature without documentation does not exist. The docs are checked
in alongside the code, reviewed in the same pull requests, and held
to the same standard: short, concrete, and honest about limitations.

The corollary: pull requests that add a flag or a public API without
updating the docs are not ready to merge.