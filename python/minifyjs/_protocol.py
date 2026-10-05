"""Defines the CLI flags used to talk to the native binary from Python.

Kept as a single place to change if/when the binary grows a richer
(e.g. JSON-over-stdio) protocol for passing diagnostics back.

The current protocol is: source on stdin, minified source on stdout,
diagnostics on stderr in "minifyjs: ..." lines. This is the same
interface a shell user gets from `cat app.js | minifyjs`.
"""

from __future__ import annotations

from typing import List

from .options import BundleOptions, Options


def build_args(opts: Options) -> List[str]:
    """Build the CLI argument list for a transform (non-bundle) call."""
    args: List[str] = []
    if opts.compress:
        args.append("--compress")
    if opts.mangle:
        args.append("--mangle")
    if opts.target:
        args.extend(["--target", opts.target])
    if opts.format:
        args.extend(["--format", opts.format])
    if opts.sourcemap:
        args.extend([f"--sourcemap={opts.sourcemap}"])
    if opts.banner:
        args.extend(["--banner", opts.banner])
    if opts.footer:
        args.extend(["--footer", opts.footer])
    if opts.legal_comments:
        args.extend(["--legal-comments", opts.legal_comments])
    return args


def build_bundle_args(opts: BundleOptions) -> List[str]:
    """Build the CLI argument list for a bundle call."""
    args: List[str] = ["--bundle"]
    args.extend(build_args(opts))
    if opts.outdir:
        args.extend(["--outdir", opts.outdir])
    if opts.outfile:
        args.extend(["--outfile", opts.outfile])
    if opts.platform:
        args.extend(["--platform", opts.platform])
    if opts.splitting:
        args.append("--splitting")
    if opts.entry_points:
        args.extend(opts.entry_points)
    return args