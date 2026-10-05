"""Custom Django management command that minifies static assets.

Install by placing this file at:

    <your_app>/management/commands/minifystatic.py

where `<your_app>` is any app listed in INSTALLED_APPS. For this
example, the app is `staticfiles_example`, registered in
settings.py.
"""

from __future__ import annotations

from pathlib import Path

from django.conf import settings
from django.core.management.base import BaseCommand

from minifyjs import Adapter


class DjangoStaticAdapter(Adapter):
    """Writes minified files to STATIC_ROOT, preserving the path
    relative to whichever STATICFILES_DIRS entry they came from."""

    def __init__(self, static_root: Path, static_dirs: list) -> None:
        super().__init__()
        self.static_root = static_root
        self.static_dirs = static_dirs

    def output_name(self, path: Path) -> Path:
        path = Path(path)
        for src in self.static_dirs:
            try:
                rel = path.resolve().relative_to(src.resolve())
            except ValueError:
                continue
            return self.static_root / rel.parent / f"{rel.stem}.min{rel.suffix}"
        return self.static_root / f"{path.stem}.min{path.suffix}"


class Command(BaseCommand):
    help = "Minify every JavaScript file in STATICFILES_DIRS into STATIC_ROOT"

    def add_arguments(self, parser):
        parser.add_argument(
            "--no-mangle",
            action="store_true",
            help="Disable identifier mangling",
        )

    def handle(self, *args, **options):
        from minifyjs import Options

        static_root = Path(settings.STATIC_ROOT)
        static_dirs = [Path(p) for p in settings.STATICFILES_DIRS]
        static_root.mkdir(parents=True, exist_ok=True)

        opts = Options(compress=True, mangle=not options["no_mangle"])
        adapter = DjangoStaticAdapter(static_root, static_dirs)
        adapter.default_options = opts

        total = 0
        saved = 0
        for d in static_dirs:
            if not d.is_dir():
                continue
            for result in adapter.minify_directory(d):
                total += 1
                saved += result.bytes_saved

        self.stdout.write(
            self.style.SUCCESS(
                f"minified {total} file(s); saved {saved} bytes"
            )
        )