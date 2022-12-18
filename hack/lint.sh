#!/usr/bin/env bash
set -euo pipefail

main() {
  local pkg
  getopts 'w:' _
  pkg="${OPTARG:-./...}"

  echo 'Linting *.go files...'
  golangci-lint run "$pkg"
  echo 'Done.'
}

main "$@"
