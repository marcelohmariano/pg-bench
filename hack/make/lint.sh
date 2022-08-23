#!/usr/bin/env bash
set -e

unset GOOS
unset GOARCH

echo 'Linting *.go files...'
golangci-lint run ./...

echo 'Linting *.sh files...'
shellcheck ./hack/make/*.sh

echo 'Done.'
