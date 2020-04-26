#!/bin/bash

DATE=$(date +"%H-%M-%S-%d-%m-%Y")
DESTFILE="cpustress-deployment-uniform-auto-$DATE.yaml"
DIST=$(python3 stresscpu-dist.py --dist uniform --cpu-load 20 --min 0 --max 35 --num 100 --sigma 10)

sed "s/DIST-PLACEHOLDER/$DIST/g" cpustressdeployment-template.yaml >"/tmp/$DESTFILE"
sed "s/DATE/$DATE/g" "/tmp/$DESTFILE" >"$DESTFILE"

kubectl apply -f "$DESTFILE"

echo "$DESTFILE"