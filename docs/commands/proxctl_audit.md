## proxctl audit

Run read-only system and proxy audit

```
proxctl audit [flags]
```

### Options

```
  -h, --help           help for audit
      --skip-updates   skip package update checks
```

### Options inherited from parent commands

```
      --config-file string   proxctl defaults file
      --json                 emit JSON where supported
      --no-color             disable color
      --provider string      default provider (default "xray")
      --service string       Xray systemd service (default "xray")
      --state-dir string     state directory (default "/var/lib/proxctl")
      --xray-bin string      Xray binary (default "xray")
      --xray-config string   Xray config path (default "/usr/local/etc/xray/config.json")
      --yes                  accept low-risk defaults
```

### SEE ALSO

* [proxctl](proxctl.md)	 - Lightweight proxy deployment and operations CLI

