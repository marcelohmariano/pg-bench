#!/usr/bin/env bash
set -euo pipefail

test_no_cov() {
  go test -race "$1"
}

test_cov() {
  trap 'rm -f cover.out' EXIT
  go test -race -coverprofile=cover.out "$1"
  go tool cover -func=cover.out
}

main() {
  local pkg
  local cov

  while getopts 'w:c:' opt; do
    case "$opt" in
    w) pkg="${OPTARG:-./...}";;
    c) cov="$OPTARG";;
    *) exit 1 ;;
    esac
  done

  echo 'Running tests...'
  [[ "$cov" == 'y' ]] && test_cov "$pkg" || test_no_cov "$pkg"
  echo 'Done.'
}

main "$@"
