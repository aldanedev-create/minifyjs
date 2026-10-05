# Flask integration

Minifying a Flask app's static assets.

## The pattern

Flask serves static files from `app.static_folder`, which defaults
to `<app>/static`. The simplest integration minifies every `.js`
file under that folder when the app starts, then serves the
`.min.js` files.

For a real deployment, the minification runs at build time and the
Flask app just serves the pre-built files.

## At startup

The fastest integration: minify when the app object is created.

```python
from pathlib import Path
from flask import Flask

from minifyjs import Adapter, Options

app = Flask(__name__)

STATIC = Path(app.static_folder)


class FlaskAdapter(Adapter):
    def output_name(self, path: Path) -> Path:
        return path.with_name(f"{path.stem}.min{path.suffix}")


with app.app_context():
    FlaskAdapter(options=Options(
        compress=True,
        mangle=True,
        target="es2015",
    )).minify_directory(STATIC)
```

Now `static/app.js` produces `static/app.min.js`, and any template
can reference `{{ url_for('static', filename='app.min.js') }}`.

## At build time

For production, do not minify at startup. Run it as a build step:

```python
#!/usr/bin/env python3
"""Build script for a Flask app."""

from pathlib import Path

from minifyjs import Adapter, Options

STATIC = Path(__file__).parent / "static"

Adapter(options=Options(
    compress=True,
    mangle=True,
    target="es2015",
    sourcemap="external",
)).minify_directory(STATIC)
```

Run it before `flask run` in development, and before containerizing
in production.

## In a Flask CLI command

If you want `flask minify-static` to do the work:

```python
import click
from flask import Flask
from flask.cli import with_appcontext

from minifyjs import Adapter, MinifyJSError, Options

app = Flask(__name__)


@app.cli.command("minify-static")
@click.option("--no-mangle", is_flag=True)
@with_appcontext
def minify_static(no_mangle: bool) -> None:
    """Minify every .js file in the static folder."""
    adapter = Adapter(options=Options(
        compress=True,
        mangle=not no_mangle,
        target="es2015",
    ))

    try:
        results = adapter.minify_directory(app.static_folder)
    except MinifyJSError as e:
        raise click.ClickException(str(e))

    total_in = sum(r.original_bytes for r in results)
    total_out = sum(r.minified_bytes for r in results)
    click.echo(f"minified {len(results)} file(s), saved {total_in - total_out} bytes")
```

Then:

```console
flask minify-static
flask minify-static --no-mangle
```

## Referencing the minified files

Flask's `url_for('static', filename=...)` does not care about the
`.min.js` suffix:

```html
<script src="{{ url_for('static', filename='app.min.js') }}"></script>
```

For a real app, use a variable so the same template works in
development and production:

```python
# config.py
MINIFIED = True   # in production
```

```jinja
{% set suffix = '.min' if config.MINIFIED else '' %}
<script src="{{ url_for('static', filename='app' + suffix + '.js') }}"></script>
```

## Cache busting

Flask serves static files with a `Last-Modified` header but does
not add a content hash to the URL by default. If you want cache
busting, the minified filename needs to change when the content
does:

```python
import hashlib
from pathlib import Path
from minifyjs import Adapter

class HashedAdapter(Adapter):
    def output_name(self, path: Path) -> Path:
        digest = hashlib.sha256(path.read_bytes()).hexdigest()[:8]
        return path.with_name(f"{path.stem}.{digest}{path.suffix}")
```

Then `static/app.js` produces `static/app.a1b2c3d4.js`. Templates
must know the hash, which usually means writing a manifest JSON.

For most Flask apps, this is overkill. `flask-static-digest` or a
CDN that handles cache busting is simpler.

## Blueprint static folders

If your app uses blueprints with their own `static_folder`, run the
adapter once per blueprint:

```python
from minifyjs import Adapter

adapter = Adapter()

for blueprint in app.blueprints.values():
    if blueprint.has_static_folder:
        adapter.minify_directory(blueprint.static_folder)
```

## What not to do

- **Do not minify inside the request handler.** A request should
  never read and rewrite files on disk. Do the work at startup or
  at build time.
- **Do not minify on every import.** If your app is imported
  multiple times (Gunicorn workers, tests), the minification runs
  multiple times. Do it once in a `flask` CLI command or a build
  script.
- **Do not commit `.min.js` files to version control.** They are
  build artifacts. Add `*.min.js` to `.gitignore` and rebuild on
  deploy.

## Complete example

```python
#!/usr/bin/env python3
"""A Flask app that minifies its static assets on startup."""

from pathlib import Path

from flask import Flask, send_from_directory

from minifyjs import Adapter


HERE = Path(__file__).parent
STATIC = HERE / "static"


class FlaskStaticAdapter(Adapter):
    def output_name(self, path: Path) -> Path:
        return path.with_name(f"{path.stem}.min{path.suffix}")


def build_static() -> None:
    if not STATIC.is_dir():
        return
    adapter = FlaskStaticAdapter()
    results = adapter.minify_directory(STATIC)
    print(f"minified {len(results)} JS file(s)")


app = Flask(__name__, static_folder=str(STATIC))


@app.route("/")
def index():
    return """
    <!doctype html>
    <html>
      <body>
        <script src="/static/app.min.js"></script>
      </body>
    </html>
    """


@app.route("/static/<path:filename>")
def static_files(filename: str):
    return send_from_directory(STATIC, filename)


if __name__ == "__main__":
    build_static()
    app.run(debug=False, port=5000)
```

## See also

- [generic-python.md](generic-python.md) — the underlying pattern
- [../python/adapter.md](../python/adapter.md) — the `Adapter` API
- [django.md](django.md), [fastapi.md](fastapi.md) — other frameworks