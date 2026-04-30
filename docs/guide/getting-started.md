# Getting Started

Proxctl is a local operations CLI for proxy services on small VPS nodes. It starts with Xray and is designed so future providers can reuse the same lifecycle model.

## Install Proxctl

```bash
curl -fsSL https://github.com/kehr/Proxctl/releases/latest/download/install.sh | sh
```

Verify the binary:

```bash
proxctl version
```

## Adopt an existing Xray node

If Xray is already installed, adopt it before making changes:

```bash
sudo proxctl adopt xray
```

This reads the current config, prints a redacted summary, and creates a baseline backup under `/var/lib/proxctl/backups`.

## Run health checks

```bash
sudo proxctl health
```

The health check validates:

- config file exists
- `xray run -test -config` succeeds
- `xray.service` is active
- configured inbound ports are owned by Xray

## Audit the node

```bash
sudo proxctl audit --skip-updates
```

Use `--skip-updates` on low-resource nodes when you want to avoid package metadata checks.

## Export a client config

```bash
proxctl client export shadowrocket \
  --provider xray \
  --server <ip-or-domain> \
  --public-key <reality-public-key> \
  --name my-node
```

Use `proxctl client list` to see supported formats.

## Review before applying changes

High-risk changes use plan/apply commands:

```bash
sudo proxctl plan rotate xray all
sudo proxctl apply rotate xray all
```

Apply commands require confirmation tokens such as `ROTATE`, `RESTORE`, `SSH`, or `FIREWALL`.
