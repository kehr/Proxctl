## proxctl

Lightweight proxy deployment and operations CLI

### Options

```
      --config-file string   proxctl defaults file
  -h, --help                 help for proxctl
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

* [proxctl adopt](proxctl_adopt.md)	 - Adopt an existing Xray installation
* [proxctl apply](proxctl_apply.md)	 - Apply high-risk planned changes
* [proxctl audit](proxctl_audit.md)	 - Run read-only system and proxy audit
* [proxctl backup](proxctl_backup.md)	 - Create a backup manifest
* [proxctl boot-check](proxctl_boot-check.md)	 - Run boot health check
* [proxctl client](proxctl_client.md)	 - Generate client configuration
* [proxctl config](proxctl_config.md)	 - Inspect provider configuration
* [proxctl docs](proxctl_docs.md)	 - Generate Markdown command docs
* [proxctl doctor](proxctl_doctor.md)	 - Run health checks and show recent Xray logs
* [proxctl firewall](proxctl_firewall.md)	 - Audit and harden firewall
* [proxctl health](proxctl_health.md)	 - Run Xray health checks
* [proxctl init](proxctl_init.md)	 - Initialize provider default config
* [proxctl install](proxctl_install.md)	 - Install provider binaries and service files
* [proxctl plan](proxctl_plan.md)	 - Show planned high-risk changes
* [proxctl restore](proxctl_restore.md)	 - Restore Xray config from backup
* [proxctl ssh](proxctl_ssh.md)	 - Audit and harden SSH
* [proxctl status](proxctl_status.md)	 - Show compact node status
* [proxctl wizard](proxctl_wizard.md)	 - Interactive setup wizard

