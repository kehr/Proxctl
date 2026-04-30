# Firewall

Firewall changes are high-risk on remote VPS nodes. Proxctl treats firewall enablement as an explicit plan/apply operation.

## Review firewall state

```bash
sudo proxctl firewall enable --plan
```

The plan audits firewalld and nftables state and describes what will be changed.

## Enable firewall

```bash
sudo proxctl firewall enable --apply
```

You must type:

```text
FIREWALL
```

The firewall flow should always allow SSH and proxy ports before enabling enforcement.

## Recommendation

For single-purpose VPS nodes, SSH hardening and config permissions are usually higher priority than enabling a new firewall backend. Use the firewall command after confirming your provider console or recovery path.
