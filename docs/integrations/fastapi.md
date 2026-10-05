# FastAPI integration

Building a static directory before mounting it.

## The pattern

FastAPI serves static files through `StaticFiles`, which points at
a directory. The MinifyJS integration builds that directory from a
source directory: minify every `.js` file, copy everything else,
then mount the built directory.

The build runs at startup in development and at build time in
production.

## At startup

The simplest version:

```python
from pathlib import Path

from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles

from minifyjs import Adapter, Options


HERE = Path(__file__).parent
SRC = HERE / "static"
BUILD = HERE / "build"


class FastAPIAdapter(Adapter):
    def __init__(self, source: Path, build: Path) -> None:
        super().__init__(options=Options(
            compress=True,
            mangle=True,
            target="es2015",
        ))
        self.source = Path(source)
        self.build = Path(build)

    def output_name(self, path: Path) -> Path:
        rel = Path(path).relative_to(self.source)
        return self.build / rel

    def build_all(self) -> int:
        import shutil
        if self.build.exists():
            shutil.rmtree(self.build)
        self.build.mkdir(parents=True)

        count = 0
        for src in self._iter_files(self.source, recursive=True):
            dest = self.output_name(src)
            dest.parent.mkdir(parents=True, exist_ok=True)
            if self.should_process(src):
                self.minify_file(src, output=dest)
                count += 1
            else:
                dest.write_bytes(src.read_bytes())
        return count


app = FastAPI()

FastAPIAdapter(source=SRC, build=BUILD).build_all()
app.mount("/static", StaticFiles(directory=str(BUILD)), name="static")
```

Every request to `/static/app.js` now serves the minified version.

## At build time

For production, do not build inside the app. Build once, then
serve:

```python
#!/usr/bin/env python3
"""Build script for a FastAPI app."""

from pathlib import Path

from minifyjs import Adapter, Options


HERE = Path(__file__).parent
SRC = HERE / "static"
BUILD = HERE / "static-build"


class BuildAdapter(Adapter):
    def __init__(self, source: Path, build: Path) -> None:
        super().__init__(options=Options(
            compress=True,
            mangle=True,
            target="es2015",
            sourcemap="external",
        ))
        self.source = source
        self.build = build

    def output_name(self, path: Path) -> Path:
        return self.build / Path(path).relative_to(self.source)

    def run(self) -> int:
        import shutil
        if self.build.exists():
            shutil.rmtree(self.build)
        self.build.mkdir(parents=True)

        for src in self._iter_files(self.source, recursive=True):
            dest = self.output_name(src)
            dest.parent.mkdir(parents=True, exist_ok=True)
            if self.should_process(src):
                self.minify_file(src, output=dest)
            else:
                dest.write_bytes(src.read_bytes())
        return 0


if __name__ == "__main__":
    raise SystemExit(BuildAdapter(SRC, BUILD).run())
```

Then in the app:

```python
from fastapi.staticfiles import StaticFiles
app.mount("/static", StaticFiles(directory="static-build"), name="static")
```

## Uvicorn and startup ordering

If you build in the app module (top-level code), the build runs
before the first request, because importing the module runs the
code. This is the right place for a small app in development.

For a large app, the import happens once per worker process.
Building inside every worker means the build runs N times for N
workers. That is wasteful; the build should be a separate step.

## Serving the built directory

```python
from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles

app = FastAPI()
app.mount("/static", StaticFiles(directory="static-build"), name="static")
```

URLs are `/static/<path>`. A file at `static-build/app.js` is
served at `/static/app.js`.

If you want the source directory to be served in development and
the build directory in production, choose at mount time:

```python
import os
from fastapi.staticfiles import StaticFiles

mount_dir = "static-build" if os.environ.get("ENV") == "production" else "static"
app.mount("/static", StaticFiles(directory=mount_dir), name="static")
```

## Cache headers

FastAPI's `StaticFiles` does not add cache-control headers by
default. For production, add them:

```python
from starlette.staticfiles import StaticFiles as StarletteStatic

class CachedStaticFiles(StarletteStatic):
    def file_response(self, *args, **kwargs):
        response = super().file_response(*args, **kwargs)
        response.headers["Cache-Control"] = "public, max-age=31536000, immutable"
        return response


app.mount("/static", CachedStaticFiles(directory="static-build"), name="static")
```

The `immutable` directive tells browsers never to revalidate. That
requires the URL to change when the content changes, which means
content-hashed filenames:

```python
import hashlib
from pathlib import Path
from minifyjs import Adapter

class HashedAdapter(Adapter):
    def output_name(self, path: Path) -> Path:
        digest = hashlib.sha256(path.read_bytes()).hexdigest()[:8]
        return Path(path).with_name(f"{path.stem}.{digest}{path.suffix}")
```

If you add hashing, you also need a manifest that maps logical
names to hashed names, so templates or API responses can emit the
right URL. That is usually more machinery than a small app wants;
long cache times with `no-cache` for the HTML is simpler.

## Code splitting and bundling

If your app has modules that import each other, use `bundle()` for
the entry point:

```python
from minifyjs import bundle

bundle(
    entry_points=["src/main.js"],
    outdir="static-build/js",
    format="esm",
    target="es2015",
    splitting=True,
)
```

Then mount `static-build/` and reference the entry point files.

`Adapter` and `bundle()` compose: use `Adapter` for standalone
scripts in `static/js/`, and `bundle()` for the module graph in
`src/`.

## Complete example

```python
#!/usr/bin/env python3
"""A FastAPI app that builds its static directory at startup."""

from pathlib import Path

from fastapi import FastAPI
from fastapi.responses import HTMLResponse
from fastapi.staticfiles import StaticFiles

from minifyjs import Adapter, Options


HERE = Path(__file__).parent
SRC = HERE / "static"
BUILD = HERE / "static-build"


class StaticBuildAdapter(Adapter):
    def __init__(self, source: Path, build: Path) -> None:
        super().__init__(options=Options(
            compress=True,
            mangle=True,
            target="es2015",
        ))
        self.source = Path(source)
        self.build = Path(build)

    def output_name(self, path: Path) -> Path:
        return self.build / Path(path).relative_to(self.source)

    def build_all(self) -> int:
        import shutil
        if self.build.exists():
            shutil.rmtree(self.build)
        self.build.mkdir(parents=True)

        count = 0
        for src in self._iter_files(self.source, recursive=True):
            dest = self.output_name(src)
            dest.parent.mkdir(parents=True, exist_ok=True)
            if self.should_process(src):
                self.minify_file(src, output=dest)
                count += 1
            else:
                dest.write_bytes(src.read_bytes())
        return count


app = FastAPI()
StaticBuildAdapter(SRC, BUILD).build_all()
app.mount("/static", StaticFiles(directory=str(BUILD)), name="static")


@app.get("/")
def index() -> HTMLResponse:
    return HTMLResponse(
        "<!doctype html><html><body>"
        '<script src="/static/app.js"></script>'
        "</body></html>"
    )


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="127.0.0.1", port=8000)
```

## See also

- [generic-python.md](generic-python.md) — the underlying pattern
- [../python/adapter.md](../python/adapter.md) — the `Adapter` API
- [../python/bundle.md](../python/bundle.md) — for module graphs
- [flask.md](flask.md), [django.md](django.md) — other frameworks