#!/bin/bash
make build-docker
docker login --username "$DOCKER_USERNAME" --password "$DOCKER_PASSWORD"
export REPO=mmichaelb/gosharexserver
docker build -f ./build/package/Dockerfile -t ${REPO}:$1 .
docker tag ${REPO}:$1 ${REPO}:latest
docker push ${REPO}
