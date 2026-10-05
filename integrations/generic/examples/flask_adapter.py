"""Flask static assets integration for MinifyJS.

Flask serves static files from ``app.static_folder`` (by default
``<app>/static``). This adapter minifies every ``.js`` file under
that folder in place, writing ``name.min.js`` next to ``name.js``.

Usage:

    from examples.flask_adapter import FlaskMinifyAdapter
    adapter = FlaskMinifyAdapter(static_folder="myapp/static")
    adapter.minify_directory("myapp/static")

Or, hooked into a Flask CLI command:

    @app.cli.command("minify-static")
    def minify_static():
        FlaskMinifyAdapter(app.static_folder).minify_directory(app.static_folder)

This module does not import Flask. It works with any folder path.
"""

from __future__ import annotations

import sys
from pathlib import Path
from typing import List, Optional

if __package__ is None or __package__ == "":
    sys.path.insert(0, str(Path(__file__).resolve().parents[3]))

from minifyjs import Adapter, Options, Result


class FlaskMinifyAdapter(Adapter):
    """An adapter that writes .min.js next to each .js file."""

    def __init__(
        self,
        static_folder: Path,
        options: Optional[Options] = None,
    ) -> None:
        super().__init__(options=options or Options(compress=True, mangle=True))
        self.static_folder = Path(static_folder)

    def output_name(self, path: Path) -> Path:
        """Write next to the source: ``app.js`` -> ``app.min.js``."""
        return Path(path).with_name(
            f"{Path(path).stem}{self.output_suffix}{Path(path).suffix}"
        )

    def minify_all(self, options: Optional[Options] = None) -> List[Result]:
        """Minify everything under ``static_folder``."""
        return self.minify_directory(self.static_folder, options=options)


def _main(argv: List[str]) -> int:
    if len(argv) != 2:
        print("usage: flask_adapter.py STATIC_FOLDER", file=sys.stderr)
        return 2

    adapter = FlaskMinifyAdapter(static_folder=argv[1])
    results = adapter.minify_all()

    total_saved = sum(r.bytes_saved for r in results)
    print(f"minified {len(results)} file(s); saved {total_saved} bytes")
    return 0


if __name__ == "__main__":
    raise SystemExit(_main(sys.argv))