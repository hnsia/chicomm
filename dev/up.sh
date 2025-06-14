#!/bin/bash

function wait_for() {
    attempts=0
    while ! eval $@ >/dev/null 2>&1; do
        if [[ $attempts -gt 99 ]] ; then
            echo "timed out waiting for $@"
            return 1
        fi
        attempts=$((attempts+1))
        sleep 3
    done
    return 0
}

# build dockerfile
echo "=> building containers" > /dev/stderr
dev/build.sh

# bring up mysql container
docker-compose -f dev/docker-compose.yaml up -d mysql

# wait for mysql container to come up
wait_for "docker-compose -f dev/docker-compose.yaml exec mysql mysql -uroot -ppassword -e 'SELECT 1;'"

# apply migration files
echo "=> applying migrations" > /dev/stderr
dev/migrate up

# bring up rest of the services
echo "=> bringing up rest of the services" > /dev/stderr
docker-compose -f dev/docker-compose.yaml up -d