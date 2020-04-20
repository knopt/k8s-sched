#!/bin/bash

DESTFILE=$1
DIST=$(python3 stresscpu-dist.py --dist uniform --cpu-load 50 --min 0 --max 100 --num 100 --sigma 15)

sed "s/DIST-PLACEHOLDER/$DIST/g" cpustresspod-template.yaml >"$DESTFILE"

kubectl apply -f "$DESTFILE"

echo "$DESTFILE"