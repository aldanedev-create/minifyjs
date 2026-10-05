"""Tests for the sourcemap module.

These exercise the V3 source map parsing, VLQ codec, inline map
extraction, and the small query helpers. They do not test esbuild's
map generation — that is covered by the Go integration tests.
"""

from __future__ import annotations

import base64
import json

import pytest

from minifyjs.sourcemap import (
    SourceMap,
    decode_mappings,
    encode_mappings,
    extract_inline_map,
    strip_inline_map,
)
from minifyjs.sourcemap import _decode_vlq, _encode_vlq


# ---------------------------------------------------------------------------
# VLQ codec
# ---------------------------------------------------------------------------


def test_vlq_encode_decode_zero():
    assert _decode_vlq(_encode_vlq([0])) == [0]


def test_vlq_encode_decode_positive():
    for n in (1, 2, 15, 16, 31, 32, 100, 1000, 100000):
        assert _decode_vlq(_encode_vlq([n])) == [n], f"failed for {n}"


def test_vlq_encode_decode_negative():
    for n in (-1, -2, -15, -16, -31, -32, -100, -1000, -100000):
        assert _decode_vlq(_encode_vlq([n])) == [n], f"failed for {n}"


def test_vlq_encode_decode_multiple():
    values = [0, 1, -1, 16, -16, 100, -100, 1000]
    assert _decode_vlq(_encode_vlq(values)) == values


def test_vlq_known_value_zero():
    # Per the spec, 0 encodes to "A".
    assert _encode_vlq([0]) == "A"


def test_vlq_known_value_one():
    # 1 encodes to "C".
    assert _encode_vlq([1]) == "C"


def test_vlq_known_value_negative_one():
    # -1 encodes to "D".
    assert _encode_vlq([-1]) == "D"


def test_vlq_decode_rejects_invalid_character():
    with pytest.raises(ValueError):
        _decode_vlq("!@#")


def test_vlq_decode_empty_string():
    assert _decode_vlq("") == []


def test_vlq_decode_truncated_continuation():
    # "g" has the continuation bit set but is the last character.
    with pytest.raises(ValueError):
        _decode_vlq("g")


# ---------------------------------------------------------------------------
# Mappings decode / encode
# ---------------------------------------------------------------------------


def test_decode_mappings_empty():
    assert decode_mappings("") == []


def test_decode_mappings_single_segment():
    # "AAAA" decodes to (0, 0, 0, 0).
    result = decode_mappings("AAAA")
    assert result == [[(0, 0, 0, 0, -1)]]


def test_decode_mappings_two_lines():
    result = decode_mappings("AAAA;AACA")
    assert len(result) == 2
    assert result[0] == [(0, 0, 0, 0, -1)]
    assert result[1] == [(0, 0, 1, 0, -1)]


def test_decode_mappings_segment_with_name():
    # 5-field segment: gen_col, src_idx, src_line, src_col, name_idx.
    # "AAAAA" -> (0, 0, 0, 0, 0).
    result = decode_mappings("AAAAA")
    assert result == [[(0, 0, 0, 0, 0)]]


def test_decode_mappings_multiple_segments_same_line():
    # Two segments on one line: "AAAA,SAAS".
    # First: (0,0,0,0). Second: gen_col += 9, src_col += 9.
    result = decode_mappings("AAAA,SAAS")
    assert len(result) == 1
    assert len(result[0]) == 2
    assert result[0][0] == (0, 0, 0, 0, -1)
    assert result[0][1] == (9, 0, 0, 9, -1)


def test_encode_mappings_empty():
    assert encode_mappings([]) == ""


def test_encode_mappings_single_line_single_segment():
    encoded = encode_mappings([[(0, 0, 0, 0, -1)]])
    assert encoded == "AAAA"


def test_encode_mappings_two_lines():
    encoded = encode_mappings([
        [(0, 0, 0, 0, -1)],
        [(0, 0, 1, 0, -1)],
    ])
    assert encoded == "AAAA;AACA"


def test_encode_decode_roundtrip():
    lines = [
        [(0, 0, 0, 0, -1), (5, 0, 0, 5, -1)],
        [(0, 0, 1, 0, -1)],
        [(2, 1, 0, 0, 0), (10, 1, 0, 5, 0)],
    ]
    encoded = encode_mappings(lines)
    decoded = decode_mappings(encoded)
    assert decoded == lines


