#!/bin/sh
# The simplest possible use of MinifyJS: read a file, write a smaller
# version of it, print the byte count.
#
# This is what a build script should look like.

set -eu
cd "$(dirname "$0")"

echo "input:  $(wc -c < input.js) bytes"

minifyjs input.js -o input.min.js

echo "output: $(wc -c < input.min.js) bytes"
echo
echo "--- input.min.js ---"
cat input.min.js
echo