#!/bin/bash
make build-docker
docker login --username "$DOCKER_USERNAME" --password "$DOCKER_PASSWORD"
# build gosharexserver image
export REPO=mmichaelb/gosharexserver
docker build -f ./build/gosharexserver/Dockerfile -t ${REPO}:$1 .
docker tag ${REPO}:$1 ${REPO}:latest
docker push ${REPO}
# build gosharexserver-user-adder tool image
export REPO=mmichaelb/gosharexserver-user-adder
docker build -f ./build/gosharexserver-user-adder/Dockerfile -t ${REPO}:$1 .
docker tag ${REPO}:$1 ${REPO}:latest
docker push ${REPO}
