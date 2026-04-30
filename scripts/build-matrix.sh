#!/usr/bin/env bash
set -euo pipefail

mode="${1:-snapshot}"
root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

if [[ "$mode" == "release" ]]; then
  version="${VERSION:-$(./scripts/version.sh release)}"
else
  version="${VERSION:-$(./scripts/version.sh snapshot)}"
fi
commit="${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo unknown)}"
date="${DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"
pkg="./src/cmd/proxctl"
ldflags="-s -w -X github.com/kehr/proxctl/src/internal/cli.Version=${version} -X github.com/kehr/proxctl/src/internal/cli.Commit=${commit} -X github.com/kehr/proxctl/src/internal/cli.Date=${date}"

rm -rf dist
mkdir -p dist

targets=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
  "windows amd64"
)

for target in "${targets[@]}"; do
  read -r goos goarch <<<"$target"
  name="proxctl_${version}_${goos}_${goarch}"
  bin="proxctl"
  if [[ "$goos" == "windows" ]]; then
    bin="proxctl.exe"
  fi
  work="dist/${name}"
  mkdir -p "$work"
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build -trimpath -ldflags="$ldflags" -o "$work/$bin" "$pkg"
  cp README.md CHANGELOG.md "$work/"
  cp -R docs configs deployments "$work/"
  if [[ "$goos" == "windows" ]]; then
    (cd dist && zip -qr "${name}.zip" "$name")
  else
    tar -C dist -czf "dist/${name}.tar.gz" "$name"
  fi
  rm -rf "$work"
done

(cd dist && shasum -a 256 proxctl_${version}_* > "checksums.txt")

if [[ "$mode" == "snapshot" ]]; then
  mkdir -p build
  cp dist/proxctl_${version}_linux_amd64.tar.gz build/ 2>/dev/null || true
fi
