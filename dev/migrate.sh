#!/bin/bash

set -e

export HOST=${HOST:-0.0.0.0}
export PORT=${PORT:-3305}
export LOGIN=${LOGIN:-root}
export PASSWORD=${PASSWORD:-password}
export DATABASE=${DATABASE:-chicomm}

docker_cmd=(docker run -it --rm --network host --volume "$(pwd)/db:/db" migrate/migrate:v4.17.0)

if ["${1:-}" == "create"]; then
    shift
    exec "${docker_cmd[@]}" create -ext sql -dir /db/migrations "$@"
fi

exec "${docker_cmd[@]}" -path=/db/migrations -database "mysql://${LOGIN}:${PASSWORD}@tcp(${HOST}:${PORT})/${DATABASE}" "$@"