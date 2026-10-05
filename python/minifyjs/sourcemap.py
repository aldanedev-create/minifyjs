"""Source map handling for the Python package.

MinifyJS delegates source-map generation to esbuild, which produces
standard V3 source maps. This module provides Python-side helpers for
reading, inspecting, and lightly manipulating those maps.

It does not re-implement the V3 format. It exists so callers who want
to check "does this map reference this file?" or "what are the sources
in this map?" can do so without depending on a third-party source-map
library.

The public surface is intentionally small:

    SourceMap           -- a parsed source map
    SourceMap.load()    -- read a .map file from disk
    SourceMap.parse()   -- parse a .map JSON string

    decode_mappings()   -- decode the base64 VLQ "mappings" field
    encode_mappings()   -- the inverse

The encode/decode functions are used by tests that compare a
generated map against a known-good reference. Application code
should not need them.
"""

from __future__ import annotations

import json
from dataclasses import dataclass, field
from pathlib import Path
from typing import Dict, List, Optional, Tuple, Union

PathLike = Union[str, Path]


# ---------------------------------------------------------------------------
# Source map data model
# ---------------------------------------------------------------------------


@dataclass
class SourceMap:
    """A parsed source map.

    Only the fields MinifyJS's tooling needs are modelled. The raw
    JSON is preserved in ``raw`` so unknown fields survive a
    round-trip through :meth:`to_json`.
    """

    version: int = 3
    file: str = ""
    source_root: str = ""
    sources: List[str] = field(default_factory=list)
    sources_content: List[Optional[str]] = field(default_factory=list)
    names: List[str] = field(default_factory=list)
    mappings: str = ""

    #: The original parsed JSON, with every field the loader saw.
    raw: Dict = field(default_factory=dict)

    # ------------------------------------------------------------------
    # Construction
    # ------------------------------------------------------------------

    @classmethod
    def parse(cls, text: str) -> "SourceMap":
        """Parse a source map from a JSON string."""
        data = json.loads(text)
        if not isinstance(data, dict):
            raise ValueError("source map must be a JSON object")
        return cls._from_dict(data)

    @classmethod
    def load(cls, path: PathLike) -> "SourceMap":
        """Load a source map from a file on disk."""
        text = Path(path).read_text(encoding="utf-8")
        return cls.parse(text)

    @classmethod
    def _from_dict(cls, data: Dict) -> "SourceMap":
        return cls(
            version=int(data.get("version", 3)),
            file=data.get("file", ""),
            source_root=data.get("sourceRoot", ""),
            sources=list(data.get("sources", [])),
            sources_content=list(data.get("sourcesContent", []) or []),
            names=list(data.get("names", [])),
            mappings=data.get("mappings", ""),
            raw=dict(data),
        )

    # ------------------------------------------------------------------
    # Serialization
    # ------------------------------------------------------------------

    def to_dict(self) -> Dict:
        """Return a JSON-serializable dict for this map.

        Unknown fields seen at parse time are preserved.
        """
        out = dict(self.raw)
        out["version"] = self.version
        if self.file:
            out["file"] = self.file
        if self.source_root:
            out["sourceRoot"] = self.source_root
        out["sources"] = list(self.sources)
        if self.sources_content:
            out["sourcesContent"] = list(self.sources_content)
        if self.names:
            out["names"] = list(self.names)
        out["mappings"] = self.mappings
        return out

    def to_json(self, indent: Optional[int] = None) -> str:
        """Serialize back to JSON. ``indent=None`` produces compact output."""
        return json.dumps(self.to_dict(), indent=indent, ensure_ascii=False)

    # ------------------------------------------------------------------
    # Queries
    # ------------------------------------------------------------------

    def has_source(self, name: str) -> bool:
        """Return True if ``name`` appears in the sources list."""
        return name in self.sources

    def source_index(self, name: str) -> int:
        """Return the index of ``name`` in sources, or -1 if absent."""
        try:
            return self.sources.index(name)
        except ValueError:
            return -1

    def name_index(self, name: str) -> int:
        """Return the index of ``name`` in names, or -1 if absent."""
        try:
            return self.names.index(name)
        except ValueError:
            return -1

    @property
    def source_count(self) -> int:
        """Number of distinct source files referenced by this map."""
        return len(self.sources)

    @property
    def has_content(self) -> bool:
        """True if the map embeds the original sources."""
        return bool(self.sources_content) and all(
            c is not None for c in self.sources_content
        )

    # ------------------------------------------------------------------
    # Mapping decoding
    # ------------------------------------------------------------------

    def decode(self) -> List[List[Tuple[int, int, int, int, int]]]:
        """Decode the ``mappings`` field into per-line segment lists.

        Each segment is a 5-tuple of integers:
        ``(gen_col, src_idx, src_line, src_col, name_idx)``.

        Segments with only 1 or 2 fields are represented with -1 for
        the missing fields. This mirrors the V3 spec's rule that the
        source/name fields are optional.
        """
        return decode_mappings(self.mappings)