def test_encode_decode_roundtrip_with_names():
    lines = [[(0, 0, 0, 0, 0), (5, 0, 0, 2, 1)]]
    encoded = encode_mappings(lines)
    decoded = decode_mappings(encoded)
    assert decoded == lines


def test_decode_mappings_rejects_oversized_segment():
    # A 6-field segment is invalid.
    with pytest.raises(ValueError):
        decode_mappings("AAAAAA")


# ---------------------------------------------------------------------------
# SourceMap.parse
# ---------------------------------------------------------------------------


def _minimal_map() -> dict:
    return {
        "version": 3,
        "file": "out.js",
        "sources": ["in.js"],
        "names": [],
        "mappings": "AAAA",
    }


def test_sourcemap_parse_minimal():
    m = SourceMap.parse(json.dumps(_minimal_map()))
    assert m.version == 3
    assert m.file == "out.js"
    assert m.sources == ["in.js"]
    assert m.mappings == "AAAA"


def test_sourcemap_parse_empty_object():
    m = SourceMap.parse("{}")
    assert m.version == 3
    assert m.sources == []
    assert m.mappings == ""


def test_sourcemap_parse_rejects_non_object():
    with pytest.raises(ValueError):
        SourceMap.parse("[]")
    with pytest.raises(ValueError):
        SourceMap.parse('"a string"')


def test_sourcemap_parse_rejects_invalid_json():
    with pytest.raises(json.JSONDecodeError):
        SourceMap.parse("not json")


def test_sourcemap_parse_preserves_unknown_fields():
    data = _minimal_map()
    data["x_google_ignoreList"] = [0]
    m = SourceMap.parse(json.dumps(data))
    assert m.raw["x_google_ignoreList"] == [0]


def test_sourcemap_parse_with_sources_content():
    data = _minimal_map()
    data["sourcesContent"] = ["const x = 1;"]
    m = SourceMap.parse(json.dumps(data))
    assert m.sources_content == ["const x = 1;"]
    assert m.has_content is True


def test_sourcemap_parse_source_root():
    data = _minimal_map()
    data["sourceRoot"] = "/src"
    m = SourceMap.parse(json.dumps(data))
    assert m.source_root == "/src"


def test_sourcemap_parse_version_string_parses_to_int():
    data = _minimal_map()
    data["version"] = "3"
    m = SourceMap.parse(json.dumps(data))
    assert m.version == 3


# ---------------------------------------------------------------------------
# SourceMap.load / to_json
# ---------------------------------------------------------------------------


def test_sourcemap_load_from_file(tmp_path):
    p = tmp_path / "out.js.map"
    p.write_text(json.dumps(_minimal_map()))
    m = SourceMap.load(p)
    assert m.sources == ["in.js"]


def test_sourcemap_to_json_roundtrip():
    m = SourceMap.parse(json.dumps(_minimal_map()))
    text = m.to_json()
    m2 = SourceMap.parse(text)
    assert m2.file == m.file
    assert m2.sources == m.sources
    assert m2.mappings == m.mappings


def test_sourcemap_to_json_pretty():
    m = SourceMap.parse(json.dumps(_minimal_map()))
    text = m.to_json(indent=2)
    assert "\n" in text


def test_sourcemap_to_json_preserves_unknown_fields():
    data = _minimal_map()
    data["x_custom"] = "value"
    m = SourceMap.parse(json.dumps(data))
    out = json.loads(m.to_json())
    assert out["x_custom"] == "value"


def test_sourcemap_to_dict_keys():
    m = SourceMap.parse(json.dumps(_minimal_map()))
    d = m.to_dict()
    assert d["version"] == 3
    assert d["file"] == "out.js"
    assert d["sources"] == ["in.js"]
    assert d["mappings"] == "AAAA"


# ---------------------------------------------------------------------------
# Query helpers
# ---------------------------------------------------------------------------


def test_sourcemap_has_source():
    m = SourceMap.parse(json.dumps({
        "version": 3,
        "sources": ["a.js", "b.js"],
        "mappings": "",
    }))
    assert m.has_source("a.js")
    assert m.has_source("b.js")
    assert not m.has_source("c.js")


