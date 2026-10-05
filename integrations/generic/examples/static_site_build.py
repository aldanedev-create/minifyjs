"""A framework-agnostic static site build step.

Given a directory of source files, this script produces a `dist/`
directory that contains:

  - every HTML file copied verbatim
  - every CSS file copied verbatim (a future MinifyJS version may
    handle CSS; today it does not)
  - every JS file minified with the .min.js suffix removed, so the
    HTML's `<script src="app.js">` keeps working

Usage:

    python static_site_build.py SOURCE_DIR [DIST_DIR]

Defaults DIST_DIR to SOURCE_DIR + "/dist".
"""

from __future__ import annotations

import shutil
import sys
from pathlib import Path
from typing import List, Optional

if __package__ is None or __package__ == "":
    sys.path.insert(0, str(Path(__file__).resolve().parents[3]))

from minifyjs import Adapter, Options, Result


class StaticSiteAdapter(Adapter):
    """Copies non-JS assets, minifies JS assets, and strips the
    .min suffix so URLs in the HTML stay unchanged."""

    def __init__(
        self,
        source_dir: Path,
        dist_dir: Path,
        options: Optional[Options] = None,
    ) -> None:
        super().__init__(options=options or Options(compress=True, mangle=True))
        self.source_dir = Path(source_dir)
        self.dist_dir = Path(dist_dir)

    def output_name(self, path: Path) -> Path:
        """Minified JS keeps its original name so the HTML doesn't
        need to change."""
        path = Path(path)
        try:
            rel = path.resolve().relative_to(self.source_dir.resolve())
        except ValueError:
            return self.dist_dir / path.name
        # No .min suffix in the output; the whole dist/ is minified.
        return self.dist_dir / rel

    def build(self, options: Optional[Options] = None) -> List[Result]:
        results: List[Result] = []
        if self.dist_dir.exists():
            shutil.rmtree(self.dist_dir)
        self.dist_dir.mkdir(parents=True)

        for src in self._iter_files(self.source_dir, recursive=True):
            dest = self.output_name(src)
            dest.parent.mkdir(parents=True, exist_ok=True)
            if self.should_process(src):
                results.append(self.minify_file(src, output=dest, options=options))
            else:
                shutil.copy2(src, dest)
        return results


def _main(argv: List[str]) -> int:
    if len(argv) < 2 or len(argv) > 3:
        print("usage: static_site_build.py SOURCE_DIR [DIST_DIR]", file=sys.stderr)
        return 2

    source_dir = Path(argv[1]).resolve()
    dist_dir = Path(argv[2]).resolve() if len(argv) == 3 else source_dir / "dist"

    if not source_dir.is_dir():
        print(f"error: {source_dir} is not a directory", file=sys.stderr)
        return 1

    adapter = StaticSiteAdapter(source_dir=source_dir, dist_dir=dist_dir)
    results = adapter.build()

    js_saved = sum(r.bytes_saved for r in results)
    print(f"built {dist_dir}")
    print(f"minified {len(results)} JS file(s); saved {js_saved} bytes")
    return 0


if __name__ == "__main__":
    raise SystemExit(_main(sys.argv))