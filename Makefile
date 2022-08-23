SHELL := /bin/bash
MAKEFLAGS += Rrs --warn-undefined-variable

.DEFAULT_GOAL := help

.ONESHELL:

export DOCKER_ENABLED ?= $(if $(wildcard /.dockerenv),0,1)

CONTAINER := build-env
INIT := ./hack/make/init.sh
RUN = $(INIT) && ./hack/make/run.sh $(CONTAINER)

.PHONY: help all clean lint test tidy seed init

# @help: Show this help message
help:
	./hack/make/help.sh $(MAKEFILE_LIST)

# @all: Build the benchmark binary
all:
	$(RUN) ./hack/make/build.sh ./cmd/benchmark

# @clean: Remove the benchmark binary
clean:
	rm -f ./bin/benchmark

# @lint: Lint the source files
lint:
	$(RUN) ./hack/make/lint.sh

# @test: Run tests
test:
	$(RUN) ./hack/make/test.sh

# @tidy: Update dependencies in `go.mod`
tidy:
	$(RUN) ./hack/make/tidy.sh

# @seed: Seed a local TimescaleDB instance with sample data
seed: CONTAINER := timescaledb
seed:
	$(RUN) ./hack/make/seed.sh

# @init: Initialize the development environment
init:
	$(INIT)
