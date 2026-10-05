# Makefile for the MinifyJS core Go module.
#
# All targets assume they are run from core/. The top-level Makefile
# delegates here.
#
# Variables:
#   VERSION  release version string (overridden by CI from the git tag)
#   BINARY   output path of the compiled binary

SHELL := /bin/sh

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BINARY  ?= bin/minifyjs
LDFLAGS := -s -w -X github.com/minifyjs/minifyjs/core/internal/version.Version=$(VERSION)

.PHONY: all build test test-race vet lint fmt tidy clean install release

all: build

build:
	@mkdir -p $(dir $(BINARY))
	go build -trimpath -ldflags '$(LDFLAGS)' -o $(BINARY) ./cmd/minifyjs

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .
	goimports -w . 2>/dev/null || true

tidy:
	go mod tidy

clean:
	rm -rf bin

install: build
	install -m 0755 $(BINARY) $(GOPATH)/bin/minifyjs

release:
	@echo "release builds are driven by build/build_release.py, not make"
	@exit 1