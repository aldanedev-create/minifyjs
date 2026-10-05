#!/usr/bin/env python3
"""Generate bench/fixtures/medium.js.

A ~50 KB file with a realistic mix of declarations, functions,
classes, and control flow. Generated rather than hand-written so the
fixture can be regenerated when the schema changes.
"""

from __future__ import annotations

from pathlib import Path

HERE = Path(__file__).parent


def generate() -> str:
    parts = [
        '"use strict";',
        "",
        "// Generated benchmark fixture. Do not edit by hand.",
        "// Regenerate with: python generate_medium.py",
        "",
    ]

    # 200 utility functions.
    for i in range(200):
        parts.append(f"function util{i}(a, b) {{")
        parts.append(f"    if (a > b) {{")
        parts.append(f"        return a * {i} + b;")
        parts.append(f"    }} else if (a === b) {{")
        parts.append(f"        return 0;")
        parts.append(f"    }} else {{")
        parts.append(f"        return b - a * {i};")
        parts.append(f"    }}")
        parts.append(f"}}")
        parts.append("")

    # 50 classes with methods.
    for i in range(50):
        parts.append(f"class Widget{i} {{")
        parts.append(f"    constructor(name) {{")
        parts.append(f"        this.name = name;")
        parts.append(f"        this.value = {i};")
        parts.append(f"    }}")
        parts.append(f"    getValue() {{")
        parts.append(f"        return this.value;")
        parts.append(f"    }}")
        parts.append(f"    setValue(v) {{")
        parts.append(f"        this.value = v;")
        parts.append(f"        return this;")
        parts.append(f"    }}")
        parts.append(f"}}")
        parts.append("")

    # A list of constants.
    parts.append("const CONSTANTS = {")
    for i in range(50):
        parts.append(f"    KEY_{i}: {i * 100},")
    parts.append("};")
    parts.append("")

    # A big loop.
    parts.append("function main() {")
    parts.append("    const results = [];")
    parts.append("    for (let i = 0; i < 100; i++) {")
    parts.append("        results.push(util0(i, i + 1));")
    parts.append("    }")
    parts.append("    return results;")
    parts.append("}")
    parts.append("")

    return "\n".join(parts)


if __name__ == "__main__":
    out = HERE / "medium.js"
    out.write_text(generate(), encoding="utf-8")
    print(f"wrote {out} ({out.stat().st_size} bytes)")