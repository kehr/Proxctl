## proxctl firewall enable

Plan or enable firewall

```
proxctl firewall enable [flags]
```

### Options

```
      --apply   apply firewall plan
  -h, --help    help for enable
      --plan    show firewall plan (default true)
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

* [proxctl firewall](proxctl_firewall.md)	 - Audit and harden firewall