# ---------------------------------------------------------------------------
# VLQ base64 codec
# ---------------------------------------------------------------------------

# Base64 alphabet used by the V3 source map format.
_B64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

# Reverse lookup for decoding.
_B64_INDEX: Dict[str, int] = {c: i for i, c in enumerate(_B64)}


def _decode_vlq(s: str) -> List[int]:
    """Decode a base64 VLQ string into a list of signed integers."""
    values: List[int] = []
    shift = 0
    value = 0
    for ch in s:
        digit = _B64_INDEX.get(ch)
        if digit is None:
            raise ValueError(f"invalid base64 VLQ character: {ch!r}")
        continuation = digit & 0x20
        digit &= 0x1F
        value += digit << shift
        if continuation:
            shift += 5
        else:
            # The low bit is the sign; shift right and negate if set.
            if value & 1:
                values.append(-(value >> 1))
            else:
                values.append(value >> 1)
            value = 0
            shift = 0
    if shift != 0:
        raise ValueError("truncated VLQ sequence")
    return values


def _encode_vlq(values: List[int]) -> str:
    """Encode a list of signed integers into base64 VLQ."""
    out: List[str] = []
    for v in values:
        # Move the sign to the low bit.
        vlq = (-v << 1) | 1 if v < 0 else v << 1
        while True:
            digit = vlq & 0x1F
            vlq >>= 5
            if vlq:
                digit |= 0x20  # continuation bit
            out.append(_B64[digit])
            if not vlq:
                break
    return "".join(out)


# ---------------------------------------------------------------------------
# Mappings encoder / decoder
# ---------------------------------------------------------------------------


def decode_mappings(
    mappings: str,
) -> List[List[Tuple[int, int, int, int, int]]]:
    """Decode the V3 ``mappings`` field into per-line segment lists.

    Returns a list whose length equals the number of generated lines.
    Each element is a list of segments for that line. A segment is a
    tuple ``(gen_col, src_idx, src_line, src_col, name_idx)``. When a
    field is absent in the encoding, the tuple uses -1 for it.

    The V3 spec delta-encodes every field except the first segment of
    a line, whose ``gen_col`` is absolute.
    """
    if not mappings:
        return []

    lines: List[List[Tuple[int, int, int, int, int]]] = []

    # Running state across the whole mappings string.
    src_idx = 0
    src_line = 0
    src_col = 0
    name_idx = 0

    for line_str in mappings.split(";"):
        line_segments: List[Tuple[int, int, int, int, int]] = []
        gen_col = 0

        for segment_str in line_str.split(","):
            if not segment_str:
                continue
            values = _decode_vlq(segment_str)
            if not 1 <= len(values) <= 5:
                raise ValueError(
                    f"segment must have 1..5 fields, got {len(values)}"
                )

            gen_col += values[0]

            if len(values) >= 4:
                src_idx += values[1]
                src_line += values[2]
                src_col += values[3]
                if len(values) == 5:
                    name_idx += values[4]
                    line_segments.append(
                        (gen_col, src_idx, src_line, src_col, name_idx)
                    )
                else:
                    line_segments.append(
                        (gen_col, src_idx, src_line, src_col, -1)
                    )
            else:
                # Segment has only gen_col (and optionally a dummy
                # source field). Represent missing fields as -1.
                line_segments.append((gen_col, -1, -1, -1, -1))

        lines.append(line_segments)

    return lines


