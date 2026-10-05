"""Defines the CLI flags used to talk to the native binary from Python.

Kept as a single place to change if/when the binary grows a richer
(e.g. JSON-over-stdio) protocol for passing diagnostics back.

The current protocol is: source on stdin, minified source on stdout,
diagnostics on stderr in "minifyjs: ..." lines. This is the same
interface a shell user gets from `cat app.js | minifyjs`.
"""

from __future__ import annotations

from .options import BundleOptions, Options


def build_args(opts: Options) -> list[str]:
    """Build the CLI argument list for a transform (non-bundle) call."""
    args: list[str] = []
    if opts.compress:
        args.append("--compress")
    elif not opts.mangle:
        args.append("--no-minify")
    else:
        args.append("--no-compress")
    if opts.mangle:
        args.append("--mangle")
    else:
        args.append("--no-mangle")
    if opts.target:
        args.extend(["--target", opts.target])
    if opts.format:
        args.extend(["--format", opts.format])
    if opts.sourcemap:
        mode = "inline" if opts.sourcemap == "external" else opts.sourcemap
        args.extend([f"--sourcemap={mode}"])
    if opts.banner:
        args.extend(["--banner", opts.banner])
    if opts.footer:
        args.extend(["--footer", opts.footer])
    if opts.legal_comments:
        args.extend(["--legal-comments", opts.legal_comments])
    return args


def build_bundle_args(opts: BundleOptions) -> list[str]:
    """Build the CLI argument list for a bundle call."""
    args: list[str] = ["--bundle"]
    args.extend(build_args(opts))
    if opts.sourcemap == "external":
        args[args.index("--sourcemap=inline")] = "--sourcemap=external"
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
