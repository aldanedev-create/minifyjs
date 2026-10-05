"""End-to-end integration tests exercising the whole stack."""

from __future__ import annotations

import subprocess
import sys

import pytest

from minifyjs import Adapter, Options, bundle, minify, optimize

from .conftest import BINARY_PATH

pytestmark = pytest.mark.binary


@pytest.fixture(autouse=True)
def _require_binary():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")


def test_python_api_and_cli_agree(tmp_path):
    """The Python API and the CLI must produce identical output for
    the same input and options."""
    src = "function add(a, b) { return a + b; }"
    api_result = minify(src, compress=True, mangle=True)

    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs", "--compress", "--mangle"],
        input=src.encode("utf-8"),
        capture_output=True,
        check=False,
    )
    assert proc.returncode == 0
    assert proc.stdout.decode("utf-8") == api_result.code


def test_adapter_and_direct_api_agree(tmp_path):
    src = "const x = 1 + 2;"
    direct = minify(src, compress=True, mangle=True)
    a = Adapter(options=Options(compress=True, mangle=True))
    in_file = tmp_path / "app.js"
    in_file.write_text(src)
    via_adapter = a.minify_file(in_file)
    assert via_adapter.code == direct.code


def test_optimize_equals_minify_with_both_flags():
    src = "const x = 1 + 2 + 3;"
    a = optimize(src)
    b = minify(src, compress=True, mangle=True)
    assert a.code == b.code


def test_bundle_then_optimize(tmp_path):
    """Bundle a small project, then minify the bundled output."""
    (tmp_path / "lib.js").write_text("export const add = (a, b) => a + b;\n")
    (tmp_path / "main.js").write_text(
        'import { add } from "./lib.js"; console.log(add(1, 2));\n'
    )
    out = tmp_path / "bundle.js"
    r = bundle(
        [str(tmp_path / "main.js")],
        outfile=str(out),
        format="esm",
        compress=True,
        mangle=True,
    )
    assert not r.has_errors
    body = out.read_text()
    # Second pass on the bundle output should be a no-op.
    r2 = optimize(body)
    assert not r2.has_errors


def test_multiple_files_via_adapter(tmp_path):
    for name in ("a.js", "b.js", "c.js"):
        (tmp_path / name).write_text("function f() { return 1; }")
    a = Adapter()
    results = a.minify_directory(tmp_path)
    assert len(results) == 3
    for name in ("a.min.js", "b.min.js", "c.min.js"):
        assert (tmp_path / name).is_file()


def test_options_from_config_to_minify(tmp_path):
    """A config file's options can be threaded into minify()."""
    from minifyjs.config import load
    cfg_path = tmp_path / "minifyjs.config.json"
    cfg_path.write_text('{"target":"es5"}')
    cfg = load(cfg_path)

    src = "const f = (a) => a + 1;"
    r = minify(src, target=cfg.target)
    assert "=>" not in r.code


def test_roundtrip_through_cache_directory(tmp_path):
    """Adapter on a directory, then again, should be idempotent."""
    (tmp_path / "app.js").write_text("const x = 1;")
    a = Adapter()
    r1 = a.minify_directory(tmp_path)
    # Second pass skips the .min.js file and processes app.js again.
    r2 = a.minify_directory(tmp_path)
    assert len(r2) == 1
    assert len(r1) == 1


def test_python_module_matches_function_api():
    """python -m minifyjs and minifyjs.minify produce the same bytes."""
    src = "const x = 1 + 2 + 3;"
    api = minify(src, compress=True)

    proc = subprocess.run(
        [sys.executable, "-m", "minifyjs", "--compress"],
        input=src.encode("utf-8"),
        capture_output=True,
        check=False,
    )
    assert proc.returncode == 0
    assert proc.stdout.decode("utf-8") == api.code
