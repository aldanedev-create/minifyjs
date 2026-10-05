# Top-level Makefile. Delegates to the module Makefiles and the
# release/packaging scripts.
#
# Contributor entry points:
#   make build       compile the CLI into core/bin/minifyjs
#   make test        run every Go and Python test
#   make lint        run linters for both languages
#   make fmt         format everything
#   make bench       run the benchmark suite
#   make docs        regenerate CLI reference and syntax matrix
#   make clean       remove build artifacts
#
# Release entry points (maintainers only, driven by CI):
#   make release     full release pipeline (see build/build_release.py)

SHELL := /bin/sh

.PHONY: all build test test-go test-python lint fmt bench docs clean release

all: build

build:
	$(MAKE) -C core build

test: test-go test-python

test-go:
	$(MAKE) -C core test

test-python:
	cd python && python -m pytest tests

lint:
	$(MAKE) -C core lint
	cd python && python -m ruff check . && python -m mypy minifyjs

fmt:
	$(MAKE) -C core fmt
	cd python && python -m ruff format .

bench:
	cd bench && python run_all.py

docs:
	python tools/docs/generate_cli_reference.py --output docs/cli/reference.md
	python tools/docs/generate_syntax_matrix.py --output docs/supported-syntax.md

clean:
	$(MAKE) -C core clean
	rm -rf build/dist build/wheelhouse

release:
	python build/build_release.py