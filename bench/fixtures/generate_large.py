#!/usr/bin/env python3
"""Generate bench/fixtures/large.js.

A ~500 KB file. Regenerate with:
    python generate_large.py
"""

from __future__ import annotations

from pathlib import Path

HERE = Path(__file__).parent


def generate() -> str:
    parts = [
        '"use strict";',
        "",
        "// Generated benchmark fixture. Do not edit by hand.",
        "",
    ]

    for i in range(2000):
        parts.append(f"function fn{i}(x) {{")
        parts.append(f"    const doubled = x * 2;")
        parts.append(f"    const tripled = x * 3;")
        parts.append(f"    if (doubled > tripled) {{")
        parts.append(f"        return doubled - {i};")
        parts.append(f"    }}")
        parts.append(f"    return tripled + {i};")
        parts.append(f"}}")
        parts.append("")

    for i in range(300):
        parts.append(f"const obj{i} = {{")
        parts.append(f"    a: {i},")
        parts.append(f"    b: {i * 2},")
        parts.append(f"    c: {i * 3},")
        parts.append(f"}};")
        parts.append("")

    return "\n".join(parts)


if __name__ == "__main__":
    out = HERE / "large.js"
    out.write_text(generate(), encoding="utf-8")
    print(f"wrote {out} ({out.stat().st_size} bytes)")