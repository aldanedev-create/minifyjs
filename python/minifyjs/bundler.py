"""bundle() runs a multi-file build through the native binary.

Unlike minify() and optimize(), bundle() writes output to disk
because esbuild's Build API produces one or more files rather than a
single stream. The entry_points and outdir/outfile arguments map
directly to the CLI's --bundle / --outdir / --outfile flags.
"""

from __future__ import annotations

from ._runner import run_bundle
from .options import BundleOptions
from .result import Result


def bundle(
    entry_points: list[str],
    *,
    outdir: str | None = None,
    outfile: str | None = None,
    working_dir: str | None = None,
    platform: str = "browser",
    splitting: bool = False,
    compress: bool = True,
    mangle: bool = True,
    target: str | None = None,
    format: str | None = "esm",
    sourcemap: str | None = None,
    banner: str | None = None,
    footer: str | None = None,
    legal_comments: str | None = None,
) -> Result:
    """Bundle and minify a multi-file JavaScript project.

    Exactly one of ``outdir`` or ``outfile`` must be given.

    Example::

        from minifyjs import bundle
        bundle(["src/main.js"], outdir="dist", format="esm")

    Raises:
        BundleError: if the bundle fails.
        MinifyJSError: if the native binary cannot be found.
    """
    if outdir is None and outfile is None:
        raise ValueError("bundle() requires either outdir or outfile")
    if outdir is not None and outfile is not None:
        raise ValueError("bundle() accepts either outdir or outfile, not both")

    opts = BundleOptions(
        entry_points=entry_points,
        outdir=outdir,
        outfile=outfile,
        working_dir=working_dir,
        platform=platform,
        splitting=splitting,
        compress=compress,
        mangle=mangle,
        target=target,
        format=format,
        sourcemap=sourcemap,
        banner=banner,
        footer=footer,
        legal_comments=legal_comments,
    )
    return run_bundle(opts)
