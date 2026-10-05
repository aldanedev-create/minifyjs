"""Single source of truth for the Python package version.

Kept separate from pyproject.toml so ``minifyjs.__version__`` works
from a source checkout that has not been pip-installed.
"""

__version__ = "0.1.1"
