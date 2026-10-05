# FastAPI: build a static directory at startup

A minimal FastAPI app that turns a `static/` directory into a
`build/` directory, minifying every `.js` file along the way, then
mounts `build/` as the static root.

## Run

```console
pip install -r requirements.txt
python main.py
```

Then open <http://127.0.0.1:8000/>.

## What it demonstrates

- Using `Adapter` to build an output directory from a source
  directory
- Copying non-JS files verbatim while minifying JS files
- Mounting the built directory with FastAPI's `StaticFiles`

## The `Adapter` subclass

```python
class FastAPIStaticAdapter(Adapter):
    def __init__(self, source: Path, build: Path):
        super().__init__()
        self.source = source
        self.build = build

    def build_all(self) -> int:
        ...
```

`build_all()` walks the source directory. `.js` files go through
`minify_file()`; everything else is copied with `Path.write_bytes`.
The result is a complete static tree that FastAPI can serve.

## Why build at startup

For a small app, building at startup means the mount point always
reflects the current source. For a production deployment, you would
run `build_all()` as a build step and mount the resulting directory,
so startup does not depend on MinifyJS being installed.