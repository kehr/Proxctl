#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."
make build
install -m 0755 build/proxctl /usr/local/bin/proxctl

