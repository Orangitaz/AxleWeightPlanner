#!/bin/sh
# Rebuilds the Windows and Linux programs from index.html
set -e
cd "$(dirname "$0")"
cp index.html program/index.html
cd program
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-H=windowsgui -s -w" -o ../releases/AxleWeightPlanner.exe .
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ../releases/axle-planner-linux .
rm index.html
echo done
