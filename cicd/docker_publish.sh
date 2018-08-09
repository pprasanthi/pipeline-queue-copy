#!/bin/sh

echo "🤖 Attempting to login"
echo "${DOCKER_TOKEN}" | docker login --username "${DOCKER_USER}" --password-stdin registry.hub.docker.com

IMAGE_NAME="$1"

docker push $IMAGE_NAME:latest

if [ -n "$MAJOR_MINOR" ]; then
  docker push $IMAGE_NAME:$MAJOR_MINOR
fi
