"""Structured diagnostics produced by the engine.

Mirrors core/internal/diagnostics.Diagnostic. The Python side parses
the JSON emitted by the CLI's --diagnostics-json flag (or falls back
to line-by-line stderr parsing for older binary versions).
"""

from __future__ import annotations

from dataclasses import dataclass


@dataclass
class Diagnostic:
    """A single message produced by the engine."""

    severity: str
    """One of "info", "warning", "error"."""

    message: str
    """Human-readable description."""

    code: str = ""
    """Stable identifier like "lex.unterminated-string"."""

    file: str = ""
    """Source file the diagnostic refers to, if known."""

    line: int = 0
    """1-based line number, if known."""

    col: int = 0
    """1-based column number, if known."""

    @property
    def is_error(self) -> bool:
        return self.severity == "error"

    @property
    def has_location(self) -> bool:
        return self.line > 0 and self.col > 0

    def __str__(self) -> str:
        loc = ""
        if self.has_location:
            if self.file:
                loc = f"{self.file}:{self.line}:{self.col}: "
            else:
                loc = f"{self.line}:{self.col}: "
        return f"{loc}{self.severity}: {self.message}"