# Python: configuration

Reading and using `minifyjs.config.json` from Python.

## Two helpers

```python
from minifyjs import find_config, load_config, Config
```

| Function | What it does |
|---|---|
| `find_config(start=".")` | Walks up from `start` looking for a config file; returns a `Path` or `None` |
| `load_config(path)` | Parses the config file at `path`; returns a `Config` |
| `Config` | The parsed config as a dataclass |

## Finding the config

```python
from minifyjs import find_config

path = find_config()
if path is None:
    print("no config file found")
else:
    print(f"using {path}")
```

`find_config()` walks up from the current directory, checking for
`minifyjs.config.json`, `.minifyjsrc.json`, and `.minifyjsrc` in
that order.

Start the search elsewhere:

```python
from pathlib import Path
from minifyjs import find_config

path = find_config(Path("/some/other/dir"))
```

## Loading the config

```python
from minifyjs import load_config

cfg = load_config("minifyjs.config.json")
print(cfg.target)            # "es2015" or None
print(cfg.compress)          # True or False
print(cfg.mangle)            # True or False
print(cfg.sourcemap)         # "inline" or None
print(cfg.banner)            # "..." or None
print(cfg.extra)             # dict of unknown keys
```

## Find and load in one call

```python
from minifyjs import find_config, load_config, optimize

path = find_config()
if path:
    cfg = load_config(path)
    result = optimize(
        source,
        target=cfg.target,
        sourcemap=cfg.sourcemap,
        banner=cfg.banner,
    )
else:
    result = optimize(source)
```

## The `Config` dataclass

```python
@dataclass
class Config:
    inputs: List[str]
    output: Optional[str]
    target: Optional[str]
    format: Optional[str]
    sourcemap: Optional[str]
    banner: Optional[str]
    footer: Optional[str]
    legal_comments: Optional[str]
    compress: bool
    mangle: bool
    extra: Dict[str, Any]
```

`compress` and `mangle` are derived from the `minify` block in the
JSON:

```json
{
  "minify": {
    "whitespace": true,
    "identifiers": true,
    "syntax": true
  }
}
```

becomes:

```python
Config(compress=True, mangle=True, ...)
```

`extra` holds any keys the loader did not recognize, so a config
written for a newer MinifyJS version round-trips without losing
information.

## Path resolution

The `Config` dataclass does not resolve paths. Paths in the config
are relative to the config file's directory, and the loader does
not rewrite them because Python callers usually want to interpret
paths in their own context.

If you want Django-style resolution:

```python
from pathlib import Path
from minifyjs import load_config

config_path = Path("minifyjs.config.json")
cfg = load_config(config_path)

base = config_path.parent
output = base / cfg.output if cfg.output else None
```

## Using the config in a build script

```python
#!/usr/bin/env python3
"""Build script that respects minifyjs.config.json."""

from pathlib import Path
from minifyjs import Adapter, find_config, load_config, Options


def main() -> int:
    config_path = find_config()
    if config_path is None:
        print("no minifyjs.config.json found; using defaults")
        opts = Options(compress=True, mangle=True)
    else:
        print(f"using {config_path}")
        cfg = load_config(config_path)
        opts = Options(
            compress=cfg.compress,
            mangle=cfg.mangle,
            target=cfg.target,
            sourcemap=cfg.sourcemap,
            banner=cfg.banner,
        )

    adapter = Adapter(options=opts)
    results = adapter.minify_directory("static/js")
    print(f"minified {len(results)} files")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
```

## Environment variables

The Python API does not read environment variables. `MINIFYJS_*`
variables affect the CLI only. If you want the same behavior in
Python, read them yourself:

```python
import os
from minifyjs import optimize

target = os.environ.get("MINIFYJS_TARGET") or None
result = optimize(source, target=target)
```

## Where the config format is documented

The JSON format is documented in
[../reference/configuration.md](../reference/configuration.md). The
Python `Config` dataclass mirrors the JSON shape but does not cover
every key; keys the loader does not recognize end up in `extra`.

## See also

- [minify.md](minify.md) — passing options to `minify()`
- [optimize.md](optimize.md) — passing options to `optimize()`
- [../reference/configuration.md](../reference/configuration.md) —
  the config file reference