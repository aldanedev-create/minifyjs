"""Tests for config file discovery and loading."""

from __future__ import annotations

import json

from minifyjs.config import CONFIG_FILENAMES, find, load, load_discovered


def test_find_returns_none_when_absent(tmp_path):
    assert find(tmp_path) is None


def test_find_locates_config_in_dir(tmp_path):
    (tmp_path / "minifyjs.config.json").write_text("{}")
    found = find(tmp_path)
    assert found is not None
    assert found.name == "minifyjs.config.json"


def test_find_walks_up(tmp_path):
    sub = tmp_path / "a" / "b"
    sub.mkdir(parents=True)
    (tmp_path / "minifyjs.config.json").write_text("{}")
    found = find(sub)
    assert found is not None
    assert found.parent == tmp_path


def test_find_prefers_first_name(tmp_path):
    (tmp_path / "minifyjs.config.json").write_text("{}")
    (tmp_path / ".minifyjsrc").write_text("{}")
    found = find(tmp_path)
    assert found.name == "minifyjs.config.json"


def test_load_parses_target(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text(json.dumps({"target": "es2015"}))
    cfg = load(p)
    assert cfg.target == "es2015"


def test_load_minify_block(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text(json.dumps({
        "minify": {"whitespace": True, "identifiers": True, "syntax": True}
    }))
    cfg = load(p)
    assert cfg.compress is True
    assert cfg.mangle is True


def test_load_preserves_unknown_keys(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text(json.dumps({"futureKey": "value"}))
    cfg = load(p)
    assert cfg.extra["futureKey"] == "value"


def test_load_discovered_returns_none(tmp_path):
    assert load_discovered(tmp_path) is None


def test_load_discovered_finds_and_loads(tmp_path):
    (tmp_path / "minifyjs.config.json").write_text(
        json.dumps({"target": "esnext"})
    )
    cfg = load_discovered(tmp_path)
    assert cfg is not None
    assert cfg.target == "esnext"


def test_config_filenames_order():
    assert CONFIG_FILENAMES[0] == "minifyjs.config.json"


# NEW tests


def test_load_empty_object(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text("{}")
    cfg = load(p)
    assert cfg.target is None
    assert cfg.compress is False


def test_load_format(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text(json.dumps({"format": "esm"}))
    assert load(p).format == "esm"


def test_load_sourcemap(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text(json.dumps({"sourcemap": "inline"}))
    assert load(p).sourcemap == "inline"


def test_load_banner(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text(json.dumps({"banner": "/*b*/"}))
    assert load(p).banner == "/*b*/"


def test_load_footer(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text(json.dumps({"footer": "/*f*/"}))
    assert load(p).footer == "/*f*/"


def test_load_legal_comments(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text(json.dumps({"legalComments": "none"}))
    assert load(p).legal_comments == "none"


def test_load_output(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text(json.dumps({"output": "dist/out.js"}))
    assert load(p).output == "dist/out.js"


def test_load_inputs(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text(json.dumps({"inputs": ["a.js", "b.js"]}))
    assert load(p).inputs == ["a.js", "b.js"]


def test_load_compress_only(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text(json.dumps({"minify": {"whitespace": True}}))
    cfg = load(p)
    assert cfg.compress is True
    assert cfg.mangle is False


def test_load_mangle_only(tmp_path):
    p = tmp_path / "minifyjs.config.json"
    p.write_text(json.dumps({"minify": {"identifiers": True}}))
    cfg = load(p)
    assert cfg.mangle is True


def test_load_rc_json_variant(tmp_path):
    p = tmp_path / ".minifyjsrc.json"
    p.write_text(json.dumps({"target": "esnext"}))
    cfg = load(p)
    assert cfg.target == "esnext"


def test_find_ignores_directories_with_same_name(tmp_path):
    # A *directory* named minifyjs.config.json should not be returned.
    (tmp_path / "minifyjs.config.json").mkdir()
    assert find(tmp_path) is None
