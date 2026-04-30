# Architecture

Proxctl is a production-oriented local control plane for proxy services on small VPS nodes.

## Design Principles

- Default to read-only commands for routine operations.
- Treat configuration changes as transactions: backup, render, validate, apply, restart, health check.
- Keep runtime dependencies low. The Go binary uses the standard library and delegates only to system tools already required for operations, such as `systemctl`, `ss`, `sshd`, and provider binaries.
- Keep provider code isolated. Xray is implemented first, but the command model is provider-neutral.
- Do not hide high-risk changes behind `--yes`. SSH, firewall, restore, and credential rotation require explicit confirmation tokens.

## Top-Level Directories

- `src/cmd/proxctl`: binary entrypoint.
- `src/internal/cli`: Cobra/Viper command routing, config binding, completion, and docs generation.
- `src/internal/xray`: Xray config parsing, rendering, health, init, rotate, and service integration.
- `src/internal/client`: client profile exporters.
- `src/internal/state`: backup, manifest, state, and restore helpers.
- `src/internal/system`: OS, SSH, firewall, and update checks.
- `src/internal/command`: external command execution boundary.
- `docs`: operator and implementation documentation.
- `configs`: default production configuration examples.
- `scripts`: build and deployment helpers.

## Lifecycle

```text
new machine:
  install xray -> init xray -> health -> client export

existing machine:
  adopt xray -> backup baseline -> audit -> plan/apply changes

ongoing operations:
  backup -> rotate/switch/harden -> health -> rollback when needed
```

## Provider Boundary

Providers should implement these conceptual operations:

- detect installation
- load config
- summarize config without leaking secrets
- test config
- start/restart service
- inspect listening ports
- init default config
- rotate credentials
- export client profile data

## CLI Framework

Proxctl uses Cobra for the command tree and shell completion, and Viper for default configuration, environment variable binding, and future `/etc/proxctl/defaults.yaml` loading.

Configuration precedence is:

```text
command flags -> environment variables -> config file -> built-in defaults
```

Environment variables use the `PROXCTL_` prefix, for example `PROXCTL_XRAY_CONFIG_PATH`.
