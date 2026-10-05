"""Tests for binary discovery."""

from __future__ import annotations

import os
import stat

import pytest

from minifyjs._paths import bin_dir, binary_path
from minifyjs._platform import binary_filename
from minifyjs.errors import BinaryNotFoundError

from .conftest import BINARY_PATH


def test_binary_path_points_to_bin_dir():
    path = binary_path()
    assert path.startswith(bin_dir())
    assert path.endswith(binary_filename())


def test_find_binary_raises_when_missing(monkeypatch, tmp_path):
    # Point the loader at an empty directory.
    monkeypatch.setattr("minifyjs._paths.bin_dir", lambda: str(tmp_path))
    from minifyjs._binary import find_binary
    with pytest.raises(BinaryNotFoundError):
        find_binary()


@pytest.mark.binary
def test_find_binary_returns_path_when_present():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")
    from minifyjs._binary import find_binary
    path = find_binary()
    assert os.path.isfile(path)


@pytest.mark.binary
def test_find_binary_makes_executable(tmp_path, monkeypatch):
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")
    # Copy the binary to a temp dir, strip the +x bit, point the
    # loader there, and verify it restores the bit.
    dst = tmp_path / binary_filename()
    dst.write_bytes(BINARY_PATH.read_bytes())
    os.chmod(dst, stat.S_IRUSR | stat.S_IWUSR)  # 0o600
    monkeypatch.setattr("minifyjs._paths.binary_path", lambda: str(dst))
    from minifyjs._binary import find_binary
    path = find_binary()
    mode = os.stat(path).st_mode
    assert mode & stat.S_IXUSR


# NEW tests


def test_bin_dir_exists():
    assert os.path.isdir(bin_dir())


def test_binary_path_absolute():
    assert os.path.isabs(binary_path())


def test_binary_filename_no_path_separators():
    name = binary_filename()
    assert "/" not in name
    assert "\\" not in name


def test_binary_path_ends_with_filename():
    assert binary_path().endswith(binary_filename())


def test_binary_not_found_error_message_mentions_platform(monkeypatch, tmp_path):
    monkeypatch.setattr("minifyjs._paths.bin_dir", lambda: str(tmp_path))
    from minifyjs._binary import find_binary
    try:
        find_binary()
    except BinaryNotFoundError as e:
        assert "platform" in str(e).lower() or "-" in str(e)
    else:
        pytest.fail("expected BinaryNotFoundError")


def test_binary_found_returns_absolute_path():
    if not BINARY_PATH.is_file():
        pytest.skip("bundled binary not present")
    from minifyjs._binary import find_binary
    assert os.path.isabs(find_binary())