## proxctl plan rotate

Rotate Xray credentials

```
proxctl plan rotate xray <uuid|shortid|reality-key|all> [flags]
```

### Options

```
  -h, --help   help for rotate
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

* [proxctl plan](proxctl_plan.md)	 - Show planned high-risk changes

