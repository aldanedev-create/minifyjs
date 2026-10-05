"""optimize() enables all three minification passes.

It is a thin wrapper over minify() that sets compress=True and
mangle=True. It exists as a named function because the two operations
are conceptually distinct: minify() is whitespace/comment removal,
optimize() is that plus syntax rewriting and identifier renaming.
"""

from __future__ import annotations

from typing import Optional

from .minifier import minify
from .result import Result


def optimize(
    source: str,
    *,
    target: Optional[str] = None,
    format: Optional[str] = None,
    sourcemap: Optional[str] = None,
    banner: Optional[str] = None,
    footer: Optional[str] = None,
    legal_comments: Optional[str] = None,
    source_name: Optional[str] = None,
) -> Result:
    """Minify, mangle, and rewrite syntax for the given target.

    Example::

        from minifyjs import optimize
        optimize("const x = 1 + 2 + 3;").code
        # 'const x=6;'
    """
    return minify(
        source,
        compress=True,
        mangle=True,
        target=target,
        format=format,
        sourcemap=sourcemap,
        banner=banner,
        footer=footer,
        legal_comments=legal_comments,
        source_name=source_name,
    )