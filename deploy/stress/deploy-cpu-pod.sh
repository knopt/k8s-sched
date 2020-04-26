#!/bin/bash

DATE=$(date +"%H-%M-%S-%d-%m-%Y")
DESTFILE="cpustress-pod-uniform-auto-$DATE.yaml"
DIST=$(python3 stresscpu-dist.py --dist uniform --cpu-load 20 --min 0 --max 40 --num 7 --sigma 15)

sed "s/DIST-PLACEHOLDER/$DIST/g" cpustresspod-template.yaml >"/tmp/$DESTFILE"
sed "s/DATE/$DATE/g" "/tmp/$DESTFILE" >"$DESTFILE"

kubectl apply -f "$DESTFILE"

echo "$DESTFILE"