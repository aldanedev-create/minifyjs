"""A minimal directory watcher using MinifyJS's per-file API.

This is a reference implementation of watch-mode outside the CLI. It
poll every N seconds, re-minifies files whose mtime changed, and
skips files that have not changed.

Usage:

    python watch_directory.py SOURCE_DIR [--interval 1.0]

Stop with Ctrl-C.

Why polling and not inotify: inotify, FSEvents, and ReadDirectoryChangesW
behave differently across platforms and fail on some network mounts.
Polling is the boring solution that works everywhere.
"""

from __future__ import annotations

import argparse
import sys
import time
from pathlib import Path
from typing import Dict, List

if __package__ is None or __package__ == "":
    sys.path.insert(0, str(Path(__file__).resolve().parents[3]))

from minifyjs import Adapter, MinifyJSError, Options


class WatchingAdapter(Adapter):
    """Minifies in place and tracks mtimes across runs."""

    def __init__(self, root: Path, options: Options) -> None:
        super().__init__(options=options)
        self.root = Path(root)
        self.mtimes: Dict[Path, float] = {}

    def output_name(self, path: Path) -> Path:
        return Path(path).with_name(
            f"{Path(path).stem}{self.output_suffix}{Path(path).suffix}"
        )

    def scan(self) -> List[Path]:
        """Return the set of processable files whose mtime changed."""
        changed: List[Path] = []
        for path in self._iter_files(self.root, recursive=True):
            if not self.should_process(path):
                continue
            mtime = path.stat().st_mtime
            if self.mtimes.get(path) != mtime:
                self.mtimes[path] = mtime
                changed.append(path)
        return changed

    def run_once(self) -> int:
        changed = self.scan()
        if not changed:
            return 0
        failures = 0
        for path in changed:
            try:
                self.minify_file(path)
                print(f"minified {path}")
            except MinifyJSError as e:
                print(f"error: {path}: {e}", file=sys.stderr)
                failures += 1
        return failures


def _main(argv: List[str]) -> int:
    parser = argparse.ArgumentParser(description="Watch a directory and minify on change")
    parser.add_argument("source_dir", type=Path)
    parser.add_argument("--interval", type=float, default=1.0)
    args = parser.parse_args(argv[1:])

    if not args.source_dir.is_dir():
        print(f"error: {args.source_dir} is not a directory", file=sys.stderr)
        return 1

    adapter = WatchingAdapter(
        root=args.source_dir,
        options=Options(compress=True, mangle=True),
    )

    # Initial pass.
    adapter.run_once()

    print(f"watching {args.source_dir} every {args.interval}s (Ctrl-C to stop)")
    try:
        while True:
            time.sleep(args.interval)
            adapter.run_once()
    except KeyboardInterrupt:
        print()
        return 0


if __name__ == "__main__":
    raise SystemExit(_main(sys.argv))