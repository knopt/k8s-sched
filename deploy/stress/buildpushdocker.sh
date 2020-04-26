#!/bin/bash

IMAGE=knopt/scipy:250401

docker build . -t "${IMAGE}"
docker push ${IMAGE}