def test_sourcemap_source_index():
    m = SourceMap.parse(json.dumps({
        "version": 3,
        "sources": ["a.js", "b.js"],
        "mappings": "",
    }))
    assert m.source_index("a.js") == 0
    assert m.source_index("b.js") == 1
    assert m.source_index("c.js") == -1


def test_sourcemap_name_index():
    m = SourceMap.parse(json.dumps({
        "version": 3,
        "names": ["foo", "bar"],
        "mappings": "",
    }))
    assert m.name_index("foo") == 0
    assert m.name_index("bar") == 1
    assert m.name_index("baz") == -1


def test_sourcemap_source_count():
    m = SourceMap.parse(json.dumps({
        "version": 3,
        "sources": ["a.js", "b.js", "c.js"],
        "mappings": "",
    }))
    assert m.source_count == 3


def test_sourcemap_has_content_false():
    m = SourceMap.parse(json.dumps(_minimal_map()))
    assert m.has_content is False


def test_sourcemap_has_content_partial():
    # Only some sources have content.
    m = SourceMap.parse(json.dumps({
        "version": 3,
        "sources": ["a.js", "b.js"],
        "sourcesContent": ["const x = 1;", None],
        "mappings": "",
    }))
    assert m.has_content is False


def test_sourcemap_decode():
    m = SourceMap.parse(json.dumps(_minimal_map()))
    assert m.decode() == [[(0, 0, 0, 0, -1)]]


# ---------------------------------------------------------------------------
# Inline source map extraction
# ---------------------------------------------------------------------------


def test_extract_inline_map_returns_none_when_absent():
    assert extract_inline_map("const x=1;") is None


def test_extract_inline_map_decodes_base64():
    payload = json.dumps(_minimal_map())
    b64 = base64.b64encode(payload.encode()).decode()
    js = f"const x=1;\n//# sourceMappingURL=data:application/json;base64,{b64}"
    extracted = extract_inline_map(js)
    assert extracted is not None
    m = SourceMap.parse(extracted)
    assert m.sources == ["in.js"]


def test_extract_inline_map_picks_last_occurrence():
    # Two inline maps: the last one is the effective one.
    m1 = json.dumps({"version": 3, "sources": ["first.js"], "mappings": ""})
    m2 = json.dumps({"version": 3, "sources": ["second.js"], "mappings": ""})
    b1 = base64.b64encode(m1.encode()).decode()
    b2 = base64.b64encode(m2.encode()).decode()
    js = (
        f"//# sourceMappingURL=data:application/json;base64,{b1}\n"
        f"//# sourceMappingURL=data:application/json;base64,{b2}"
    )
    extracted = extract_inline_map(js)
    m = SourceMap.parse(extracted)
    assert m.sources == ["second.js"]


def test_extract_inline_map_returns_none_on_bad_base64():
    js = "//# sourceMappingURL=data:application/json;base64,!!!notbase64!!!"
    # base64.b64decode with validate=False may not raise; but the
    # decoded bytes are not valid UTF-8 in this case, so we return None.
    assert extract_inline_map(js) is None


# ---------------------------------------------------------------------------
# Inline source map stripping
# ---------------------------------------------------------------------------


def test_strip_inline_map_no_op_when_absent():
    js = "const x=1;"
    assert strip_inline_map(js) == js


def test_strip_inline_map_removes_comment():
    js = "const x=1;\n//# sourceMappingURL=data:application/json;base64,AAAA"
    stripped = strip_inline_map(js)
    assert "sourceMappingURL" not in stripped
    assert stripped == "const x=1;"


def test_strip_inline_map_handles_no_newline_before():
    js = "const x=1;//# sourceMappingURL=data:application/json;base64,AAAA"
    stripped = strip_inline_map(js)
    assert "sourceMappingURL" not in stripped


def test_strip_inline_map_leaves_other_comments():
    js = "// regular comment\nconst x=1;\n//# sourceMappingURL=data:application/json;base64,AAAA"
    stripped = strip_inline_map(js)
    assert "regular comment" in stripped
    assert "sourceMappingURL" not in stripped


def test_strip_inline_map_empty_input():
    assert strip_inline_map("") == ""


def test_strip_inline_map_preserves_trailing_semicolon():
    js = "const x=1;\n//# sourceMappingURL=data:application/json;base64,AAAA"
    stripped = strip_inline_map(js)
    assert stripped.endswith(";")