#!/bin/bash

DIST=$(python3 stresscpu-dist.py --dist uniform --cpu-load 50 --min 0 --max 100 --num 9 --sigma 30)
R=""+$RANDOM
R=${R:1}
DATE=$(date +"%H-%M-%S-%d-%m-%Y")
UNIQE="$DATE-$R"
DESTFILE="cpustress-pod-uniform-auto-$UNIQE.yaml"

sed "s/DIST-PLACEHOLDER/$DIST/g" cpustresspod-template.yaml >"/tmp/$DESTFILE"
sed "s/DATE/$UNIQE/g" "/tmp/$DESTFILE" >"$DESTFILE"

#kubectl apply -f "$DESTFILE"

echo "$DESTFILE"