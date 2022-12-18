#!/usr/bin/env bash
set -euo pipefail

format() {
  for pkg in "$@"; do
    goimports-reviser -rm-unused -format -recursive "$pkg"
  done
}

main() {
  getopts 'b:w:' _
  local what=("${OPTARG}")

  if [[ "${#what}" -eq 0 ]]; then
    what=(cmd internal)
  fi

  echo 'Formatting *.go files...'
  format "${what[@]}"
  echo 'Done.'
}

main "$@"
