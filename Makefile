SHELL := /bin/bash
MAKEFLAGS += Rrs

.DEFAULT_GOAL := help

.ONESHELL:
.SUFFIXES:

CMD := cmd

export BIN := bin
export GOBIN := $(PWD)/$(BIN)
export PATH := $(BIN):$(PATH)

# The bench binary
BINARY := pg-bench

# Tools used during development
GOFMT := $(BIN)/goimports-reviser
TOOLS += $(GOFMT)

GOLANGCI_LINT := $(BIN)/golangci-lint
TOOLS += $(GOLANGCI_LINT)

TOOLS_GO_MOD := tools/go.mod

# @help: Show this help message. This is the default target.
.PHONY: help
help:
	./hack/help.sh $(MAKEFILE_LIST)

# @all: Build the bench binary.
.PHONY: all
all: build

.PHONY: build
build: $(BINARY)

# @clean: Remove the bench binary.
.PHONY: clean
clean:
	rm -f $(BIN)/$(BINARY)

# @fmt: Format Go source files that don't follow a standard formatting style.
# @fmt: Pass `WHAT=path/to/file/or/package` to format only a file or package.
.PHONY: fmt
fmt: $(GOFMT)
	./hack/fmt.sh -w '$(WHAT)'

# @lint: Lint Go source files.
# @lint: Pass `WHAT=path/to/file/or/package` to lint only a file or package.
.PHONY: lint
lint: $(GOLANGCI_LINT)
	./hack/lint.sh -w '$(WHAT)'

# @test: Run tests.
# @test: Pass `WHAT=path/to/package` to test only a package.
# @test: Pass `COV=y` to enable coverage analysis.
.PHONY: test
test: gen
	./hack/test.sh -w '$(WHAT)' -c '$(COV)'

.PHONY: gen
gen: $(MOCKGEN)
	go generate ./...

.PHONY: tidy
tidy:
	./hack/tidy.sh -w '$(WHAT)'

.PHONY: tools
tools: $(TOOLS)

.PHONY: $(BINARY)
$(BINARY):
	./hack/build.sh -p ./$(CMD)/$@

$(GOFMT): $(TOOLS_GO_MOD)
	./hack/build.sh -tp 'github.com/incu6us/goimports-reviser/v3'

$(GOLANGCI_LINT): $(TOOLS_GO_MOD)
	./hack/build.sh -tp 'github.com/golangci/golangci-lint/cmd/golangci-lint'
