# Credential Rotation

Credential rotation is a high-risk operation because existing clients may stop working until they receive updated profiles.

## Review the plan

```bash
sudo proxctl plan rotate xray all
```

Targets:

- `uuid`
- `shortid`
- `reality-key`
- `all`

## Apply rotation

```bash
sudo proxctl apply rotate xray all
```

You must type:

```text
ROTATE
```

The intended operation sequence is:

```text
backup -> render config -> xray test -> atomic replace -> restart -> healthcheck -> client export
```

## Export updated clients

After rotation, export updated client configs:

```bash
proxctl client export shadowrocket \
  --provider xray \
  --name my-node
```

Pass `--server <domain>` if clients should use a domain instead of the auto-detected public IP.
