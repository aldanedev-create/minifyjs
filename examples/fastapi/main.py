#!/usr/bin/env python3
"""A FastAPI app that builds its static directory at startup.

FastAPI serves static files through `StaticFiles`, which points at a
directory. This example builds that directory from a source
directory: every `.js` file is minified, everything else is copied.

Run:

    pip install -r requirements.txt
    python main.py

Then visit http://127.0.0.1:8000/.
"""

from __future__ import annotations

from pathlib import Path

from fastapi import FastAPI
from fastapi.responses import HTMLResponse
from fastapi.staticfiles import StaticFiles

from minifyjs import Adapter


HERE = Path(__file__).parent
SRC = HERE / "static"
BUILD = HERE / "build"


class FastAPIStaticAdapter(Adapter):
    """Minifies .js and copies everything else into the build dir."""

    def __init__(self, source: Path, build: Path) -> None:
        super().__init__()
        self.source = source
        self.build = build

    def build_all(self) -> int:
        if self.build.exists():
            import shutil
            shutil.rmtree(self.build)
        self.build.mkdir(parents=True)

        count = 0
        for src_file in self.source.rglob("*"):
            if not src_file.is_file():
                continue
            rel = src_file.relative_to(self.source)
            dest = self.build / rel
            dest.parent.mkdir(parents=True, exist_ok=True)
            if src_file.suffix == ".js":
                self.minify_file(src_file, output=dest)
                count += 1
            else:
                dest.write_bytes(src_file.read_bytes())
        return count


def build_static() -> None:
    if not SRC.is_dir():
        return
    adapter = FastAPIStaticAdapter(source=SRC, build=BUILD)
    n = adapter.build_all()
    print(f"built {BUILD} ({n} JS file(s) minified)")


app = FastAPI()


@app.get("/")
def index() -> HTMLResponse:
    return HTMLResponse(
        """
        <!doctype html>
        <html>
          <head><title>MinifyJS + FastAPI</title></head>
          <body>
            <h1>MinifyJS + FastAPI</h1>
            <script src="/static/app.js"></script>
          </body>
        </html>
        """
    )


def main() -> int:
    build_static()
    app.mount("/static", StaticFiles(directory=str(BUILD)), name="static")

    import uvicorn
    print("serving on http://127.0.0.1:8000/")
    uvicorn.run(app, host="127.0.0.1", port=8000, log_level="warning")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())