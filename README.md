# Proxctl

Proxctl is a lightweight proxy service deployment and operations CLI.

The first provider is Xray, with production-oriented workflows for single-purpose VPS nodes:

- install/init/adopt lifecycle
- read-only status, audit, health, and doctor checks
- transactional backup and restore
- credential rotation plans and apply flows
- Reality target switching
- client export for Shadowrocket, Surge, Mihomo/Stash, sing-box, v2rayN, and v2rayNG
- SSH and firewall hardening plans
- boot health checks

## Layout

```text
src/cmd/proxctl/       CLI entrypoint
src/internal/          private implementation packages
configs/              default configuration examples
docs/                 operator and architecture docs
scripts/              build and deployment helpers
build/                local build output
```

## Build

```bash
make test
make build
```

The binary is written to `build/proxctl`.

## Common Usage

```bash
proxctl status
proxctl audit --skip-updates
proxctl health
proxctl config summary
proxctl adopt xray
proxctl backup baseline
proxctl client export shadowrocket --provider xray --server <ip-or-domain> --public-key <reality-public-key>
proxctl plan rotate xray all
proxctl apply rotate xray all
proxctl ssh harden --plan
proxctl firewall enable --plan
```

High-risk actions require an explicit confirmation token even when `--yes` is set.

## State

Default production paths:

```text
/usr/local/bin/proxctl
/etc/proxctl/defaults.yaml
/var/lib/proxctl/
/var/lib/proxctl/backups/
```

## Provider Roadmap

The provider boundary is designed so Xray is only the first backend. Future providers can add sing-box, Hysteria, TUIC, or other proxy services without changing the CLI lifecycle model.

