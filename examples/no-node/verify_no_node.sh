#!/bin/sh
# Verify that MinifyJS works without Node.js.
#
# This script:
#
#   1. Scrubs any node / npm / npx directory from PATH.
#   2. Checks that `node` is genuinely unavailable.
#   3. Runs minifyjs on app.js.
#   4. Checks that the output is non-empty and looks minified.
#
# If all four steps succeed, MinifyJS has demonstrated that it does
# not depend on Node.js.

set -eu
cd "$(dirname "$0")"

# Strip any PATH entry that contains node, npm, or npx.
CLEAN_PATH=""
OLD_IFS="$IFS"
IFS=":"
for entry in $PATH; do
    case "$entry" in
        *node*|*npm*|*npx*) continue ;;
    esac
    if [ -z "$CLEAN_PATH" ]; then
        CLEAN_PATH="$entry"
    else
        CLEAN_PATH="$CLEAN_PATH:$entry"
    fi
done
IFS="$OLD_IFS"

export PATH="$CLEAN_PATH"

echo "== 1. scrubbed PATH"
echo "   PATH=$PATH"

echo
echo "== 2. checking that node is unavailable"
if command -v node >/dev/null 2>&1; then
    echo "   FAIL: node is still on PATH"
    exit 1
fi
if command -v npm >/dev/null 2>&1; then
    echo "   FAIL: npm is still on PATH"
    exit 1
fi
echo "   OK: node and npm are not on PATH"

echo
echo "== 3. running minifyjs"
minifyjs app.js -o /tmp/no-node-output.js

echo
echo "== 4. checking output"
size=$(wc -c < /tmp/no-node-output.js)
if [ "$size" -eq 0 ]; then
    echo "   FAIL: output is empty"
    exit 1
fi
if grep -q $'\n' /tmp/no-node-output.js; then
    echo "   FAIL: output contains newlines"
    exit 1
fi
echo "   OK: output is $size bytes and single-line"

echo
echo "PASS: MinifyJS works without Node.js"