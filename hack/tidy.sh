#!/usr/bin/env bash
set -euo pipefail

getopts 'w:' _

echo 'Updating dependencies...'
[[ -n "$OPTARG" ]] && cd "$OPTARG"
go mod tidy
echo 'Done.'
