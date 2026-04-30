#!/usr/bin/env bash
set -euo pipefail

mode="${1:-snapshot}"
base="$(tr -d '[:space:]' < VERSION)"
short_sha="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"

case "$mode" in
  release)
    printf '%s\n' "$base"
    ;;
  snapshot)
    printf '%s-dev.%s\n' "$base" "$short_sha"
    ;;
  tag-check)
    tag="${2:-}"
    if [[ -z "$tag" ]]; then
      echo "missing tag" >&2
      exit 2
    fi
    if [[ "$tag" != "$base" ]]; then
      echo "release tag $tag does not match VERSION $base" >&2
      exit 1
    fi
    ;;
  *)
    echo "usage: scripts/version.sh <snapshot|release|tag-check> [tag]" >&2
    exit 2
    ;;
esac
