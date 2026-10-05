#!/usr/bin/env python3
"""A Flask app that minifies its static assets at startup.

This is the pattern most Flask users want: run MinifyJS once when the
app starts, serve the minified files, forget about it.

Run:

    pip install -r requirements.txt
    python app.py

Then visit http://127.0.0.1:5000/.
"""

from __future__ import annotations

from pathlib import Path

from flask import Flask, send_from_directory

from minifyjs import Adapter


HERE = Path(__file__).parent
STATIC = HERE / "static"


class FlaskStaticAdapter(Adapter):
    """Minifies .js files in place. The .min.js files are what the
    HTML actually references."""

    def output_name(self, path: Path) -> Path:
        return path.with_name(f"{path.stem}.min{path.suffix}")


def build_static() -> None:
    """Minify every JS file under static/."""
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
      <head><title>MinifyJS + Flask</title></head>
      <body>
        <h1>MinifyJS + Flask</h1>
        <p>Open the console. The page loaded a minified script.</p>
        <script src="/static/app.min.js"></script>
      </body>
    </html>
    """


@app.route("/static/<path:filename>")
def static_files(filename: str):
    return send_from_directory(STATIC, filename)


def main() -> int:
    build_static()
    print("serving on http://127.0.0.1:5000/")
    app.run(debug=False, port=5000)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())