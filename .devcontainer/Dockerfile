# The development container image.
#
# This is not the same as the release Dockerfile
# (packaging/docker/Dockerfile), which builds a minimal runtime
# image. This one is for contributors: it includes the Go toolchain,
# the Python toolchain, linters, and test runners.
#
# Build:
#
#     docker build -f .devcontainer/Dockerfile -t minifyjs-dev .
#
# The image is normally built automatically by VS Code or GitHub
# Codespaces, which reads .devcontainer/devcontainer.json.

ARG GO_VERSION=1.23
ARG PYTHON_VERSION=3.11

# ---------------------------------------------------------------------
# Stage 1: Go toolchain and Go-based development tools.
# ---------------------------------------------------------------------
FROM golang:${GO_VERSION}-bookworm AS go-tools

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        git \
 && rm -rf /var/lib/apt/lists/*

# Install golangci-lint. The version is pinned so all contributors
# use the same one; CI uses the same version.
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b /usr/local/bin v1.61.0

# Install goimports, used by the pre-commit hooks.
RUN go install golang.org/x/tools/cmd/goimports@latest

# ---------------------------------------------------------------------
# Stage 2: Python toolchain and Python-based development tools.
# ---------------------------------------------------------------------
FROM python:${PYTHON_VERSION}-bookworm AS python-tools

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        git \
 && rm -rf /var/lib/apt/lists/*

# Install the tools CI uses. Versions are pinned to match CI.
RUN pip install --no-cache-dir \
        pytest==8.3.3 \
        pytest-cov==6.0.0 \
        mypy==1.11.2 \
        ruff==0.6.9 \
        pre-commit==3.8.0 \
        build==1.2.2 \
        twine==5.1.1 \
        psutil==6.0.0 \
        rjsmin==1.2.2

# ---------------------------------------------------------------------
# Stage 3: The final development image.
# ---------------------------------------------------------------------
FROM python:${PYTHON_VERSION}-bookworm

ARG GO_VERSION=1.23

# System packages a contributor is likely to want.
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        wget \
        git \
        bash-completion \
        zsh \
        less \
        jq \
        tree \
        make \
        shellcheck \
        gnupg \
        lsb-release \
        procps \
        htop \
        unzip \
        zip \
 && rm -rf /var/lib/apt/lists/*

# Go toolchain. Copied from the official image so it matches
# whatever Go version the project builds against.
COPY --from=golang:${GO_VERSION}-bookworm /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/home/vscode/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Development tools from the earlier stages.
COPY --from=go-tools /usr/local/bin/golangci-lint /usr/local/bin/golangci-lint
COPY --from=go-tools /root/go/bin/goimports /usr/local/bin/goimports
COPY --from=python-tools /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=python-tools /usr/local/bin/pytest /usr/local/bin/pytest
COPY --from=python-tools /usr/local/bin/ruff /usr/local/bin/ruff
COPY --from=python-tools /usr/local/bin/mypy /usr/local/bin/mypy
COPY --from=python-tools /usr/local/bin/pre-commit /usr/local/bin/pre-commit
COPY --from=python-tools /usr/local/bin/build /usr/local/bin/build
COPY --from=python-tools /usr/local/bin/twine /usr/local/bin/twine

# Node.js for the benchmark suite. Not needed for the engine, but
# the benchmark scripts use it to install esbuild, terser, and
# uglify-js.
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
 && apt-get install -y --no-install-recommends nodejs \
 && rm -rf /var/lib/apt/lists/*

# Create the vscode user. This is the user the container runs as.
# The UID is 1000, which is the default for the first non-root user
# on most Linux systems.
RUN useradd --create-home --shell /bin/bash --uid 1000 vscode \
 && mkdir -p /home/vscode/go \
 && chown -R vscode:vscode /home/vscode

# Switch to the vscode user for the rest of the image.
USER vscode
WORKDIR /home/vscode

# Configure Go for the vscode user.
RUN go env -w GOTOOLCHAIN=local

# Configure the shell. The container's default shell is bash;
# zsh is available for contributors who prefer it.
RUN echo 'export PATH="$HOME/go/bin:$PATH"' >> /home/vscode/.bashrc \
 && echo 'export PATH="$HOME/.local/bin:$PATH"' >> /home/vscode/.bashrc \
 && echo 'alias ll="ls -alF"' >> /home/vscode/.bashrc

USER root

# The workspace is mounted by the devcontainer spec at
# /workspaces/minifyjs. This directory is just the mount point.
RUN mkdir -p /workspaces && chown vscode:vscode /workspaces

USER vscode

# The default command. The container stays alive until the editor
# disconnects.
CMD ["sleep", "infinity"]