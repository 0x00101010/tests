#!/usr/bin/env bash

# usage is download-logs.sh <circleci artifacts URL, pull this from the artifacts tab on the build> <folder-name>

set -euo pipefail

url="$1"
name="$2"

mkdir -p "$HOME/Downloads/logs/$2"

cd "$HOME/Downloads/logs/$2"

echo "Downloading logs..."
curl -o raw.log -L "$1"
echo "Stripping terminal characters..."
cat raw.log | sed -E 's/\x1B\[([0-9]{1,3}(;[0-9]{1,3})*)?[mGKHf]//g; s/\x1B\][0-9];.*\x07//g; s/\r//g' >cleaned.log
echo "Extracting log output..."
cat cleaned.log | jq -r '.Output' >output.log
