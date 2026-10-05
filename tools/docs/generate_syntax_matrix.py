#!/usr/bin/env python3
"""Generate docs/supported-syntax.md from the syntax matrix.

Runs the Go tool `tools/syntax_matrix` to produce a markdown table,
then wraps it in the documentation front matter.

Usage:

    python tools/docs/generate_syntax_matrix.py --output docs/supported-syntax.md
    python tools/docs/generate_syntax_matrix.py             # prints to stdout
"""

from __future__ import annotations

import argparse
import subprocess
import sys
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[2]


def run_syntax_matrix() -> str:
    """Shell out to `go run ./tools/syntax_matrix` and capture output."""
    proc = subprocess.run(
        ["go", "run", "./tools/syntax_matrix"],
        cwd=REPO_ROOT,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
    )
    if proc.returncode != 0:
        print(f"error: syntax_matrix exited {proc.returncode}", file=sys.stderr)
        print(proc.stderr.decode("utf-8"), file=sys.stderr)
        sys.exit(1)
    return proc.stdout.decode("utf-8")


def wrap(table: str) -> str:
    return (
        "# Supported syntax\n\n"
        "MinifyJS accepts every JavaScript feature that esbuild accepts.\n"
        "The table below lists the fixtures in\n"
        "`core/testdata/fixtures/` that exercise each feature and the\n"
        "ECMAScript version the feature was introduced in.\n\n"
        "This file is generated. Do not edit it by hand. Regenerate with:\n\n"
        "```console\n"
        "python tools/docs/generate_syntax_matrix.py --output docs/supported-syntax.md\n"
        "```\n\n"
        "## Table\n\n"
        + table.split("\n", 4)[-1]  # drop syntax_matrix's own header
    )


def main() -> int:
    parser = argparse.ArgumentParser(description="Generate supported syntax docs")
    parser.add_argument("--output", default=None)
    args = parser.parse_args()

    table = run_syntax_matrix()
    content = wrap(table)

    if args.output is None:
        sys.stdout.write(content)
        return 0

    out = Path(args.output)
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(content, encoding="utf-8")
    print(f"wrote {out}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())