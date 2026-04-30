## proxctl client export

Export client config

```
proxctl client export <format> [flags]
```

### Options

```
  -h, --help                help for export
      --name string         profile name (default "proxctl-node")
      --provider string     provider (default "xray")
      --public-key string   Reality public key; derives from config privateKey when omitted
      --server string       server address; auto-detects public IP when omitted
```

### Options inherited from parent commands

```
      --config-file string   proxctl defaults file
      --json                 emit JSON where supported
      --no-color             disable color
      --service string       Xray systemd service (default "xray")
      --state-dir string     state directory (default "/var/lib/proxctl")
      --xray-bin string      Xray binary (default "xray")
      --xray-config string   Xray config path (default "/usr/local/etc/xray/config.json")
      --yes                  accept low-risk defaults
```

### SEE ALSO

* [proxctl client](proxctl_client.md)	 - Generate client configuration

