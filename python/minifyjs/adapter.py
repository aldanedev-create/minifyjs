"""Generic adapter for build systems and frameworks.

MinifyJS ships no framework-specific code. Instead, this module
defines a small, opinionated base class that any framework adapter
(Django staticfiles, Flask asset pipelines, FastAPI startup hooks,
custom build scripts) can subclass to get minification integrated
with a few lines of code.

The adapter is deliberately thin:

  - It knows how to take a directory of JavaScript files and produce
    minified versions alongside them.
  - It knows how to take a single file and return its minified form.
  - It does not know about templates, URLs, cache invalidation, or
    any framework concept.

Framework-specific behavior lives in the subclass. MinifyJS never
imports a framework and never depends on one.

Example::

    from minifyjs.adapter import Adapter

    class MyBuildAdapter(Adapter):
        def output_name(self, path):
            return path.with_suffix(".min.js")

    adapter = MyBuildAdapter()
    adapter.minify_directory("static/js")

See docs/integrations/generic-python.md for the full guide.
"""

from __future__ import annotations

import os
from pathlib import Path
from typing import Iterable, List, Optional, Union

from .minifier import minify
from .options import Options
from .result import Result

PathLike = Union[str, Path]


class Adapter:
    """Base class for framework-specific MinifyJS integrations.

    Subclasses typically override :meth:`output_name` to control
    where minified output lands, and may override
    :meth:`should_process` to filter files.

    All other methods have sensible defaults that work for the common
    case of "minify every .js file under a directory, write the
    result next to the input with a .min.js suffix".
    """

    #: File extension the adapter processes by default.
    extension: str = ".js"

    #: Suffix inserted before the extension in output file names.
    output_suffix: str = ".min"

    #: Default options passed to minify() unless overridden.
    default_options: Options

    def __init__(self, options: Optional[Options] = None) -> None:
        self.default_options = options or Options(compress=True, mangle=True)

    # ------------------------------------------------------------------
    # Hooks: override these in subclasses to change behavior.
    # ------------------------------------------------------------------

    def should_process(self, path: Path) -> bool:
        """Return True if the adapter should minify ``path``.

        The default implementation processes files with the adapter's
        extension that are not already minified.
        """
        if path.suffix != self.extension:
            return False
        # Skip files that already end with the output suffix.
        stem = path.stem
        if stem.endswith(self.output_suffix):
            return False
        return True

    def output_name(self, path: Path) -> Path:
        """Return the output path for a given input path.

        The default inserts ``output_suffix`` before the extension:
        ``app.js`` becomes ``app.min.js``.
        """
        return path.with_name(f"{path.stem}{self.output_suffix}{path.suffix}")

    # ------------------------------------------------------------------
    # Public API: usable without subclassing for the common case.
    # ------------------------------------------------------------------

    def minify_file(
        self,
        path: PathLike,
        output: Optional[PathLike] = None,
        options: Optional[Options] = None,
    ) -> Result:
        """Minify a single file and write the result.

        If ``output`` is None, ``self.output_name(path)`` is used.
        Returns the Result from the underlying minify() call.
        """
        src_path = Path(path)
        out_path = Path(output) if output is not None else self.output_name(src_path)
        opts = options or self.default_options

        source = src_path.read_text(encoding="utf-8")
        result = minify(
            source,
            compress=opts.compress,
            mangle=opts.mangle,
            target=opts.target,
            format=opts.format,
            sourcemap=opts.sourcemap,
            banner=opts.banner,
            footer=opts.footer,
            legal_comments=opts.legal_comments,
            source_name=str(src_path),
        )
        out_path.write_text(result.code, encoding="utf-8")
        return result

    def minify_directory(
        self,
        directory: PathLike,
        options: Optional[Options] = None,
        recursive: bool = True,
    ) -> List[Result]:
        """Minify every processable file under ``directory``.

        Returns the list of Results, in the order files were visited.
        Files that :meth:`should_process` returns False for are
        skipped silently.
        """
        root = Path(directory)
        paths = self._iter_files(root, recursive)
        return [self.minify_file(p, options=options) for p in paths if self.should_process(p)]

    # ------------------------------------------------------------------
    # Internal helpers.
    # ------------------------------------------------------------------

    def _iter_files(self, root: Path, recursive: bool) -> Iterable[Path]:
        if recursive:
            for dirpath, dirnames, filenames in os.walk(root):
                # Skip common noise directories.
                dirnames[:] = [
                    d for d in dirnames
                    if d not in ("node_modules", ".git", "__pycache__")
                ]
                for name in filenames:
                    yield Path(dirpath) / name
        else:
            for entry in root.iterdir():
                if entry.is_file():
                    yield entry