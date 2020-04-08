#!/bin/bash

pushd ../ 

IMAGE=knopt/k8s-ext:3110

docker build . -t "${IMAGE}"
docker push ${IMAGE}
