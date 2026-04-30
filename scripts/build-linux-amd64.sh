#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o build/proxctl-linux-amd64 ./src/cmd/proxctl

