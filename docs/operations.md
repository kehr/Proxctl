# Operations Guide

## Adopt Existing Xray

```bash
proxctl adopt xray
```

This reads the current Xray config, prints a redacted summary, and creates a baseline backup under `/var/lib/proxctl/backups`.

## Health Check

```bash
proxctl health
```

Checks:

- config file exists
- `xray run -test -config` passes
- `xray.service` is active
- inbound ports are owned by the Xray service

## Credential Rotation

Always review first:

```bash
proxctl plan rotate xray all
```

Apply requires typing `ROTATE`:

```bash
proxctl apply rotate xray all
```

After rotation, export updated client config:

```bash
proxctl client export shadowrocket --provider xray --server <ip-or-domain> --public-key <new-public-key>
```

## Restore

Restore requires typing `RESTORE`:

```bash
proxctl restore latest
```

The restored config is syntax-tested before the service is restarted.

## SSH Hardening

Plan first:

```bash
proxctl ssh harden --plan
```

Apply requires typing `SSH`:

```bash
proxctl ssh harden --apply
```

Before applying SSH hardening, keep an existing SSH session open and verify key-based login works in a second session.

