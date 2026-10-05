"""Locates (and sanity-checks) the native minifyjs binary bundled with
this package.

The binary is installed by the release pipeline into
``minifyjs/bin/``. This module never invokes ``go build``; if the
binary is not present, the package was either installed from source
without a build step or is running on an unsupported platform.
"""

from __future__ import annotations

import os
import stat

from ._errors import (
    binary_not_executable_message,
    binary_not_found_message,
)
from ._paths import binary_path
from ._platform import platform_tag
from .errors import BinaryNotFoundError, MinifyJSError


def find_binary() -> str:
    """Return the path to the native binary, making it executable if
    needed. Raises BinaryNotFoundError if the binary is absent."""
    path = binary_path()
    if not os.path.isfile(path):
        raise BinaryNotFoundError(binary_not_found_message(path, platform_tag()))

    # pip does not always preserve the +x bit when unpacking a wheel on
    # every platform/filesystem. Restore it if missing.
    mode = os.stat(path).st_mode
    if not mode & stat.S_IXUSR:
        try:
            os.chmod(path, mode | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)
        except OSError as e:
            raise MinifyJSError(binary_not_executable_message(path)) from e

    return path