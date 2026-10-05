"""optimize() enables all three minification passes.

It is a thin wrapper over minify() that sets compress=True and
mangle=True. It exists as a named function because the two operations
are conceptually distinct: minify() is whitespace/comment removal,
optimize() is that plus syntax rewriting and identifier renaming.
"""

from __future__ import annotations

from .minifier import minify
from .result import Result


def optimize(
    source: str,
    *,
    target: str | None = None,
    format: str | None = None,
    sourcemap: str | None = None,
    banner: str | None = None,
    footer: str | None = None,
    legal_comments: str | None = None,
    source_name: str | None = None,
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
