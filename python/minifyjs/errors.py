"""Public exception types raised by the minifyjs package."""

from __future__ import annotations


class MinifyJSError(Exception):
    """Base class for all errors raised by minifyjs."""


class BinaryNotFoundError(MinifyJSError):
    """Raised when the precompiled native engine cannot be located."""


class MinifyError(MinifyJSError):
    """Raised when the native engine reports it could not process the
    given source (as opposed to a warning/diagnostic, which is
    returned on Result.diagnostics instead)."""


class BundleError(MinifyError):
    """Raised when a bundling build fails."""