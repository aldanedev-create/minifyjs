#!/usr/bin/env python3
"""Run the full benchmark suite and write results/latest.md.

    python run_all.py

Skips any tool that is not installed. Always runs MinifyJS.
"""

from __future__ import annotations

import io
import shutil
import subprocess
import sys
from contextlib import redirect_stdout
from datetime import datetime, timezone
from pathlib import Path
import os

# Put bench/node_modules/.bin on PATH so the local install of
# esbuild / terser / uglify-js is found without a global npm install.
_BENCH_DIR = Path(__file__).parent
_NODE_BIN = _BENCH_DIR / "node_modules" / ".bin"
if _NODE_BIN.is_dir():
    os.environ["PATH"] = str(_NODE_BIN) + os.pathsep + os.environ.get("PATH", "")

HERE = Path(__file__).parent
RESULTS = HERE / "results"


def capture(script: str, *args: str) -> str:
    buf = io.StringIO()
    with redirect_stdout(buf):
        subprocess.run(
            [sys.executable, str(HERE / script), *args],
            check=False,
        )
    return buf.getvalue()


def main() -> int:
    RESULTS.mkdir(exist_ok=True)

    sections = []

    sections.append("# MinifyJS benchmark results\n")
    sections.append(
        f"Generated: {datetime.now(timezone.utc).strftime('%Y-%m-%d %H:%M UTC')}\n"
    )

    sections.append("\n## Compression\n\n")
    sections.append(capture("measure_compression.py"))

    sections.append("\n## Throughput (medium.js)\n\n")
    sections.append(capture("measure_throughput.py", "--fixture", "medium.js"))

    sections.append("\n## Startup time\n\n")
    sections.append(capture("measure_startup.py"))

    sections.append("\n## Peak memory (large.js)\n\n")
    sections.append(capture("measure_memory.py"))

    sections.append(
        "\n## Methodology\n\n"
        "See `bench/README.md`. Each measurement is the median of "
        "multiple runs in isolated subprocesses.\n"
    )

    out = RESULTS / "latest.md"
    out.write_text("".join(sections), encoding="utf-8")
    print(f"wrote {out}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())