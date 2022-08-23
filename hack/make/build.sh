#!/usr/bin/env bash
set -e

GOOS="$(tr '[:upper:]' '[:lower:]' <<<"$GOOS")"
GOARCH="$(tr '[:upper:]' '[:lower:]' <<<"$GOARCH")"

if ! go tool dist list &>/dev/null; then
  unset GOARCH
fi

echo "Building for ${GOOS:-$(uname -s)}/${GOARCH:-$(uname -m)}..."
go build -ldflags "-s -w" -o ./bin/ "$@"
echo 'Done.'
