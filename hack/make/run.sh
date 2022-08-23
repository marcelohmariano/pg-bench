#!/usr/bin/env bash
set -e

if [[ $DOCKER_ENABLED -ne 1 ]]; then
  set -- "${@:2}"
else
  set -- docker compose exec \
    -u "$(id -u):$(id -g)" \
    -e GOOS="${GOOS:-$(uname -s)}" \
    -e GOARCH="${GOARCH:-$(uname -m)}" \
    "$1" "${@:2}"
fi

exec "$@"
