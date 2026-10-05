# Adding a platform

MinifyJS ships wheels for six platform and architecture pairs. This
page is the checklist for adding a seventh.

## What "adding a platform" means

The Go binary is compiled once per target and bundled into a
platform-specific Python wheel. Adding a platform means:

1. Adding an entry to `build/platforms.json`.
2. Adding a build step to the CI workflow for that platform.
3. Teaching the Python package to recognize the platform tag.
4. Adding a test.

That is the whole checklist. It is about thirty minutes of work
plus a CI run.

## Prerequisites

Before you start:

- **The platform must be able to run Go 1.22+.** If Go does not
  build for it, MinifyJS cannot ship for it.
- **The platform must have a PyPI-recognized wheel tag.** The
  `manylinux`, `macosx`, `win`, and `musllinux` families cover most
  cases.
- **You must be able to run the test suite on it.** A platform
  without CI is a platform without tests, and a platform without
  tests is a platform that will break silently.

## Step 1: Add an entry to `build/platforms.json`

Open `build/platforms.json` and add an object to the `platforms`
array:

```json
{
  "wheel_tag": "musllinux_1_1_x86_64",
  "human": "Alpine Linux x86-64",
  "goos": "linux",
  "goarch": "amd64",
  "binary_name": "minifyjs",
  "github_runner": "ubuntu-22.04"
}
```

| Field | Meaning |
|---|---|
| `wheel_tag` | The platform tag as it appears in the wheel filename |
| `human` | A short human-readable name for error messages |
| `goos` | Go's `GOOS` value (`linux`, `darwin`, `windows`, `freebsd`, ...) |
| `goarch` | Go's `GOARCH` value (`amd64`, `arm64`, `386`, ...) |
| `binary_name` | `minifyjs` on Unix, `minifyjs.exe` on Windows |
| `github_runner` | Which GitHub Actions runner builds it |

The `wheel_tag` must be a valid Python wheel platform tag. See
[PEP 425](https://peps.python.org/pep-0425/) for the rules.

## Step 2: Add a build step to CI

The CI workflow that builds binaries lives at
`.github/workflows/build-*.yml`. Each platform family has its own
workflow:

- `build-linux.yml` — builds the Linux wheels
- `build-macos.yml` — builds the macOS wheels
- `build-windows.yml` — builds the Windows wheels

If the new platform uses an existing runner (e.g. a musl wheel on
an Ubuntu runner), add a step to the existing workflow:

```yaml
- name: Build Alpine wheel
  run: python build/build_wheels.py --platform musllinux_1_1_x86_64 --version ${{ github.ref_name }}
```

If the new platform needs a new runner (e.g. `windows-11-arm64` for
Windows on ARM), add a new workflow file following the pattern of
the existing ones.

## Step 3: Teach the Python package the new tag

Open `python/minifyjs/_platform.py`. If the new platform's binary
name is different (e.g. `minifyjs.exe`), the `binary_filename()`
function already handles it via `sys.platform`.

If the platform needs a special case in `platform_tag()` (which is
only used in error messages), update it:

```python
def platform_tag() -> str:
    system = platform.system().lower()
    machine = platform.machine().lower()
    if system == "darwin":
        system = "macos"
    # ... existing code ...
    if system == "linux" and is_musl():
        system = "musl"
    return f"{system}-{machine}"
```

The `is_musl()` helper would need to be added. Detecting musl
reliably requires checking `platform.libc_ver()` or scanning the
filesystem; either is fine.

**Important:** the `platform_tag()` function is only used in error
messages. The actual binary selection is done by pip, which picks
the right wheel based on the platform tag. So this step is about
giving users useful error messages, not about functionality.

## Step 4: Add the platform to the release manifest

`build/generate_manifest.py` reads `build/platforms.json` and
enumerates every platform. Nothing to change here; the file is
generated from the JSON.

## Step 5: Add a test

The tests in `python/tests/test_platform.py` check the platform
tag's shape. If the new platform has a distinct shape (e.g. it
includes `musl` in the tag), add a test:

```python
def test_platform_tag_musl():
    # Only runs on musl systems.
    if not is_musl():
        pytest.skip("not running on musl")
    assert "musl" in platform_tag()
```

## Step 6: Update the docs

Three places:

- [`../compatibility.md`](../compatibility.md) — add the platform
  to the supported list.
- [`../python/platform-support.md`](../python/platform-support.md) —
  same.
- [`../getting-started/installation.md`](../getting-started/installation.md) —
  if the platform has special install instructions (e.g. Alpine
  needs a source install today), update those.

## Step 7: Cut a release

New platforms ship with the next release. There is no separate
"release" for a single platform. Add the platform to the JSON, get
the CI green, and the next release publishes it.

## What you cannot do

### Add a platform without CI

A platform without a CI job is a platform without tests. If you add
a wheel for a platform that CI cannot build, the platform will
break silently when an esbuild upgrade changes something.

Do not add a platform without a CI runner for it.

### Add a platform without a test

Same reasoning. Every platform must have at least one test that
verifies the binary runs.

The test can be as simple as:

```python
def test_binary_runs():
    from minifyjs import minify
    result = minify("const x = 1;")
    assert result.code == "const x=1;"
```

If that passes on the platform, the platform is supported.

### Add a platform that Go cannot target

If Go does not build for the platform, MinifyJS cannot ship for it.
No workaround exists that does not involve cross-compiling from a
different platform, which is not what any of the current workflows
do.

### Add a platform for a Python version

Platform and Python version are independent. Adding Python 3.13 to
the support matrix is a different task (see
[`../compatibility.md`](../compatibility.md)).

## Debugging a new platform

### "binary not executable"

The wheel was built but the binary does not have the executable
bit. On Unix, the fix is to ensure the file mode is preserved
through the wheel build. The wheel format stores file modes;
something in the build chain is stripping them.

Check `build/build_wheels.py`. The file is added with:

```python
entries.append((bin_arcname, binary_path.read_bytes()))
```

The `zipfile` module's `writestr` does not preserve mode by
default. To preserve it, write the file with an explicit mode:

```python
import zipfile
info = zipfile.ZipInfo(bin_arcname, date_time=(1980, 1, 1, 0, 0, 0))
info.external_attr = 0o755 << 16  # rwxr-xr-x
info.compress_type = zipfile.ZIP_DEFLATED
zf.writestr(info, binary_path.read_bytes())
```

This is a known issue with `build_wheels.py` today. If you add a
platform and the binary is not executable, this is the fix.

### "wheel not installable"

The wheel tag is wrong. Verify it against the platform's expected
tag. On the target platform:

```console
python -c "from packaging.tags import sys_tags; print(next(sys_tags()))"
```

The output is the wheel tag pip will look for. If it does not
match, the wheel will not install.

### "wheel installable but binary does not run"

The binary was cross-compiled for a different architecture. Check
`platforms.json`'s `goos` and `goarch` against the actual target:

```console
file /path/to/minifyjs
```

The output should name the correct architecture.

## See also

- [../compatibility.md](../compatibility.md) — the current support
  matrix
- [`build/README.md`](https://github.com/aldanedev-create/minifyjs/tree/main/build) —
  the release pipeline
- [release-process.md](release-process.md) — the release steps