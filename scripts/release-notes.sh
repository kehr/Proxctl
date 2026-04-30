#!/usr/bin/env bash
set -euo pipefail

version="${1:-$(./scripts/version.sh release)}"
previous="${2:-}"

if [[ -z "$previous" ]]; then
  previous="$(git tag --sort=-creatordate | grep -E '^v[0-9]' | grep -v "^${version}$" | head -n 1 || true)"
fi

range=""
if [[ -n "$previous" ]]; then
  range="${previous}..${version}"
else
  range="${version}"
fi

cat <<EOF
## Proxctl ${version}

Proxctl is a lightweight proxy deployment and operations CLI. This release provides versioned multi-platform binaries and operational tooling for Xray-based VPS nodes.

Documentation: https://proxctl.kaixuan.ai

## Highlights

- Xray lifecycle commands for install, init, adopt, status, audit, health, and doctor workflows.
- Transaction-oriented operations for backup, restore, credential rotation, and Reality target switching.
- Client config export for Shadowrocket, Surge, Mihomo/Stash, sing-box, v2rayN, and v2rayNG.
- Cobra/Viper command framework with shell completion and generated command reference.
- Versioned release artifacts for Linux, macOS, and Windows.

## Changes

EOF

if git rev-parse "$version" >/dev/null 2>&1; then
  if [[ -n "$previous" ]]; then
    git log --pretty='- %s (%h)' "$range"
  else
    git log --pretty='- %s (%h)' "$version"
  fi
else
  git log --pretty='- %s (%h)' -n 20
fi

cat <<EOF

## Install

One-line installer:

\`\`\`bash
curl -fsSL https://github.com/kehr/Proxctl/releases/latest/download/install.sh | sh
\`\`\`

Linux amd64:

\`\`\`bash
curl -LO https://github.com/kehr/Proxctl/releases/download/${version}/proxctl-${version}-linux-amd64.tar.gz
tar -xzf proxctl-${version}-linux-amd64.tar.gz
sudo install -m 0755 proxctl-${version}-linux-amd64/proxctl /usr/local/bin/proxctl
proxctl version
\`\`\`

macOS arm64:

\`\`\`bash
curl -LO https://github.com/kehr/Proxctl/releases/download/${version}/proxctl-${version}-darwin-arm64.tar.gz
tar -xzf proxctl-${version}-darwin-arm64.tar.gz
sudo install -m 0755 proxctl-${version}-darwin-arm64/proxctl /usr/local/bin/proxctl
proxctl version
\`\`\`

## Assets

- \`proxctl-${version}-linux-amd64.tar.gz\`
- \`proxctl-${version}-linux-arm64.tar.gz\`
- \`proxctl-${version}-darwin-amd64.tar.gz\`
- \`proxctl-${version}-darwin-arm64.tar.gz\`
- \`proxctl-${version}-windows-amd64.zip\`
- \`checksums.txt\`

## Verify Checksums

\`\`\`bash
shasum -a 256 -c checksums.txt
\`\`\`

EOF

if [[ -n "$previous" ]]; then
  cat <<EOF
## Full Changelog

https://github.com/kehr/Proxctl/compare/${previous}...${version}
EOF
else
  cat <<EOF
## Full Changelog

https://github.com/kehr/Proxctl/commits/${version}
EOF
fi
