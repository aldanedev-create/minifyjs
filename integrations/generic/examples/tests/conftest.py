"""pytest fixtures for the integration example tests.

Tests in this package do not require Django, Flask, or FastAPI to be
installed. The examples only use file paths, so tests exercise them
against real temp directories.
"""

from __future__ import annotations

import sys
from pathlib import Path

# Put the repository's python/ directory on sys.path so
# `from minifyjs import Adapter` works without an install.
_REPO_ROOT = Path(__file__).resolve().parents[4]
_PYTHON_DIR = _REPO_ROOT / "python"
if str(_PYTHON_DIR) not in sys.path:
    sys.path.insert(0, str(_PYTHON_DIR))

# Also expose the examples package so tests can import it directly.
_EXAMPLES_ROOT = Path(__file__).resolve().parents[1]
if str(_EXAMPLES_ROOT) not in sys.path:
    sys.path.insert(0, str(_EXAMPLES_ROOT))