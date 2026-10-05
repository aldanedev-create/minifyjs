# Makefile for the release automation. Every target delegates to a
# Python script; the Makefile exists so common operations are one
# short command instead of a long one.

SHELL := /bin/sh

# The version we are building. Populated from the git tag when
# running in CI; falls back to "dev" for local runs.
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

# Where intermediate artifacts go.
DIST ?= dist
WHEELHOUSE ?= wheelhouse

.PHONY: all build-all wheels release clean verify checksums manifest

all: build-all

build-all:
	./cross_compile.sh

wheels: build-all
	python build_wheels.py --all --version $(VERSION) --out $(WHEELHOUSE)

release:
	python build_release.py --version $(VERSION)

verify:
	@for f in $(DIST)/*; do \
		python scripts/verify_binary.py $$f || exit 1; \
	done

checksums:
	python generate_checksums.py --dir $(WHEELHOUSE) --out $(WHEELHOUSE)/checksums.txt

manifest:
	python generate_manifest.py --dir $(WHEELHOUSE) --out $(WHEELHOUSE)/manifest.json --version $(VERSION)

clean:
	rm -rf $(DIST) $(WHEELHOUSE)