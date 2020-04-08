#!/bin/bash

kubectl -n kube-system logs deploy/my-scheduler -c my-scheduler-extender-ctr -f
