"""The value minify()/optimize()/bundle() return."""

from __future__ import annotations

from dataclasses import dataclass, field

from .diagnostics import Diagnostic


@dataclass
class Result:
    """The result of a minification or bundling operation."""

    #: The minified JavaScript source.
    code: str

    #: The source map, if one was requested. Empty otherwise.
    map: str = ""

    #: The byte length of the original input.
    original_bytes: int = 0

    #: The byte length of the minified output.
    minified_bytes: int = 0

    #: Any diagnostics produced by the engine.
    diagnostics: list[Diagnostic] = field(default_factory=list)

    @property
    def ratio(self) -> float:
        """Compression ratio as minified/original (0..1); 0 if empty."""
        if self.original_bytes == 0:
            return 0.0
        return self.minified_bytes / self.original_bytes

    @property
    def bytes_saved(self) -> int:
        """original_bytes - minified_bytes. Negative if output grew."""
        return self.original_bytes - self.minified_bytes

    @property
    def has_errors(self) -> bool:
        """True if any diagnostic is an error."""
        return any(d.is_error for d in self.diagnostics)

    @property
    def has_warnings(self) -> bool:
        """True if any diagnostic is a warning."""
        return any(d.severity == "warning" for d in self.diagnostics)
