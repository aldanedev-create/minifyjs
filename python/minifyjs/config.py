"""Python-side loading of minifyjs.config.json.

Mirrors the Go loader in core/internal/config. Exists so Python tools
that want to inspect or override the config can do so without
shelling out.
"""

from __future__ import annotations

import json
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any, Dict, List, Optional, Union

#: Config file names, in discovery order. First match wins.
CONFIG_FILENAMES = (
    "minifyjs.config.json",
    ".minifyjsrc.json",
    ".minifyjsrc",
)


@dataclass
class Config:
    """The subset of config values the Python package cares about.

    Unknown keys are preserved in ``extra`` so a config written for a
    newer version round-trips without loss.
    """

    inputs: List[str] = field(default_factory=list)
    output: Optional[str] = None
    target: Optional[str] = None
    format: Optional[str] = None
    sourcemap: Optional[str] = None
    banner: Optional[str] = None
    footer: Optional[str] = None
    legal_comments: Optional[str] = None
    compress: bool = False
    mangle: bool = False
    extra: Dict[str, Any] = field(default_factory=dict)


def find(start: Union[str, Path] = ".") -> Optional[Path]:
    """Walk up from ``start`` looking for a config file.

    Returns the path of the first match, or None if none was found.
    """
    current = Path(start).resolve()
    while True:
        for name in CONFIG_FILENAMES:
            candidate = current / name
            if candidate.is_file():
                return candidate
        parent = current.parent
        if parent == current:
            return None
        current = parent


def load(path: Union[str, Path]) -> Config:
    """Load a specific config file."""
    p = Path(path)
    with p.open("r", encoding="utf-8") as f:
        data = json.load(f)
    return _from_dict(data)


def load_discovered(start: Union[str, Path] = ".") -> Optional[Config]:
    """Load the discovered config, or return None if none is found."""
    found = find(start)
    if found is None:
        return None
    return load(found)


def _from_dict(data: Dict[str, Any]) -> Config:
    minify_block = data.get("minify", {}) or {}
    cfg = Config(
        inputs=list(data.get("inputs", []) or []),
        output=data.get("output"),
        target=data.get("target"),
        format=data.get("format"),
        sourcemap=data.get("sourcemap"),
        banner=data.get("banner"),
        footer=data.get("footer"),
        legal_comments=data.get("legalComments"),
        compress=bool(minify_block.get("whitespace") or minify_block.get("syntax")),
        mangle=bool(minify_block.get("identifiers")),
    )
    # Preserve unknown top-level keys.
    known = {
        "inputs", "output", "target", "format", "sourcemap",
        "banner", "footer", "legalComments", "minify",
    }
    for key, value in data.items():
        if key not in known:
            cfg.extra[key] = value
    return cfg