# Proxctl

Proxctl is a lightweight CLI for deploying and operating proxy services on small VPS nodes.

The first supported provider is Xray. The project is designed to expand to other proxy backends over time.

## Features

- Xray install/init/adopt lifecycle
- Status, audit, health, and doctor checks
- Backup, restore, and safer change workflows
- Credential rotation and Reality target switching
- Client config export for Shadowrocket, Surge, Mihomo/Stash, sing-box, v2rayN, and v2rayNG
- SSH and firewall hardening plans
- Shell completion and generated command reference

## Install

One-line installer:

```bash
curl -fsSL https://raw.githubusercontent.com/kehr/Proxctl/main/install.sh | sh
```

Install a specific version:

```bash
curl -fsSL https://raw.githubusercontent.com/kehr/Proxctl/main/install.sh | PROXCTL_VERSION=v0.2.0 sh
```

Download the archive for your platform from GitHub Releases, then install the binary:

```bash
tar -xzf proxctl-v0.2.0-linux-amd64.tar.gz
sudo install -m 0755 proxctl-v0.2.0-linux-amd64/proxctl /usr/local/bin/proxctl
```

Verify:

```bash
proxctl version
```

## Quick Start

Adopt an existing Xray installation:

```bash
sudo proxctl adopt xray
sudo proxctl health
sudo proxctl audit --skip-updates
```

Export a client profile:

```bash
proxctl client export shadowrocket \
  --provider xray \
  --server <ip-or-domain> \
  --public-key <reality-public-key> \
  --name my-node
```

Review high-risk changes before applying:

```bash
sudo proxctl plan rotate xray all
sudo proxctl apply rotate xray all
```

## Documentation

- [Operations guide](docs/operations.md)
- [Architecture](docs/architecture.md)
- [Development guide](docs/development.md)
- [Command reference](docs/commands/proxctl.md)
