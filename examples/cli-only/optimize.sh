#!/bin/sh
# The same file, but with full optimization: syntax rewriting,
# dead-code elimination, and identifier mangling.

set -eu
cd "$(dirname "$0")"

echo "input:  $(wc -c < input.js) bytes"

minifyjs input.js --compress --mangle -o input.opt.js

echo "output: $(wc -c < input.opt.js) bytes"
echo
echo "--- input.opt.js ---"
cat input.opt.js
echo