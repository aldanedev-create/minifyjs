#!/usr/bin/env python3
"""A minimal Django management entry point.

This is the standard Django `manage.py` with one addition: a custom
management command `minifystatic` that minifies every `.js` file in
every `STATICFILES_DIRS` entry.

Run:

    pip install -r requirements.txt
    python manage.py minifystatic
"""

from __future__ import annotations

import os
import sys
from pathlib import Path


def main() -> int:
    os.environ.setdefault("DJANGO_SETTINGS_MODULE", "settings")
    try:
        from django.core.management import execute_from_command_line
    except ImportError as e:
        raise ImportError(
            "Django is not installed. Run: pip install -r requirements.txt"
        ) from e
    execute_from_command_line(sys.argv)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())