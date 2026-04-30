#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."
make linux-amd64
