# MinifyJS examples

Every directory here is a **runnable project**, not a code snippet.
Copy one, run its command, and it works. Each example is deliberately
minimal so the interesting lines are not buried under boilerplate.

## What each example shows

| Example | Framework | What it demonstrates |
|---|---|---|
| `cli-only/` | none | Using the `minifyjs` binary directly from a shell |
| `python-basic/` | none | Using the Python API for a single file |
| `python-bundler/` | none | Using the Python API to bundle a multi-file project |
| `flask/` | Flask | Minifying a Flask app's static assets at startup |
| `django/` | Django | Integrating MinifyJS with Django's staticfiles pipeline |
| `fastapi/` | FastAPI | Building a static directory before mounting it |
| `no-node/` | none | Proving MinifyJS works on a machine with no Node.js |

## Requirements

- **Python 3.8+** for every Python example.
- **The `minifyjs` binary or the `minifyjs` Python package** for
  every example that runs MinifyJS.
- **The framework** for the framework examples (`flask`, `django`,
  `fastapi`). Each framework example has its own
  `requirements.txt`.

## Installing MinifyJS

If you do not have MinifyJS installed yet:

```console
pip install minifyjs
```

Or, from a source checkout, build the binary first:

```console
cd core
go build -o ../python/minifyjs/bin/minifyjs ./cmd/minifyjs
cd ..
pip install -e python
```

## Running an example

```console
cd examples/python-basic
python app.py
```

Every example prints something useful and exits 0. None of them
require any configuration.

## None of these examples require Node.js

The whole point of MinifyJS is that it works without Node.js. The
`no-node/` example verifies this explicitly by running MinifyJS with
`PATH` scrubbed of any Node-related binary.

## Adding a new example

1. Create `examples/<name>/`.
2. Include a `README.md` that explains what it demonstrates and how
   to run it.
3. Include a `requirements.txt` if it has Python dependencies.
4. Keep it under 100 lines of code. If it needs more, it belongs in
   `docs/`, not here.