# No Node.js required

This example proves that MinifyJS works on a machine with no Node.js
installed. It exists because "no Node.js required" is a claim worth
verifying, not just asserting.

## Run

```console
chmod +x verify_no_node.sh
./verify_no_node.sh
```

## What it does

1. Removes every `node`, `npm`, or `npx` directory from `PATH`.
2. Fails loudly if `node` or `npm` is still callable.
3. Runs `minifyjs app.js -o /tmp/no-node-output.js`.
4. Checks that the output is non-empty and single-line.

If all four steps succeed, the script prints `PASS: MinifyJS works
without Node.js`.

## Why this matters

Most JavaScript minifiers require Node.js: terser, uglify-js, and
the npm wrapper around esbuild all need a JavaScript runtime to
execute. MinifyJS does not. The engine is a Go binary that is
compiled into the Python wheel and installed by pip. No JavaScript
runtime is involved at any point.

The `verify_no_node.sh` script exists so that this property is
testable, and so that a regression that accidentally introduces a
Node dependency is caught immediately.

## Running on Windows

The script is POSIX shell. On Windows, run it under WSL, Git Bash,
or MSYS2. The equivalent PowerShell version would be:

```powershell
$env:PATH = ($env:PATH -split ';' | Where-Object { $_ -notmatch 'node|npm' }) -join ';'
if (Get-Command node -ErrorAction SilentlyContinue) { throw "node still on PATH" }
minifyjs app.js -o $env:TEMP\no-node-output.js
```

The shell version is the canonical one because the CI workflow that
runs it uses Linux.