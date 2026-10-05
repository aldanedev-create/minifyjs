"""Runs the native binary as a subprocess, feeding it source on stdin
and reading generated code from stdout.

This is the only module in the package that spawns processes. Every
public API function (minify, optimize, bundle) ultimately funnels
through one of the two functions here.
"""

from __future__ import annotations

import subprocess
from typing import List

from ._binary import find_binary
from ._protocol import build_args, build_bundle_args
from .diagnostics import Diagnostic
from .errors import BundleError, MinifyError
from .options import BundleOptions, Options
from .result import Result

_DIAG_PREFIX = "minifyjs: "


def run_minify(source: str, opts: Options) -> Result:
    """Run a transform (non-bundle) call through the native binary."""
    binary = find_binary()
    args = [binary, *build_args(opts)]

    proc = subprocess.run(
        args,
        input=source.encode("utf-8"),
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
    )

    diagnostics = _parse_diagnostics(proc.stderr.decode("utf-8", "replace"))

    if proc.returncode != 0:
        raise MinifyError(
            _format_error_message(proc.returncode, diagnostics)
        )

    code = proc.stdout.decode("utf-8")
    return Result(
        code=code,
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
    binary = find_binary()
    args = [binary, *build_bundle_args(opts)]

    proc = subprocess.run(
        args,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
    )

    diagnostics = _parse_diagnostics(proc.stderr.decode("utf-8", "replace"))

    if proc.returncode != 0:
        raise BundleError(
            _format_error_message(proc.returncode, diagnostics)
        )

    code = proc.stdout.decode("utf-8")
    return Result(
        code=code,
        original_bytes=0,
        minified_bytes=len(code.encode("utf-8")),
        diagnostics=diagnostics,
    )


def _parse_diagnostics(stderr: str) -> List[Diagnostic]:
    """Parse the CLI's stderr into Diagnostic objects.

    The format is one diagnostic per line:

        minifyjs: file:line:col: severity: message
        minifyjs: severity: message

    Lines that do not match either shape are ignored.
    """
    out: List[Diagnostic] = []
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


def _format_error_message(returncode: int, diags: List[Diagnostic]) -> str:
    errs = [d for d in diags if d.is_error]
    if errs:
        return "; ".join(str(d) for d in errs)
    if diags:
        return "; ".join(str(d) for d in diags)
    return f"minifyjs exited with status {returncode}"