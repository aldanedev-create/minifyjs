"""pytest configuration for the minifyjs test suite."""

from __future__ import annotations

import sys
from pathlib import Path

# Ensure the source tree is importable when running pytest from
# python/ without installing the package first.
sys.path.insert(0, str(Path(__file__).parent.parent))

# Path to the bundled binary. Tests that need it skip if absent.
BINARY_PATH = Path(__file__).parent.parent / "minifyjs" / "bin" / "minifyjs"


def pytest_configure(config):
    config.addinivalue_line(
        "markers",
        "binary: test requires the bundled native binary",
    )
