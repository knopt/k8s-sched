#!/bin/bash

IMAGE=knopt/scipy:2604

docker build . -t "${IMAGE}"
docker push ${IMAGE}
