# Xray Provider

Xray is the first provider implemented by Proxctl. The provider manages the local Xray binary, systemd service, and `/usr/local/etc/xray/config.json`.

## Lifecycle commands

### Install

```bash
sudo proxctl install xray --plan
sudo proxctl install xray --apply
```

The install flow is intentionally explicit. It uses the official XTLS install script and does not overwrite existing business config.

### Init

```bash
sudo proxctl init xray --plan
sudo proxctl init xray --apply
```

`init` creates a default VLESS + TCP + Reality config for a new node.

Defaults:

- port `443`
- `loglevel: warning`
- access log disabled
- VLESS Vision flow
- Reality target `www.microsoft.com:443`

If a config already exists, use `adopt xray` instead. Overwrite requires `--force`.

### Adopt

```bash
sudo proxctl adopt xray
```

Adoption is the correct path for existing servers. It reads current config, prints a redacted summary, and creates a baseline backup.

## Config summary

```bash
proxctl config summary
```

Sensitive fields such as Reality `privateKey` are redacted.

## Health

```bash
sudo proxctl health
```

Health checks verify systemd state, config syntax, and port ownership.
