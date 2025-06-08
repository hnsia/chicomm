#!/bin/bash

mkdir -p dev/dist
buildpids=""

for f in chicomm-api chicomm-grpc chicomm-notification; do
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dev/dist/$f "./cmd/$f" &
	buildpids+=" $!"
done

for pid in $buildpids; do
    echo "process id: $pid"
    wait $pid
done

image="hns.test/chicomm:latest"
docker build -t "$image" -f dev/Dockerfile.dev .
echo "=> dev image built" > /dev/stderr