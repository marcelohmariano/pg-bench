#!/usr/bin/env bash
set -e

[[ $DOCKER_ENABLED -ne 1 ]] && exit 0

docker ps >/dev/null || exit 1

build_env_image='timescaledb-benchmark-build-env'
build_env_service='build-env'

timescaledb_image='timescale/timescaledb-ha:pg14-latest'
timescaledb_service='timescaledb'

if [[ -z "$(docker images -q "$build_env_image")" ]]; then
  docker compose build \
    --build-arg UID="$(id -u)" \
    --build-arg GID="$(id -g)" \
    "$build_env_service"
fi

if [[ -z "$(docker images -q "$timescaledb_image")" ]]; then
  docker compose pull "$timescaledb_service"
fi

docker compose up -d --remove-orphans --wait 2>/dev/null
