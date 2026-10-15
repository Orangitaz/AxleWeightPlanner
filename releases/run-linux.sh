#!/bin/sh
# Starts the Axle Weight Planner on Linux from this folder (works from a thumb drive
# even when the drive can't run programs). Setups are kept in setups.json in this folder.
DIR="$(cd "$(dirname "$0")" && pwd)"
TMP="$(mktemp -d)"
cp "$DIR/axle-planner-linux" "$TMP/axle-planner" && chmod +x "$TMP/axle-planner" || { echo "Could not copy the program"; exit 1; }
AXLE_DATA_DIR="$DIR" "$TMP/axle-planner"
rm -rf "$TMP"
