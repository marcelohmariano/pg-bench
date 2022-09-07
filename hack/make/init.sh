#!/usr/bin/env bash
set -e

[[ $DOCKER_ENABLED -ne 1 ]] && exit 0

docker ps >/dev/null || exit 1

for image in $(docker compose config --images); do
  [[ -n "$(docker images -q "$image")" ]] && continue
  docker compose build --build-arg UID="$(id -u)" --build-arg GID="$(id -g)"

  [[ -n "$(docker images -q "$image")" ]] && continue
  docker compose pull
done

docker compose up -d --remove-orphans --wait 2>/dev/null
