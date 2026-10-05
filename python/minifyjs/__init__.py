"""minifyjs: a native JavaScript minifier, driven from Python.

    from minifyjs import minify

    result = minify("function add(a, b) { return a + b; }")
    print(result.code)

No Node.js or npm installation is required: this package bundles a
precompiled native binary and talks to it over stdin/stdout.

The public surface is:

    minify(source, **options) -> Result
    optimize(source, **options) -> Result
    bundle(entry_points, **options) -> Result

    Options, BundleOptions
    Result
    Diagnostic

    Adapter (for build-system integration)
    Config, load_discovered (for config file access)

    MinifyJSError and subclasses.
"""

from ._version import __version__
from .adapter import Adapter
from .bundler import bundle
from .config import Config
from .config import find as find_config
from .config import load as load_config
from .diagnostics import Diagnostic
from .errors import BinaryNotFoundError, BundleError, MinifyError, MinifyJSError
from .minifier import minify
from .optimizer import optimize
from .options import BundleOptions, Options
from .result import Result

__all__ = [
    "Adapter",
    "BinaryNotFoundError",
    "BundleError",
    "BundleOptions",
    "Config",
    "Diagnostic",
    "MinifyError",
    "MinifyJSError",
    "Options",
    "Result",
    "__version__",
    "bundle",
    "find_config",
    "load_config",
    "minify",
    "optimize",
]
