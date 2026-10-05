"""Runs the native binary as a subprocess, feeding it source on stdin
and reading generated code from stdout.

This is the only module in the package that spawns processes. Every
public API function (minify, optimize, bundle) ultimately funnels
through one of the two functions here.
"""

from __future__ import annotations

import contextlib
import json
import subprocess
import tempfile
from dataclasses import replace
from pathlib import Path

from ._binary import find_binary
from ._protocol import build_args, build_bundle_args
from .diagnostics import Diagnostic
from .errors import BundleError, MinifyError
from .options import BundleOptions, Options
from .result import Result
from .sourcemap import extract_inline_map, strip_inline_map

_DIAG_PREFIX = "minifyjs: "


def run_minify(source: str, opts: Options) -> Result:
    """Run a transform (non-bundle) call through the native binary."""
    binary = find_binary()
    args = [binary, *build_args(opts)]

    proc = subprocess.run(
        args,
        input=source.encode("utf-8"),
        capture_output=True,
        check=False,
    )

    diagnostics = _parse_diagnostics(proc.stderr.decode("utf-8", "replace"))
    if opts.source_name:
        for diagnostic in diagnostics:
            if diagnostic.file == "<stdin>":
                diagnostic.file = opts.source_name

    if proc.returncode != 0:
        raise MinifyError(
            _format_error_message(proc.returncode, diagnostics)
        )

    code = proc.stdout.decode("utf-8")
    map_text = ""
    if opts.sourcemap == "external":
        map_text = extract_inline_map(code) or ""
        code = strip_inline_map(code)
        if map_text:
            with contextlib.suppress(json.JSONDecodeError):
                map_text = json.dumps(json.loads(map_text), separators=(",", ":"))
    return Result(
        code=code,
        map=map_text,
        original_bytes=len(source.encode("utf-8")),
        minified_bytes=len(code.encode("utf-8")),
        diagnostics=diagnostics,
    )


def run_bundle(opts: BundleOptions) -> Result:
    """Run a bundle call through the native binary.

    The native binary writes bundle output to disk (there is no way
    to stream a multi-file bundle through stdout). The returned
    Result carries diagnostics and, when outfile was used, the code.
    """
    working_dir = Path(opts.working_dir or Path.cwd()).resolve()
    with tempfile.TemporaryDirectory(prefix="minifyjs-metadata-") as temp:
        metadata_path = Path(opts.metafile) if opts.metafile else Path(temp) / "meta.json"
        if not metadata_path.is_absolute():
            metadata_path = working_dir / metadata_path
        effective = replace(opts, working_dir=str(working_dir), metafile=str(metadata_path))
        proc = subprocess.run([find_binary(), *build_bundle_args(effective)], capture_output=True, check=False)
        diagnostics = _parse_diagnostics(proc.stderr.decode("utf-8", "replace"))
        if proc.returncode != 0:
            raise BundleError(_format_error_message(proc.returncode, diagnostics))
        metadata = json.loads(metadata_path.read_text(encoding="utf-8"))
        outputs = [{"path": str((working_dir / path).resolve()), "bytes": info["bytes"]}
                   for path, info in metadata.get("outputs", {}).items()]
        code = ""
        map_text = ""
        if opts.outfile:
            output = (working_dir / opts.outfile).resolve()
            code = output.read_text(encoding="utf-8")
            map_path = Path(str(output) + ".map")
            if opts.sourcemap in ("external", "both") and map_path.is_file():
                map_text = map_path.read_text(encoding="utf-8")
        return Result(code=code, map=map_text, original_bytes=0,
                      minified_bytes=sum(f["bytes"] for f in outputs),
                      diagnostics=diagnostics, output_files=outputs, metafile=metadata)


def _parse_diagnostics(stderr: str) -> list[Diagnostic]:
    """Parse the CLI's stderr into Diagnostic objects.

    The format is one diagnostic per line:

        minifyjs: file:line:col: severity: message
        minifyjs: severity: message

    Lines that do not match either shape are ignored.
    """
    out: list[Diagnostic] = []
    for line in stderr.splitlines():
        line = line.rstrip()
        if not line:
            continue
        if not line.startswith(_DIAG_PREFIX):
            continue
        body = line[len(_DIAG_PREFIX):]
        d = _parse_one(body)
        if d is not None:
            out.append(d)
    return out


def _parse_one(body: str) -> Diagnostic | None:
    """Parse a single diagnostic line body (after the "minifyjs: "
    prefix) into a Diagnostic, or return None if it does not look
    like a diagnostic."""
    # Try "file:line:col: severity: message" first.
    parts = body.split(":", 4)
    if len(parts) >= 4:
        # Heuristic: parts[1] and parts[2] must be integers.
        try:
            line = int(parts[1])
            col = int(parts[2])
        except ValueError:
            pass
        else:
            sev = parts[3].strip()
            msg = parts[4].strip() if len(parts) > 4 else ""
            if sev in ("info", "warning", "error"):
                return Diagnostic(
                    severity=sev,
                    message=msg,
                    file=parts[0].strip(),
                    line=line,
                    col=col,
                )

    # Fall back to "severity: message".
    parts = body.split(":", 1)
    if len(parts) == 2:
        sev = parts[0].strip()
        if sev in ("info", "warning", "error"):
            return Diagnostic(severity=sev, message=parts[1].strip())

    return None


def _format_error_message(returncode: int, diags: list[Diagnostic]) -> str:
    errs = [d for d in diags if d.is_error]
    if errs:
        return "; ".join(str(d) for d in errs)
    if diags:
        return "; ".join(str(d) for d in diags)
    return f"minifyjs exited with status {returncode}"
