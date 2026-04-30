# Release Notes

This page is generated at documentation build time from the repository `VERSION` file.

- Current version: `v0.2.4`
- Release: [GitHub v0.2.4](https://github.com/kehr/Proxctl/releases/tag/v0.2.4)

## Proxctl v0.2.4

Proxctl is a lightweight proxy deployment and operations CLI. This release provides versioned multi-platform binaries and operational tooling for Xray-based VPS nodes.

Documentation: https://proxctl.kaixuan.ai

## Highlights

- Xray lifecycle commands for install, init, adopt, status, audit, health, and doctor workflows.
- Transaction-oriented operations for backup, restore, credential rotation, and Reality target switching.
- Client config export for Shadowrocket, Surge, Mihomo/Stash, sing-box, v2rayN, and v2rayNG.
- Cobra/Viper command framework with shell completion and generated command reference.
- Versioned release artifacts for Linux, macOS, and Windows.

## Changes

- fix: parse current xray reality key output (35f9c66)
- docs: update release metadata for v0.2.3 (1e3c6c6)
- chore: release v0.2.3 (a8a481b)
- feat: auto-detect client export defaults (9933af3)
- fix: keep docs version metadata stable (548572b)
- ci: use node 24 for docs actions (38d7e69)
- docs: generate versioned release metadata (ba4e798)
- chore: release v0.2.2 (1101a85)
- fix: apply ssh hardening before vendor drop-ins (da8bbbe)
- fix: publish installer as release asset (edb66f1)
- fix: preserve installer archive path (aa7ebfb)
- chore: release v0.2.1 (9188a03)
- docs: add repository agent guidelines (7dd89fc)
- fix: allow untracked files during release (c6b991d)
- chore: add make release workflow (10a93e3)
- fix: harden operational config changes (fa69604)
- docs: add VitePress documentation site (02d89c8)
- chore: move installer into scripts (617a218)
- feat: add one-line installer (4b50051)
- ci: generate structured release notes (0370251)

## Install

One-line installer:

```bash
curl -fsSL https://github.com/kehr/Proxctl/releases/latest/download/install.sh | sh
```

Linux amd64:

```bash
curl -LO https://github.com/kehr/Proxctl/releases/download/v0.2.4/proxctl-v0.2.4-linux-amd64.tar.gz
tar -xzf proxctl-v0.2.4-linux-amd64.tar.gz
sudo install -m 0755 proxctl-v0.2.4-linux-amd64/proxctl /usr/local/bin/proxctl
proxctl version
```

macOS arm64:

```bash
curl -LO https://github.com/kehr/Proxctl/releases/download/v0.2.4/proxctl-v0.2.4-darwin-arm64.tar.gz
tar -xzf proxctl-v0.2.4-darwin-arm64.tar.gz
sudo install -m 0755 proxctl-v0.2.4-darwin-arm64/proxctl /usr/local/bin/proxctl
proxctl version
```

## Assets

- `proxctl-v0.2.4-linux-amd64.tar.gz`
- `proxctl-v0.2.4-linux-arm64.tar.gz`
- `proxctl-v0.2.4-darwin-amd64.tar.gz`
- `proxctl-v0.2.4-darwin-arm64.tar.gz`
- `proxctl-v0.2.4-windows-amd64.zip`
- `checksums.txt`
- `install.sh`

## Verify Checksums

```bash
shasum -a 256 -c checksums.txt
```

## Full Changelog

https://github.com/kehr/Proxctl/compare/v0.2.3...v0.2.4
