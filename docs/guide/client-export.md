# Client Export

Proxctl can generate client profiles from the current Xray config.

## Supported formats

```bash
proxctl client list
```

Current formats:

- `generic-uri`
- `shadowrocket`
- `surge`
- `stash`
- `mihomo`
- `sing-box`
- `v2rayn`
- `v2rayng`

## Shadowrocket

On a VPS, `proxctl` can read the local Xray config, derive the Reality public key from `privateKey`, and auto-detect the public server IP:

```bash
proxctl client export shadowrocket --name my-node
```

Use `--server` when clients should connect through a domain name instead of the detected public IP:

```bash
proxctl client export shadowrocket \
  --provider xray \
  --server <ip-or-domain> \
  --name my-node
```

Use `--public-key` only when the local config is unavailable or you want to override the derived key.

## Surge

```bash
proxctl client export surge \
  --name my-node
```

## Mihomo or Stash

```bash
proxctl client export mihomo \
  --name my-node
```

## Security note

Client exports include credentials such as UUID and shortId. Treat exported links and files as sensitive. Proxctl never prints the server Reality private key in client output.
