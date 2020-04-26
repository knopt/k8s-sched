#!/bin/bash

kubectl -n kube-system logs deploy/tknopik-scheduler -c tknopik-scheduler-extender-ctr -f