def encode_mappings(
    lines: List[List[Tuple[int, int, int, int, int]]],
) -> str:
    """Encode per-line segment lists into a V3 ``mappings`` string.

    This is the inverse of :func:`decode_mappings`. It is used by
    tests that build a map by hand.
    """
    src_idx = 0
    src_line = 0
    src_col = 0
    name_idx = 0

    out_lines: List[str] = []

    for line_segments in lines:
        gen_col = 0
        out_segments: List[str] = []
        for seg in line_segments:
            gen_col_delta = seg[0] - gen_col
            gen_col = seg[0]

            if seg[1] < 0:
                # Only a generated column.
                out_segments.append(_encode_vlq([gen_col_delta]))
                continue

            src_idx_delta = seg[1] - src_idx
            src_idx = seg[1]
            src_line_delta = seg[2] - src_line
            src_line = seg[2]
            src_col_delta = seg[3] - src_col
            src_col = seg[3]

            if seg[4] < 0:
                out_segments.append(
                    _encode_vlq([
                        gen_col_delta,
                        src_idx_delta,
                        src_line_delta,
                        src_col_delta,
                    ])
                )
            else:
                name_idx_delta = seg[4] - name_idx
                name_idx = seg[4]
                out_segments.append(
                    _encode_vlq([
                        gen_col_delta,
                        src_idx_delta,
                        src_line_delta,
                        src_col_delta,
                        name_idx_delta,
                    ])
                )

        out_lines.append(",".join(out_segments))

    return ";".join(out_lines)


# ---------------------------------------------------------------------------
# Inline source map helpers
# ---------------------------------------------------------------------------


def extract_inline_map(js_source: str) -> Optional[str]:
    """If ``js_source`` ends with an inline source map, return it as
    a decoded JSON string. Otherwise return None.

    Inline maps are data URLs:

        //# sourceMappingURL=data:application/json;base64,<b64>

    This function strips the marker, decodes the base64 payload, and
    returns the JSON text. It does not parse the JSON.
    """
    import base64
    import re

    marker = "sourceMappingURL=data:application/json"
    idx = js_source.rfind(marker)
    if idx < 0:
        return None

    rest = js_source[idx + len(marker):]
    # The remainder should start with either ";base64,<data>" or
    # ",<url-encoded data>". Only base64 is emitted by esbuild.
    if rest.startswith(";base64,"):
        payload = rest[len(";base64,"):]
        # Strip any trailing whitespace/newline.
        payload = payload.splitlines()[0].strip()
        try:
            return base64.b64decode(payload).decode("utf-8")
        except (ValueError, UnicodeDecodeError):
            return None

    return None


def strip_inline_map(js_source: str) -> str:
    """Return ``js_source`` with any trailing inline source map removed.

    The data URL is a comment in the source. This function removes
    the whole comment, including the ``//#`` prefix and any trailing
    newline, so the returned string is exactly the minified code.
    """
    marker = "//# sourceMappingURL=data:"
    idx = js_source.rfind(marker)
    if idx < 0:
        return js_source
    # Trim any preceding newline as well.
    cut = idx
    while cut > 0 and js_source[cut - 1] in ("\n", "\r"):
        cut -= 1
    return js_source[:cut]