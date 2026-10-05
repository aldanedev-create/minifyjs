# Python: Minify a Single File

The smallest useful example of the MinifyJS Python API. It reads a JavaScript file, minifies it two ways (whitespace-only and full optimization), prints the compression ratio, and writes the optimized version to disk.

## Run

```console
pip install -r requirements.txt
python app.py
```

## What It Demonstrates

* `minifyjs.minify()` for whitespace-only minification
* `minifyjs.optimize()` for full optimization and mangling
* Reading `Result.bytes_saved`, `Result.ratio`, and `Result.code`
* Writing the output to disk yourself (MinifyJS does not perform file I/O for you in the Python API)

## Expected Output

```text
input:  465 bytes
minify: 361 bytes (77.6% of input)
optim:  210 bytes (45.2% of input, saved 255 bytes)
wrote input.min.js
```

The exact numbers depend on the input file.
