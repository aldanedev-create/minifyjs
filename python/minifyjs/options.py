"""Options accepted by minify(), optimize(), and bundle().

Mirrors the Go API in core/api.Options so the two surfaces stay in
sync. Adding a field here means adding it on the Go side and vice
versa.
"""

from __future__ import annotations

from dataclasses import dataclass


@dataclass
class Options:
    """Options for minify() and optimize().

    All fields have safe defaults. ``Options()`` is a no-op transform
    (modulo esbuild's normalizing re-print).
    """

    #: Enable whitespace removal.
    compress: bool = False

    #: Enable identifier mangling.
    mangle: bool = False

    #: ECMAScript target for output ("es2015", "esnext", "es2020,chrome90").
    #: Empty means no lowering.
    target: str | None = None

    #: Module format ("iife", "cjs", "esm"). Empty means "preserve".
    format: str | None = None

    #: Source map mode ("inline", "external", "both"). Empty means none.
    sourcemap: str | None = None

    #: Text prepended to the output. Useful for license headers.
    banner: str | None = None

    #: Text appended to the output.
    footer: str | None = None

    #: Legal comments mode ("none", "inline", "eof", "external").
    legal_comments: str | None = None

    #: Source file name shown in diagnostics. Empty means "<stdin>".
    source_name: str | None = None


@dataclass
class BundleOptions(Options):
    """Options for bundle().

    Extends Options with the fields specific to multi-file builds.
    """

    #: Entry point files to bundle. At least one is required.
    entry_points: list[str] | None = None

    #: Output directory. Either outdir or outfile must be set.
    outdir: str | None = None

    #: Output file. Cannot be set together with outdir.
    outfile: str | None = None

    #: Directory all relative paths are resolved against.
    working_dir: str | None = None

    #: Target platform ("browser", "node", "neutral").
    platform: str = "browser"

    #: Enable code splitting. Requires format="esm".
    splitting: bool = False
    entry_names: str | None = None
    chunk_names: str | None = None
    asset_names: str | None = None
    metafile: str | None = None
    external: list[str] | None = None
    packages: str = "bundle"
    tree_shaking: bool | None = None
    charset: str = "utf8"
    define: dict[str, str] | None = None
    drop: list[str] | None = None
