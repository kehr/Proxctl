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

```bash
proxctl client export shadowrocket \
  --provider xray \
  --server <ip-or-domain> \
  --public-key <reality-public-key> \
  --name my-node
```

## Surge

```bash
proxctl client export surge \
  --provider xray \
  --server <ip-or-domain> \
  --public-key <reality-public-key> \
  --name my-node
```

## Mihomo or Stash

```bash
proxctl client export mihomo \
  --provider xray \
  --server <ip-or-domain> \
  --public-key <reality-public-key> \
  --name my-node
```

## Security note

Client exports include credentials such as UUID and shortId. Treat exported links and files as sensitive. Proxctl never prints the server Reality private key in client output.
