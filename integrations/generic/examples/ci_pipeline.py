"""Using MinifyJS from a CI pipeline.

MinifyJS is designed to be the smallest possible step in a CI job:
no Node.js install, no npm install, no toolchain setup. This script
is what a CI step looks like.

Typical GitHub Actions workflow:

    - uses: actions/setup-python@v5
      with:
        python-version: "3.11"
    - run: pip install minifyjs
    - run: python integrations/generic/examples/ci_pipeline.py dist/js

The script:

  1. Walks a target directory for .js files.
  2. Minifies each in place with the .min.js suffix.
  3. Fails the CI step if any file has an error.
  4. Prints a summary of bytes saved.

Exit codes:
  0  all files minified successfully
  1  one or more files failed
  2  usage error
"""

from __future__ import annotations

import sys
from pathlib import Path
from typing import List, Tuple

if __package__ is None or __package__ == "":
    sys.path.insert(0, str(Path(__file__).resolve().parents[3]))

from minifyjs import Adapter, MinifyJSError, Options


class CIPipelineAdapter(Adapter):
    """In-place minification with failure tracking."""

    def __init__(self, root: Path, options: Options) -> None:
        super().__init__(options=options)
        self.root = Path(root)
        self.failures: List[Tuple[Path, str]] = []

    def should_process(self, path: Path) -> bool:
        if not super().should_process(path):
            return False
        # Skip generated artifacts.
        if path.name.endswith(".min.js"):
            return False
        return True

    def output_name(self, path: Path) -> Path:
        return Path(path).with_name(
            f"{Path(path).stem}{self.output_suffix}{Path(path).suffix}"
        )

    def run(self) -> int:
        processed = 0
        saved = 0
        for path in self._iter_files(self.root, recursive=True):
            if not self.should_process(path):
                continue
            try:
                result = self.minify_file(path)
                processed += 1
                saved += result.bytes_saved
            except MinifyJSError as e:
                self.failures.append((path, str(e)))

        print(f"processed: {processed} file(s)")
        print(f"saved:     {saved} bytes")

        if self.failures:
            print(f"failures:  {len(self.failures)}")
            for path, msg in self.failures:
                print(f"  {path}: {msg}")
            return 1
        return 0


def _main(argv: List[str]) -> int:
    if len(argv) != 2:
        print("usage: ci_pipeline.py TARGET_DIR", file=sys.stderr)
        return 2

    root = Path(argv[1])
    if not root.is_dir():
        print(f"error: {root} is not a directory", file=sys.stderr)
        return 1

    adapter = CIPipelineAdapter(
        root=root,
        options=Options(compress=True, mangle=True),
    )
    return adapter.run()


if __name__ == "__main__":
    raise SystemExit(_main(sys.argv))