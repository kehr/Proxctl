#!/usr/bin/env bash
set -euo pipefail

target="${1:-vps}"

cd "$(dirname "$0")/.."
make linux-amd64
scp build/proxctl-linux-amd64 "$target:/tmp/proxctl"
ssh "$target" 'install -m 0755 /tmp/proxctl /usr/local/bin/proxctl && rm -f /tmp/proxctl && proxctl version'

