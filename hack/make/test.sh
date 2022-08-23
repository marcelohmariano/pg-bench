#!/usr/bin/env bash
set -e

unset GOOS
unset GOARCH

echo 'Running tests...'
go test -race ./...
echo 'Done.'
