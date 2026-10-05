# Django integration

Minifying a Django project's static assets.

## The pattern

Django's staticfiles framework collects files from
`STATICFILES_DIRS` into `STATIC_ROOT` when you run `collectstatic`.
The MinifyJS integration adds a step before that: minify every
`.js` file, then let `collectstatic` do its thing.

The MinifyJS work goes in a **management command**, not in a
`STATICFILES_STORAGE` backend. Storage backends are for serving
files, not for rewriting them.

## A `minifystatic` management command

Create
`<yourapp>/management/commands/minifystatic.py`:

```python
"""Minify every JavaScript file in STATICFILES_DIRS into STATIC_ROOT."""

from __future__ import annotations

from pathlib import Path

from django.conf import settings
from django.core.management.base import BaseCommand

from minifyjs import Adapter, MinifyJSError, Options


class DjangoAdapter(Adapter):
    """Writes output to STATIC_ROOT, preserving relative paths."""

    def __init__(self, static_root: Path, static_dirs: list[Path]) -> None:
        super().__init__()
        self.static_root = Path(static_root)
        self.static_dirs = [Path(p) for p in static_dirs]

    def output_name(self, path: Path) -> Path:
        path = Path(path)
        for src in self.static_dirs:
            try:
                rel = path.resolve().relative_to(src.resolve())
            except ValueError:
                continue
            return self.static_root / rel.parent / f"{rel.stem}.min{rel.suffix}"
        return self.static_root / f"{path.stem}.min{path.suffix}"


class Command(BaseCommand):
    help = "Minify every .js file in STATICFILES_DIRS into STATIC_ROOT"

    def add_arguments(self, parser):
        parser.add_argument(
            "--no-mangle",
            action="store_true",
            help="Disable identifier mangling.",
        )
        parser.add_argument(
            "--target",
            default="es2015",
            help="ECMAScript target (default: es2015).",
        )

    def handle(self, *args, **options):
        static_root = Path(settings.STATIC_ROOT)
        static_dirs = [Path(p) for p in settings.STATICFILES_DIRS]

        static_root.mkdir(parents=True, exist_ok=True)

        adapter = DjangoAdapter(static_root, static_dirs)
        adapter.default_options = Options(
            compress=True,
            mangle=not options["no_mangle"],
            target=options["target"],
        )

        total = 0
        saved = 0
        failures = []

        for src in static_dirs:
            if not src.is_dir():
                continue
            try:
                for result in adapter.minify_directory(src):
                    total += 1
                    saved += result.bytes_saved
            except MinifyJSError as e:
                failures.append((src, str(e)))

        self.stdout.write(
            self.style.SUCCESS(
                f"minified {total} file(s); saved {saved} bytes"
            )
        )

        if failures:
            for path, msg in failures:
                self.stderr.write(f"error in {path}: {msg}")
            raise SystemExit(1)
```

## Running it

```console
python manage.py minifystatic
python manage.py minifystatic --no-mangle
python manage.py minifystatic --target es5
python manage.py collectstatic --no-input
```

Run `minifystatic` before `collectstatic`. The minified files go
into the same `STATICFILES_DIRS` entries, and `collectstatic` then
copies both `app.js` and `app.min.js` into `STATIC_ROOT`.

## What `STATICFILES_DIRS` looks like

```python
# settings.py
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent.parent

STATIC_URL = "/static/"
STATIC_ROOT = BASE_DIR / "static_root"
STATICFILES_DIRS = [
    BASE_DIR / "static",
]
```

The command reads `STATICFILES_DIRS` and `STATIC_ROOT` from
settings. There is no other configuration.

## Referencing the minified files

In templates:

```jinja
{% load static %}
<script src="{% static 'js/app.min.js' %}"></script>
```

Django's `{% static %}` tag does not care about the `.min` suffix.

For a project that wants `DEBUG=True` to use unminified files:

```jinja
{% if debug %}
  <script src="{% static 'js/app.js' %}"></script>
{% else %}
  <script src="{% static 'js/app.min.js' %}"></script>
{% endif %}
```

Or a variable in settings:

```python
# settings.py
MINIFIED = not DEBUG
```

```jinja
{% load static %}
{% if minified %}
  <script src="{% static 'js/app.min.js' %}"></script>
{% else %}
  <script src="{% static 'js/app.js' %}"></script>
{% endif %}
```

## The alternate pattern: separate `dist/`

Instead of writing `.min.js` next to the source, some projects
prefer a separate output tree. Override `output_name`:

```python
class DistDjangoAdapter(Adapter):
    def __init__(self, source_dir: Path, dist_dir: Path) -> None:
        super().__init__()
        self.source_dir = source_dir
        self.dist_dir = dist_dir

    def output_name(self, path: Path) -> Path:
        rel = Path(path).resolve().relative_to(self.source_dir.resolve())
        return self.dist_dir / rel
```

Then point `STATICFILES_DIRS` at the `dist/` directory instead of
the source directory. This has the advantage that the source tree
stays clean.

## The `collectstatic` interaction

Django's `collectstatic` copies files from `STATICFILES_DIRS` into
`STATIC_ROOT`. It does not read file contents. So:

1. `minifystatic` writes `app.min.js` alongside `app.js`.
2. `collectstatic` copies both to `STATIC_ROOT`.
3. Templates reference `STATIC_ROOT/js/app.min.js` via
   `{% static %}`.

Both files are present in `STATIC_ROOT`. That is fine; the
unminified one costs bandwidth only if a template references it,
which is usually a `DEBUG`-only path.

## Running in CI

A typical GitHub Actions step:

```yaml
- uses: actions/setup-python@v5
  with:
    python-version: "3.11"
- run: pip install -r requirements.txt
- run: python manage.py minifystatic
- run: python manage.py collectstatic --no-input
- uses: actions/upload-artifact@v4
  with:
    name: static
    path: static_root/
```

The `minifystatic` step runs in seconds. It is safe to run on every
push.

## Common mistakes

### Writing to `STATIC_ROOT` before `collectstatic` deletes it

`collectstatic --clear` deletes everything in `STATIC_ROOT` before
copying. If you run `minifystatic` after `collectstatic --clear` but
before `collectstatic`, the minified files will be deleted.

Run `minifystatic` **before** `collectstatic`:

```console
python manage.py minifystatic
python manage.py collectstatic --no-input
```

### Using a `STATICFILES_STORAGE` backend

Storage backends are for serving static files from storage systems
(S3, GCS, etc.). They are not for rewriting file contents. Do not
try to hook MinifyJS into `STATICFILES_STORAGE`. Use a management
command.

### Forgetting to add the command's app to `INSTALLED_APPS`

Django only finds management commands in apps that are listed in
`INSTALLED_APPS`. If `python manage.py minifystatic` says "Unknown
command", the app containing `management/commands/minifystatic.py`
is not in `INSTALLED_APPS`.

## Complete example

The full example lives in
[`examples/django/`](https://github.com/minifyjs/minifyjs/tree/main/examples/django).
It is a runnable Django project with the `minifystatic` command
already wired up.

## See also

- [generic-python.md](generic-python.md) — the underlying pattern
- [../python/adapter.md](../python/adapter.md) — the `Adapter` API
- [flask.md](flask.md), [fastapi.md](fastapi.md) — other frameworks