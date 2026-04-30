# Backup and Restore

Proxctl stores backups under `/var/lib/proxctl/backups`.

## Create a backup

```bash
sudo proxctl backup baseline
```

Each backup stores:

- Xray config
- manifest metadata
- Xray service text when available
- effective SSH config when available

## Restore latest

```bash
sudo proxctl restore latest
```

Restore requires typing:

```text
RESTORE
```

After restoring the config, Proxctl runs an Xray config test before restarting the service.

## Restore a specific backup

```bash
sudo proxctl restore <backup-id>
```

Backup IDs are timestamped and visible under:

```bash
sudo ls -1 /var/lib/proxctl/backups
```
