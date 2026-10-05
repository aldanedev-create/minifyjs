"""The main entry point most users call: minify().

See docs/python/minify.md for the full API reference. This module
never re-implements minification; it prepares arguments, invokes the
native binary, and turns the response into a Result.
"""

from __future__ import annotations

from typing import Optional

from ._runner import run_minify
from .options import Options
from .result import Result


def minify(
    source: str,
    *,
    compress: bool = False,
    mangle: bool = False,
    target: Optional[str] = None,
    format: Optional[str] = None,
    sourcemap: Optional[str] = None,
    banner: Optional[str] = None,
    footer: Optional[str] = None,
    legal_comments: Optional[str] = None,
    source_name: Optional[str] = None,
) -> Result:
    """Minify a single JavaScript source string.

    Example::

        from minifyjs import minify
        minify("function add(a, b) { return a + b; }").code
        # 'function add(a,b){return a+b;}'

    Raises:
        MinifyJSError: if the native binary is not found or fails.
        MinifyError: if the engine cannot process the source.
    """
    opts = Options(
        compress=compress,
        mangle=mangle,
        target=target,
        format=format,
        sourcemap=sourcemap,
        banner=banner,
        footer=footer,
        legal_comments=legal_comments,
        source_name=source_name,
    )
    return run_minify(source, opts)