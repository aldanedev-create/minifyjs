# ECMAScript compatibility matrix

One subdirectory per ECMAScript version, from `es2015/` through
`es2025/`. Each contains:

- `input.js` — a file that uses every feature introduced in that
  version and every earlier version.
- `expected.min.js` — the output MinifyJS produces with
  `--target es<VERSION>`.
- `README.md` — notes on any feature the version does not support.

The tests run `minifyjs --target es<VERSION>` against each input and
compare against the expected output. This guarantees that a change
to the CLI's option mapping does not silently change which features
get lowered.

A new subdirectory is added when a new ECMAScript version is
published, not before. Guessing at features before they are
finalized is out of scope.