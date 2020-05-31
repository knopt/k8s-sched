#!/bin/bash

pushd ../ 

IMAGE=knopt/k8s-ext:100501

docker build . -t "${IMAGE}"
docker push ${IMAGE}
