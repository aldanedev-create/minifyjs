# First project

A complete walkthrough: from an empty directory to a project that
minifies its JavaScript on every build. Fifteen minutes, start to
finish.

## What you will build

A small Python web application with:

- One JavaScript file with some common patterns
- A build script that minifies it
- A config file that applies the same options every time
- A check that the minified file is smaller than the original

## Step 1: Set up

```console
mkdir my-project
cd my-project
python -m venv .venv
. .venv/bin/activate       # on Windows: .venv\Scripts\activate
pip install minifyjs
```

Verify:

```console
minifyjs --version
```

## Step 2: Write some JavaScript

Create `static/app.js`:

```javascript
const API_URL = "https://api.example.com/v1";

function fetchUser(id) {
    return fetch(API_URL + "/users/" + id).then(function (response) {
        return response.json();
    });
}

function renderUser(user) {
    const el = document.getElementById("user");
    el.textContent = user.name + " (" + user.email + ")";
}

document.addEventListener("DOMContentLoaded", function () {
    fetchUser(1).then(renderUser);
});
```

## Step 3: Try minifying it manually

```console
minifyjs static/app.js --compress --mangle -o static/app.min.js
```

Look at the output:

```console
cat static/app.min.js
```

It should be one long line, with short variable names, and no
comments. The `wc -c` comparison is instructive:

```console
wc -c static/app.js static/app.min.js
```

You should see something like:

```
  427 static/app.js
  187 static/app.min.js
```

That is a 56% reduction.

## Step 4: Add a config file

Rather than typing the same flags every time, put them in
`minifyjs.config.json`:

```json
{
  "target": "es2015",
  "minify": {
    "whitespace": true,
    "identifiers": true,
    "syntax": true
  },
  "sourcemap": "external"
}
```

Now the same command produces the same output with fewer flags:

```console
minifyjs static/app.js -o static/app.min.js
```

The config is picked up automatically because it is in the current
directory. MinifyJS will also find it from any subdirectory.

## Step 5: Add a build script

Create `build.py`:

```python
#!/usr/bin/env python3
"""Build script for my-project."""

from __future__ import annotations

import sys
from pathlib import Path

from minifyjs import Adapter, MinifyJSError, Options


HERE = Path(__file__).parent
STATIC = HERE / "static"


def main() -> int:
    adapter = Adapter(options=Options(
        compress=True,
        mangle=True,
        target="es2015",
        sourcemap="external",
    ))

    try:
        results = adapter.minify_directory(STATIC)
    except MinifyJSError as e:
        print(f"build failed: {e}", file=sys.stderr)
        return 1

    total_before = sum(r.original_bytes for r in results)
    total_after = sum(r.minified_bytes for r in results)

    print(f"built {len(results)} file(s)")
    print(f"  before: {total_before} bytes")
    print(f"  after:  {total_after} bytes")
    print(f"  saved:  {total_before - total_after} bytes "
          f"({(1 - total_after / total_before) * 100:.1f}% smaller)")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
```

Run it:

```console
python build.py
```

Expected output:

```
built 1 file(s)
  before: 427 bytes
  after:  187 bytes
  saved:  240 bytes (56.2% smaller)
```

## Step 6: Verify the output is smaller

A build step that produces *bigger* output is a bug worth catching.
Add a check to `build.py`:

```python
if total_after >= total_before:
    print("error: minified output is not smaller", file=sys.stderr)
    return 1
```

This matters because with a `--banner` or a very small input, the
output can grow. The check turns a silent regression into a loud
failure.

## Step 7: Wire it into a Makefile

Create `Makefile`:

```makefile
.PHONY: build clean

build:
	python build.py

clean:
	rm -f static/*.min.js static/*.min.js.map

watch:
	minifyjs static/app.js -o static/app.min.js --watch
```

Now `make build` runs the build, `make clean` removes the artifacts,
and `make watch` re-minifies on every save.

## Step 8: Wire it into CI

Create `.github/workflows/build.yml`:

```yaml
name: build

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: "3.11"
      - run: pip install minifyjs
      - run: python build.py
      - uses: actions/upload-artifact@v4
        with:
          name: static
          path: static/*.min.js
```

Every push now produces a minified bundle as a downloadable artifact.

## What you built

A complete build pipeline with:

- A `minifyjs.config.json` that makes the options reproducible
- A `build.py` that minifies a directory and checks the result
- A `Makefile` that wraps it for humans
- A CI workflow that runs it on every push

None of it requires Node.js.

## Where to go next

- **Need bundling?** Replace `Adapter` with `bundle()` in `build.py`.
  See [python/bundle.md](../python/bundle.md).
- **Deploying with Flask?** See
  [integrations/flask.md](../integrations/flask.md).
- **Deploying with Django?** See
  [integrations/django.md](../integrations/django.md).
- **Deploying with FastAPI?** See
  [integrations/fastapi.md](../integrations/fastapi.md).