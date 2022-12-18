#!/usr/bin/env bash
set -euo pipefail

build() {
  go build -ldflags='-s -w' -o "${BIN}/" "$1"
}

build_tool() {
  cd tools && go install "$1"
}

main() {
  local pkg
  local tool=false

  while getopts 'p:t' opt; do
    case "$opt" in
    p) pkg="$OPTARG" ;;
    t) tool=true ;;
    *) exit 1 ;;
    esac
  done

  $tool || echo "Building for $(go env GOOS)/$(go env GOARCH)"
  $tool && build_tool "$pkg" || build "$pkg"
  $tool || echo 'Done.'
}

main "$@"
