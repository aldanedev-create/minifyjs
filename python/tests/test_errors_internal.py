"""Tests for internal error message builders."""

from __future__ import annotations

from minifyjs._errors import (
    binary_not_executable_message,
    binary_not_found_message,
)


def test_not_found_message_mentions_path():
    msg = binary_not_found_message("/path/to/bin", "linux-x86_64")
    assert "/path/to/bin" in msg


def test_not_found_message_mentions_platform():
    msg = binary_not_found_message("/path/to/bin", "linux-x86_64")
    assert "linux-x86_64" in msg


def test_not_found_message_mentions_docs():
    msg = binary_not_found_message("/p", "linux-x86_64")
    assert "docs/" in msg or "installation" in msg.lower()


def test_not_executable_message_mentions_path():
    msg = binary_not_executable_message("/path/to/bin")
    assert "/path/to/bin" in msg


def test_not_executable_message_mentions_permission():
    msg = binary_not_executable_message("/p")
    assert "permission" in msg.lower() or "executable" in msg.lower()


def test_messages_are_strings():
    assert isinstance(binary_not_found_message("/p", "linux-x86_64"), str)
    assert isinstance(binary_not_executable_message("/p"), str)