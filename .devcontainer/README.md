# Development container

A ready-to-use development environment for MinifyJS. Open the
repository in VS Code or GitHub Codespaces and the editor offers to
reopen it inside a container that already has Go, Python, the
linters, and the test tools installed.

## What you get

- **Go 1.23** with `golangci-lint` and `goimports`
- **Python 3.11** with `pytest`, `mypy`, `ruff`, and `pre-commit`
- **Node.js 20** for the benchmark suite
- **VS Code extensions** for Go, Python, YAML, and Markdown
- **Editor settings** that match the repository's `.editorconfig`
- **Port forwarding** for the Flask (`5000`) and FastAPI (`8000`)
  examples

Nothing is installed on your host machine except Docker and an
editor.

## Opening the container

### In VS Code

1. Install the **Dev Containers** extension.
2. Open the repository folder.
3. VS Code shows a notification: "Reopen in Container". Click it.

Or use the command palette: `Dev Containers: Reopen in Container`.

The first open builds the image, which takes a few minutes. Every
subsequent open is instant.

### In GitHub Codespaces

1. Open the repository on GitHub.
2. Click **Code** → **Codespaces** → **Create codespace on main**.

The codespace uses the same container definition. Setup is
automatic.

### From the command line

The `devcontainer` CLI builds and runs the container without an
editor:

```console
npm install -g @devcontainers/cli
devcontainer up --workspace-folder .
devcontainer exec --workspace-folder . bash
```

## What runs after the container starts

`.devcontainer/post-create.sh` runs once, as the `vscode` user. It:

1. Configures git for the mounted repository.
2. Runs `go work sync` and downloads Go modules.
3. Installs the Python package in editable mode with dev extras.
4. Builds the CLI into `core/bin/minifyjs`.
5. Copies the CLI into `python/minifyjs/bin/` so Python tests can
   find it.
6. Installs the pre-commit hooks.
7. Runs `npm install` in `bench/` for the benchmark tools.

When the script finishes, the development environment is ready.
The last lines print the common commands.

## Common commands

Inside the container:

```console
make build      # build the CLI
make test       # run the full test suite
make lint       # run linters
make bench      # run the benchmark suite
make fmt        # format everything
```

Or directly:

```console
cd core
go build ./...
go test ./...

cd ../python
python -m pytest tests

cd ../integrations/generic/tests
python -m pytest tests
```

## Port forwarding

The container forwards two ports automatically:

| Port | What runs there |
|---|---|
| 5000 | The Flask example (`examples/flask/app.py`) |
| 8000 | The FastAPI example (`examples/fastapi/main.py`) |

Running `python app.py` in `examples/flask/` inside the container
exposes the app on `http://localhost:5000/` on the host.

## Rebuilding the container

If the Dockerfile changes, rebuild:

- **VS Code:** Command palette → `Dev Containers: Rebuild Container`.
- **Codespaces:** Delete the codespace and create a new one.
- **CLI:** `devcontainer up --workspace-folder . --remove-existing-container`.

## Customizing

If you want to add tools, edit `.devcontainer/Dockerfile` and
rebuild. Do not edit `devcontainer.json` to install tools via
`postCreateCommand`; that makes the setup slower for everyone.

Changes to the Dockerfile are part of the repository and go through
the same review process as any other change.

## Troubleshooting

### The container takes too long to build

The first build downloads Go, Python, Node, and the toolchains.
This takes 5–10 minutes depending on network speed. Subsequent
builds use the Docker layer cache and take seconds.

If a rebuild is slow when it should be fast, the layer cache was
invalidated. This happens when a file copied into the image
changes, so anything before that `COPY` rebuilds.

### `go: command not found` in the container

The `PATH` is set in the Dockerfile via `ENV PATH=...`. If it is
not picking up, the container may have started with a stale shell.
Close and reopen the terminal, or run `exec bash`.

### Python tests skip with "bundled binary not present"

The binary was not copied into `python/minifyjs/bin/`. The
`post-create.sh` script copies it, but the copy step has three
fallback paths depending on platform. If none matched, run:

```console
cd core
go build -o ../python/minifyjs/bin/minifyjs ./cmd/minifyjs
```

Or on Windows:

```console
cd core
go build -o ../python/minifyjs/bin/minifyjs.exe ./cmd/minifyjs
```

### Ports are already in use on the host

If something on the host is using 5000 or 8000, VS Code will not be
able to forward the port. Either stop the conflicting process or
change the port the example listens on. The container's
`devcontainer.json` maps the container's port to the host port
one-to-one; the host port can be changed in VS Code's Ports panel.

### The container cannot access the internet

Some corporate networks block Docker's outbound traffic. The
`post-create.sh` script needs to reach `proxy.golang.org` (Go
modules), `pypi.org` (Python packages), and `registry.npmjs.org`
(Node modules). If any of those are blocked, the script fails and
the environment is incomplete.

Configure a proxy in `.devcontainer/devcontainer.json`'s
`containerEnv` section if needed.

## What the container is not

- **Not a security boundary.** The container runs as your user, has
  access to the network, and can write to the mounted workspace.
  Do not run untrusted code inside it.
- **Not the release build environment.** The release pipeline runs
  on GitHub's own runners, not in this container. The container is
  for interactive development.
- **Not required.** A contributor who has Go and Python installed
  locally can develop MinifyJS without ever using the container.
  This is a convenience, not a requirement.

## See also

- [../docs/development/setup.md](../docs/development/setup.md) —
  the manual setup instructions
- [../CONTRIBUTING.md](../CONTRIBUTING.md) — the contribution
  workflow