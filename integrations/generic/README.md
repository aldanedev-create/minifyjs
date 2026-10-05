# MinifyJS generic integrations

MinifyJS ships no framework-specific code. Every integration lives
outside the engine: the framework's build step subclasses
`minifyjs.Adapter` and lets MinifyJS handle the actual minification.

This directory contains working examples for the most common
patterns. Each example is a self-contained script you can copy into
your own project and adapt.

## Why this design

Frameworks change. MinifyJS does not want to grow a `django/` folder,
then a `flask/` folder, then a `fastapi/` folder, then a folder for
whatever comes next. Instead, MinifyJS exposes one small base class:

```python
from minifyjs import Adapter

class MyAdapter(Adapter):
    def output_name(self, path):
        return path.with_suffix(".min.js")
```

That is the entire integration surface. Anything more — knowing
about Django's `STATIC_ROOT`, Flask's `static_folder`, FastAPI's
mount points, or a CI system's artifact paths — belongs in the
framework adapter, not in MinifyJS.

## The examples

| Example | What it shows |
|---|---|
| `django_adapter.py` | Subclassing `Adapter` for Django's `STATICFILES_DIRS` / `STATIC_ROOT` |
| `flask_adapter.py` | Subclassing `Adapter` for Flask's `app.static_folder` |
| `fastapi_adapter.py` | Subclassing `Adapter` for FastAPI's `StaticFiles` mount points |
| `static_site_build.py` | A framework-agnostic build step for a static site |
| `ci_pipeline.py` | Using MinifyJS inside a GitHub Actions / GitLab CI job |
| `watch_directory.py` | A long-running watcher using MinifyJS's per-file API |

Each example is runnable:

```console
python integrations/generic/examples/static_site_build.py path/to/static
```

The framework examples (`django_adapter.py`, etc.) exit with a clear
message if the framework is not installed, so they are safe to run
anywhere.

## Writing your own adapter

The base class has three extension points:

1. **`extension`** — which files to process (default `.js`).
2. **`output_suffix`** — inserted before the extension (default `.min`).
3. **`should_process(path)`** — return True/False for a given file.
4. **`output_name(path)`** — where the minified output goes.

The only method most adapters override is `output_name`, because
different frameworks have different conventions for where build
artifacts live.

## The full API

```python
from minifyjs import Adapter, Options

class MyAdapter(Adapter):
    extension = ".js"
    output_suffix = ".min"

    def __init__(self, target_dir, **kwargs):
        super().__init__(options=Options(compress=True, mangle=True))
        self.target_dir = target_dir

    def should_process(self, path):
        # Skip files the framework itself generates.
        return super().should_process(path) and path.name not in {"service-worker.js"}

    def output_name(self, path):
        return self.target_dir / path.name.replace(".js", ".min.js")

adapter = MyAdapter(target_dir="static/dist")
adapter.minify_directory("static/js")
```

See `minifyjs/adapter.py` for the full base-class documentation.

## Tests

The tests under `integrations/generic/tests/` verify that each
example adapter behaves correctly without requiring the framework
to be installed. The Django example, for instance, is tested with a
stub `STATIC_ROOT` path — the test does not import Django.