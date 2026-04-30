# Configuration

Proxctl uses command flags, environment variables, config files, and built-in defaults.

Precedence:

```text
command flags -> environment variables -> config file -> built-in defaults
```

## Global flags

```bash
proxctl --help
```

Common flags:

- `--xray-config`
- `--state-dir`
- `--xray-bin`
- `--service`
- `--provider`
- `--yes`
- `--json`
- `--no-color`

## Environment variables

Viper binds environment variables with the `PROXCTL_` prefix.

Examples:

```bash
export PROXCTL_XRAY_CONFIG_PATH=/usr/local/etc/xray/config.json
export PROXCTL_OPERATIONS_STATE_DIR=/var/lib/proxctl
```

## Defaults file

Example defaults live at:

```text
configs/defaults.yaml
```

The production location is:

```text
/etc/proxctl/defaults.yaml
```
