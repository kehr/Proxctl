# Release Notes

This page is generated at documentation build time from the repository `VERSION` file.

- Current version: `v0.2.2`
- Release: [GitHub v0.2.2](https://github.com/kehr/Proxctl/releases/tag/v0.2.2)

## Proxctl v0.2.2

Proxctl is a lightweight proxy deployment and operations CLI. This release provides versioned multi-platform binaries and operational tooling for Xray-based VPS nodes.

Documentation: https://proxctl.kaixuan.ai

## Highlights

- Xray lifecycle commands for install, init, adopt, status, audit, health, and doctor workflows.
- Transaction-oriented operations for backup, restore, credential rotation, and Reality target switching.
- Client config export for Shadowrocket, Surge, Mihomo/Stash, sing-box, v2rayN, and v2rayNG.
- Cobra/Viper command framework with shell completion and generated command reference.
- Versioned release artifacts for Linux, macOS, and Windows.

## Changes

- chore: release v0.2.2 (1101a85)
- fix: apply ssh hardening before vendor drop-ins (da8bbbe)

## Install

One-line installer:

```bash
curl -fsSL https://github.com/kehr/Proxctl/releases/latest/download/install.sh | sh
```

Linux amd64:

```bash
curl -LO https://github.com/kehr/Proxctl/releases/download/v0.2.2/proxctl-v0.2.2-linux-amd64.tar.gz
tar -xzf proxctl-v0.2.2-linux-amd64.tar.gz
sudo install -m 0755 proxctl-v0.2.2-linux-amd64/proxctl /usr/local/bin/proxctl
proxctl version
```

macOS arm64:

```bash
curl -LO https://github.com/kehr/Proxctl/releases/download/v0.2.2/proxctl-v0.2.2-darwin-arm64.tar.gz
tar -xzf proxctl-v0.2.2-darwin-arm64.tar.gz
sudo install -m 0755 proxctl-v0.2.2-darwin-arm64/proxctl /usr/local/bin/proxctl
proxctl version
```

## Assets

- `proxctl-v0.2.2-linux-amd64.tar.gz`
- `proxctl-v0.2.2-linux-arm64.tar.gz`
- `proxctl-v0.2.2-darwin-amd64.tar.gz`
- `proxctl-v0.2.2-darwin-arm64.tar.gz`
- `proxctl-v0.2.2-windows-amd64.zip`
- `checksums.txt`
- `install.sh`

## Verify Checksums

```bash
shasum -a 256 -c checksums.txt
```

## Full Changelog

https://github.com/kehr/Proxctl/compare/v0.2.1...v0.2.2
