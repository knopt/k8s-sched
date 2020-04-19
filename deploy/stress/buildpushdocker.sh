#!/bin/bash

IMAGE=knopt/scipy:1804

docker build . -t "${IMAGE}"
docker push ${IMAGE}
