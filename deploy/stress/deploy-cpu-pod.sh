#!/bin/bash

DATE=$(date +"%H-%M-%S-%d-%m-%Y")
DESTFILE="cpustress-pod-uniform-auto-$DATE.yaml"
DIST=$(python3 stresscpu-dist.py --dist uniform --cpu-load 50 --min 0 --max 100 --num 9 --sigma 0)

sed "s/DIST-PLACEHOLDER/$DIST/g" cpustresspod-template.yaml >"/tmp/$DESTFILE"
sed "s/DATE/$DATE/g" "/tmp/$DESTFILE" >"$DESTFILE"

kubectl apply -f "$DESTFILE"

echo "$DESTFILE"