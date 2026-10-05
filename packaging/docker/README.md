# Docker image

A minimal Docker image that contains the MinifyJS binary, for CI
environments that prefer a container over a Python install.

## What is in the image

- The `minifyjs` binary, statically linked.
- Alpine Linux 3.20 as the base.
- A non-root user (`minifyjs`).
- `ca-certificates` for HTTPS.

Nothing else. The image is about 12 MB compressed.

## Building

From the repository root:

```console
docker build -f packaging/docker/Dockerfile -t minifyjs .
```

Or with a version:

```console
docker build \
    --build-arg VERSION=0.1.0 \
    -f packaging/docker/Dockerfile \
    -t minifyjs:0.1.0 .
```

## Running

Read from stdin, write to stdout:

```console
docker run --rm -i minifyjs < app.js > app.min.js
```

Minify a file on the host:

```console
docker run --rm -v "$PWD:/src" minifyjs /src/app.js -o /src/app.min.js
```

On Windows, the volume mount needs a different path syntax:

```powershell
docker run --rm -v "${PWD}:/src" minifyjs /src/app.js -o /src/app.min.js
```

With options:

```console
docker run --rm -i minifyjs --compress --mangle < app.js > app.min.js
```

## The published image

The release pipeline publishes `minifyjs/minifyjs` to Docker Hub
with two tags per release:

- `minifyjs/minifyjs:latest` — the newest release
- `minifyjs/minifyjs:0.1.0` — a specific version

Use a specific tag in CI. `latest` is convenient but not
reproducible.

```yaml
- name: Minify
  run: |
    docker run --rm -i minifyjs/minifyjs:0.1.0 \
        --compress --mangle < static/app.js > static/app.min.js
```

## Why Alpine

Alpine's base image is 5 MB. Debian's is 120 MB. For a static Go
binary, there is no reason to pay for glibc.

The binary is compiled with `CGO_ENABLED=0`, so it does not link
against any system library. Alpine's musl libc is not used.

## Why non-root

The image runs as a non-root user (`minifyjs`, uid 1000). This
matters when the container is used in CI: a file written by the
container should be owned by the CI user, not by root.

If you mount a volume, the ownership of the output file depends on
the mount. The container does not change ownership.

## Docker Compose

For a project with a compose file:

```yaml
services:
  minify:
    image: minifyjs/minifyjs:0.1.0
    volumes:
      - ./static:/src
    working_dir: /src
    command: ["--compress", "--mangle", "app.js", "-o", "app.min.js"]
```

## Alternatives

If you are already installing Python dependencies, use the Python
package instead of the Docker image. `pip install minifyjs` pulls
the same binary into a Python environment and avoids the container
overhead.

The Docker image exists for environments where Python is not
available or where a container is the preferred isolation boundary.

## Size

The final image is about 12 MB. Breakdown:

| Layer | Size |
|---|---|
| Alpine base | ~5 MB |
| ca-certificates | ~1 MB |
| The binary | ~6 MB |
| Metadata | < 1 MB |

The binary is about 6 MB because esbuild is compiled into it.
esbuild itself is about 5 MB of Go code.