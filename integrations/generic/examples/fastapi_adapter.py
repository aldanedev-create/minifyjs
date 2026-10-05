"""FastAPI static assets integration for MinifyJS.

FastAPI serves static files through ``StaticFiles``, which points at
a directory. This adapter minifies the contents of that directory
into a build output directory so the app can mount the build dir
instead of the source dir.

Usage:

    from examples.fastapi_adapter import FastAPIMinifyAdapter
    adapter = FastAPIMinifyAdapter(
        source_dir="app/static",
        build_dir="app/static/dist",
    )
    adapter.build()

Wiring into FastAPI itself is a one-liner the app owns:

    from fastapi.staticfiles import StaticFiles
    app.mount("/static", StaticFiles(directory="app/static/dist"), name="static")

This module does not import FastAPI. It only needs the two paths.
"""

from __future__ import annotations

import sys
from pathlib import Path
from typing import List, Optional

if __package__ is None or __package__ == "":
    sys.path.insert(0, str(Path(__file__).resolve().parents[3]))

from minifyjs import Adapter, Options, Result


class FastAPIMinifyAdapter(Adapter):
    """Copies every non-JS file and minifies every JS file.

    FastAPI apps often serve a directory that mixes HTML, CSS, and
    JS. This adapter preserves the layout: non-JS files are copied
    verbatim, JS files are minified.
    """

    def __init__(
        self,
        source_dir: Path,
        build_dir: Path,
        options: Optional[Options] = None,
    ) -> None:
        super().__init__(options=options or Options(compress=True, mangle=True))
        self.source_dir = Path(source_dir)
        self.build_dir = Path(build_dir)

    def should_process(self, path: Path) -> bool:
        # The base class only matches .js files; the copy step below
        # handles everything else.
        return super().should_process(path)

    def output_name(self, path: Path) -> Path:
        """Preserve the relative path under build_dir, adding the
        minify suffix only to JS files."""
        path = Path(path)
        try:
            rel = path.resolve().relative_to(self.source_dir.resolve())
        except ValueError:
            # Not under source_dir; drop at build_dir root.
            return self.build_dir / f"{path.stem}{self.output_suffix}{path.suffix}"

        if path.suffix == self.extension:
            return self.build_dir / rel.parent / f"{rel.stem}{self.output_suffix}{rel.suffix}"
        return self.build_dir / rel

    def build(self, options: Optional[Options] = None) -> List[Result]:
        """Build the whole output tree."""
        results: List[Result] = []
        for src in self._iter_files(self.source_dir, recursive=True):
            dest = self.output_name(src)
            dest.parent.mkdir(parents=True, exist_ok=True)
            if self.should_process(src):
                results.append(self.minify_file(src, output=dest, options=options))
            else:
                dest.write_bytes(src.read_bytes())
        return results


def _main(argv: List[str]) -> int:
    if len(argv) != 3:
        print(
            "usage: fastapi_adapter.py SOURCE_DIR BUILD_DIR",
            file=sys.stderr,
        )
        return 2

    adapter = FastAPIMinifyAdapter(
        source_dir=argv[1],
        build_dir=argv[2],
    )
    results = adapter.build()

    total_saved = sum(r.bytes_saved for r in results)
    print(
        f"minified {len(results)} JS file(s); saved {total_saved} bytes"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(_main(sys.argv))