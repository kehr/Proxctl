# Boot Check

Boot checks verify that the proxy service recovered after a reboot.

## Run manually

```bash
sudo proxctl boot-check
```

## Record a boot check

```bash
sudo proxctl boot-check --record
```

The current implementation writes a boot-check marker under the Proxctl state directory.

## systemd timer

Systemd templates are provided under:

```text
deployments/systemd/
```

Install them when you want a boot-time health check:

```bash
sudo install -m 0644 deployments/systemd/proxctl-boot-check.service /etc/systemd/system/
sudo install -m 0644 deployments/systemd/proxctl-boot-check.timer /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now proxctl-boot-check.timer
```
