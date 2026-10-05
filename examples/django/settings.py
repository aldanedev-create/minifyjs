"""Minimal Django settings for the MinifyJS example."""

from pathlib import Path

BASE_DIR = Path(__file__).parent

SECRET_KEY = "minifyjs-example-key-do-not-use-in-production"

DEBUG = True

ALLOWED_HOSTS = ["*"]

INSTALLED_APPS = [
    "django.contrib.contenttypes",
    "django.contrib.staticfiles",
]

STATIC_URL = "/static/"

STATICFILES_DIRS = [
    BASE_DIR / "static",
]

STATIC_ROOT = BASE_DIR / "static_root"

USE_TZ = True