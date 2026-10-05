"""Django staticfiles integration for MinifyJS.

Django's staticfiles app collects static assets into STATIC_ROOT
when you run `manage.py collectstatic`. This adapter hooks into that
flow: it walks STATICFILES_DIRS, minifies every .js file, and writes
the minified version into STATIC_ROOT/.../.min.js.

Usage from a Django management command or a build script:

    from examples.django_adapter import DjangoMinifyAdapter
    DjangoMinifyAdapter(
        static_root="/path/to/STATIC_ROOT",
        staticfiles_dirs=["/path/to/static"],
    ).minify_directory("/path/to/static")

This module does not import Django. It only needs the paths Django
would give you. That keeps the example usable from a plain Python
build script and testable without Django installed.
"""

from __future__ import annotations

import sys
from pathlib import Path
from typing import Iterable, List, Optional

# Allow running as a script: `python django_adapter.py`.
if __package__ is None or __package__ == "":
    sys.path.insert(0, str(Path(__file__).resolve().parents[3]))

from minifyjs import Adapter, Options, Result


class DjangoMinifyAdapter(Adapter):
    """An adapter that mirrors Django's staticfiles layout.

    Minified output goes to ``static_root`` with the same relative
    path as the input under ``staticfiles_dirs``. That is exactly
    where ``{% static 'app/app.min.js' %}`` will find it after
    ``collectstatic`` runs.
    """

    def __init__(
        self,
        static_root: Path,
        staticfiles_dirs: Optional[Iterable[Path]] = None,
        options: Optional[Options] = None,
    ) -> None:
        super().__init__(options=options or Options(compress=True, mangle=True))
        self.static_root = Path(static_root)
        self.staticfiles_dirs = [Path(p) for p in (staticfiles_dirs or [])]

    def output_name(self, path: Path) -> Path:
        """Mirror the input path under STATIC_ROOT.

        If the input is inside one of the STATICFILES_DIRS, its path
        is preserved relative to that directory. Otherwise the file
        name is used directly.
        """
        path = Path(path)
        for source_dir in self.staticfiles_dirs:
            try:
                rel = path.resolve().relative_to(source_dir.resolve())
            except ValueError:
                continue
            return self.static_root / rel.parent / f"{rel.stem}{self.output_suffix}{rel.suffix}"

        # Not inside any known static dir; drop it at the root of
        # STATIC_ROOT.
        return self.static_root / f"{path.stem}{self.output_suffix}{path.suffix}"

    def minify_directory(
        self,
        directory: Path,
        options: Optional[Options] = None,
        recursive: bool = True,
    ) -> List[Result]:
        """Minify every matching file, creating output directories.

        The base class does not create directories; Django's
        collectstatic assumes they exist. This override creates them
        so a build script can call the adapter against a fresh
        STATIC_ROOT.
        """
        root = Path(directory)
        results: List[Result] = []
        for path in self._iter_files(root, recursive):
            if not self.should_process(path):
                continue
            out = self.output_name(path)
            out.parent.mkdir(parents=True, exist_ok=True)
            results.append(self.minify_file(path, output=out, options=options))
        return results


def _main(argv: List[str]) -> int:
    if len(argv) < 3:
        print(
            "usage: django_adapter.py STATIC_ROOT STATIC_DIR [STATIC_DIR ...]",
            file=sys.stderr,
        )
        return 2

    static_root = Path(argv[1])
    static_dirs = [Path(p) for p in argv[2:]]

    adapter = DjangoMinifyAdapter(
        static_root=static_root,
        staticfiles_dirs=static_dirs,
    )

    total = 0
    saved = 0
    for d in static_dirs:
        for result in adapter.minify_directory(d):
            total += 1
            saved += result.bytes_saved

    print(f"minified {total} file(s); saved {saved} bytes")
    return 0


if __name__ == "__main__":
    raise SystemExit(_main(sys.argv))