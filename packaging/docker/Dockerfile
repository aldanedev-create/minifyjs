# A minimal Docker image that contains the MinifyJS binary.
#
# The image is based on Alpine Linux so it is small (under 10 MB
# for the final layer). The binary is a static Go build, so it does
# not need a glibc base image.
#
# Build:
#
#     docker build -f packaging/docker/Dockerfile -t minifyjs .
#
# Run:
#
#     docker run --rm -i minifyjs < app.js > app.min.js
#     docker run --rm -v "$PWD:/src" minifyjs /src/app.js -o /src/app.min.js
#
# This image is published to Docker Hub as `minifyjs/minifyjs`. It
# is intended for CI environments that prefer a container over a
# Python install.

FROM golang:1.23-alpine AS builder

# Build dependencies. Only needed in the builder stage.
RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Copy the Go module first so the module download layer is cached.
COPY core/go.mod core/go.sum ./core/
RUN cd core && go mod download

# Copy the rest of the source.
COPY core ./core

# Build a static binary for the current architecture.
RUN cd core && \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath \
        -ldflags "-s -w -X github.com/minifyjs/minifyjs/core/internal/version.Version=$(cat ../version.txt 2>/dev/null || echo dev)" \
        -o /minifyjs \
        ./cmd/minifyjs

# ---------------------------------------------------------------------

FROM alpine:3.20

# The binary is static, so the runtime image only needs the
# certificate bundle (for HTTPS) and ca-certificates (for the
# system trust store). Nothing else.
RUN apk add --no-cache ca-certificates && \
    adduser -D -u 1000 minifyjs

COPY --from=builder /minifyjs /usr/local/bin/minifyjs

# Run as a non-root user.
USER minifyjs

# Default behavior: read from stdin, write to stdout.
ENTRYPOINT ["/usr/local/bin/minifyjs"]
CMD []