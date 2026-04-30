package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kehr/proxctl/src/internal/client"
	"github.com/kehr/proxctl/src/internal/state"
	"github.com/kehr/proxctl/src/internal/system"
	"github.com/kehr/proxctl/src/internal/xray"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func newStatusCommand(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: "status", Short: "Show compact node status", RunE: func(cmd *cobra.Command, args []string) error {
		p := rt.Printer()
		p.Section("Status")
		sys := system.Collect(cmd.Context(), rt.Runner)
		p.Info("Host: %s", sys.Hostname)
		p.Info("Uptime: %s", sys.Uptime)
		cfg, err := rt.Xray().LoadConfig()
		if err != nil {
			p.Warn("Xray config not readable: %v", err)
			return nil
		}
		fmt.Fprint(rt.Out, cfg.Summary())
		return nil
	}}
}

func newAuditCommand(rt *Runtime) *cobra.Command {
	var skipUpdates bool
	c := &cobra.Command{Use: "audit", Short: "Run read-only system and proxy audit", RunE: func(cmd *cobra.Command, args []string) error {
		p := rt.Printer()
		p.Section("System")
		sys := system.Collect(cmd.Context(), rt.Runner)
		p.Info("Host: %s", sys.Hostname)
		p.Info("OS: %s", sys.OS)
		p.Info("Kernel: %s", sys.Kernel)
		p.Info("Memory: %s", strings.TrimSpace(sys.Memory))
		p.Info("Disk: %s", strings.TrimSpace(sys.Disk))
		p.Section("Xray")
		_ = runHealth(cmd, rt)
		p.Section("SSH")
		system.AuditSSH(cmd.Context(), rt.Runner, p)
		p.Section("Firewall")
		system.AuditFirewall(cmd.Context(), rt.Runner, p)
		p.Section("Updates")
		if skipUpdates {
			p.Info("Skipped package update check.")
		} else {
			system.CheckUpdates(cmd.Context(), rt.Runner, p)
		}
		return nil
	}}
	c.Flags().BoolVar(&skipUpdates, "skip-updates", false, "skip package update checks")
	return c
}

func newHealthCommand(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: "health", Short: "Run Xray health checks", RunE: func(cmd *cobra.Command, args []string) error {
		return runHealth(cmd, rt)
	}}
}

func runHealth(cmd *cobra.Command, rt *Runtime) error {
	p := rt.Printer()
	prov := rt.Xray()
	if _, err := os.Stat(prov.ConfigPath); err != nil {
		p.Fail("Config missing: %s", prov.ConfigPath)
		return err
	}
	p.Pass("Config exists: %s", prov.ConfigPath)
	if err := prov.TestConfig(cmd.Context()); err != nil {
		p.Fail("Xray config syntax failed: %v", err)
	} else {
		p.Pass("Xray config syntax OK.")
	}
	if prov.ServiceActive(cmd.Context()) {
		p.Pass("Service is active: %s", prov.Service)
	} else {
		p.Fail("Service is not active: %s", prov.Service)
	}
	ports, _ := prov.Ports()
	for _, port := range ports {
		if prov.PortOwnedByService(cmd.Context(), port) {
			p.Pass("Port %d is listened by %s.", port, prov.Service)
		} else {
			p.Fail("Port %d is not confirmed as %s.", port, prov.Service)
		}
	}
	return nil
}

func newDoctorCommand(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: "doctor", Short: "Run health checks and show recent Xray logs", RunE: func(cmd *cobra.Command, args []string) error {
		if err := runHealth(cmd, rt); err != nil {
			return err
		}
		rt.Printer().Section("Recent Xray Journal")
		res := rt.Runner.Run(cmd.Context(), "journalctl", "-u", rt.Config.Service, "-n", "80", "--no-pager")
		fmt.Fprint(rt.Out, res.Stdout)
		fmt.Fprint(rt.Err, res.Stderr)
		return nil
	}}
}

func newConfigCommand(rt *Runtime) *cobra.Command {
	root := &cobra.Command{Use: "config", Short: "Inspect provider configuration"}
	root.AddCommand(&cobra.Command{Use: "summary", Short: "Print redacted Xray config summary", RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := rt.Xray().LoadConfig()
		if err != nil {
			return err
		}
		fmt.Fprint(rt.Out, cfg.Summary())
		return nil
	}})
	return root
}

func newClientCommand(rt *Runtime) *cobra.Command {
	root := &cobra.Command{Use: "client", Short: "Generate client configuration"}
	root.AddCommand(&cobra.Command{Use: "list", Short: "List supported client export formats", Run: func(cmd *cobra.Command, args []string) {
		for _, f := range client.Formats() {
			fmt.Fprintln(rt.Out, f)
		}
	}})
	var provider, server, pub, name string
	export := &cobra.Command{Use: "export <format>", Short: "Export client config", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireProvider(provider); err != nil {
			return err
		}
		if server == "" {
			return fmt.Errorf("--server is required")
		}
		cfg, err := rt.Xray().LoadConfig()
		if err != nil {
			return err
		}
		profile, err := cfg.ClientProfile(server, pub, name)
		if err != nil {
			return err
		}
		text, err := client.Export(args[0], profile)
		if err != nil {
			return err
		}
		fmt.Fprintln(rt.Out, text)
		return nil
	}}
	export.Flags().StringVar(&provider, "provider", "xray", "provider")
	export.Flags().StringVar(&server, "server", "", "server address")
	export.Flags().StringVar(&pub, "public-key", "", "Reality public key")
	export.Flags().StringVar(&name, "name", "proxctl-node", "profile name")
	root.AddCommand(export)
	return root
}

func newBackupCommand(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: "backup [label]", Short: "Create a backup manifest", Args: cobra.MaximumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		label := "manual"
		if len(args) > 0 {
			label = args[0]
		}
		b, err := backupXray(cmd, rt, label)
		if err != nil {
			return err
		}
		rt.Printer().Pass("Backup created: %s", b.Path)
		return nil
	}}
}

func backupXray(cmd *cobra.Command, rt *Runtime, label string) (state.Backup, error) {
	prov := rt.Xray()
	return rt.State().Backup(label, state.BackupInput{
		ConfigPath:      prov.ConfigPath,
		XrayVersion:     prov.Version(cmd.Context()),
		ServiceUnitText: prov.ServiceUnit(cmd.Context()),
		SSHDEffective:   system.SSHDEffective(cmd.Context(), rt.Runner),
		CommandLine:     os.Args,
	})
}

func newRestoreCommand(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: "restore [id|latest]", Short: "Restore Xray config from backup", Args: cobra.MaximumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		id := "latest"
		if len(args) > 0 {
			id = args[0]
		}
		if err := rt.Prompt().Danger("restore xray config from backup "+id, "RESTORE"); err != nil {
			return err
		}
		src, err := rt.State().RestoreConfig(id, rt.Config.ConfigPath)
		if err != nil {
			return err
		}
		if err := rt.Xray().TestConfig(cmd.Context()); err != nil {
			return err
		}
		if err := rt.Xray().Restart(cmd.Context()); err != nil {
			return err
		}
		rt.Printer().Pass("Restored %s to %s", src, rt.Config.ConfigPath)
		return runHealth(cmd, rt)
	}}
}

func newAdoptCommand(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: "adopt xray", Short: "Adopt an existing Xray installation", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireProvider(args[0]); err != nil {
			return err
		}
		cfg, err := rt.Xray().LoadConfig()
		if err != nil {
			return err
		}
		fmt.Fprint(rt.Out, cfg.Summary())
		b, err := backupXray(cmd, rt, "baseline")
		if err != nil {
			return err
		}
		rt.Printer().Pass("Baseline backup created: %s", b.Path)
		return nil
	}}
}

func newInstallCommand(rt *Runtime) *cobra.Command {
	var apply bool
	c := &cobra.Command{Use: "install xray", Short: "Install provider binaries and service files", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireProvider(args[0]); err != nil {
			return err
		}
		p := rt.Printer()
		p.Section("Install Xray")
		p.Info("Method: official XTLS install script")
		p.Info("Will not overwrite existing config.")
		if !apply {
			return nil
		}
		if err := rt.Prompt().Danger("install xray using remote official script", "INSTALL"); err != nil {
			return err
		}
		res := rt.Runner.Run(cmd.Context(), "bash", "-lc", `curl -L https://github.com/XTLS/Xray-install/raw/main/install-release.sh | bash -s -- install`)
		if res.Code != 0 {
			return fmt.Errorf("xray install script failed: %s", res.Stderr)
		}
		return nil
	}}
	c.Flags().BoolVar(&apply, "apply", false, "apply install plan")
	c.Flags().Bool("plan", true, "show install plan")
	return c
}

func newInitCommand(rt *Runtime) *cobra.Command {
	var apply, force bool
	c := &cobra.Command{Use: "init xray", Short: "Initialize provider default config", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireProvider(args[0]); err != nil {
			return err
		}
		p := rt.Printer()
		p.Section("Init Xray")
		p.Info("Default protocol: VLESS + TCP + Reality")
		p.Info("Default port: 443")
		p.Info("Default log: warning, access log off")
		if !apply {
			return nil
		}
		if _, err := os.Stat(rt.Config.ConfigPath); err == nil && !force {
			return fmt.Errorf("config exists; use adopt xray or init xray --apply --force")
		}
		if err := rt.Prompt().Danger("initialize xray config", "INIT"); err != nil {
			return err
		}
		return rt.Xray().InitDefault(cmd.Context())
	}}
	c.Flags().BoolVar(&apply, "apply", false, "apply init plan")
	c.Flags().Bool("plan", true, "show init plan")
	c.Flags().BoolVar(&force, "force", false, "overwrite existing config")
	return c
}

func newPlanCommand(rt *Runtime) *cobra.Command {
	root := &cobra.Command{Use: "plan", Short: "Show planned high-risk changes"}
	root.AddCommand(newRotateCommand(rt, false), newSwitchCommand(rt, false))
	return root
}

func newApplyCommand(rt *Runtime) *cobra.Command {
	root := &cobra.Command{Use: "apply", Short: "Apply high-risk planned changes"}
	root.AddCommand(newRotateCommand(rt, true), newSwitchCommand(rt, true))
	return root
}

func newRotateCommand(rt *Runtime, apply bool) *cobra.Command {
	return &cobra.Command{Use: "rotate xray <uuid|shortid|reality-key|all>", Short: "Rotate Xray credentials", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireProvider(args[0]); err != nil {
			return err
		}
		target := args[1]
		p := rt.Printer()
		p.Section("Rotate Xray Credentials")
		p.Info("Target: %s", target)
		p.Info("Sequence: backup -> render -> xray test -> atomic replace -> restart -> healthcheck -> client export")
		if !apply {
			return nil
		}
		if err := rt.Prompt().Danger("rotate xray "+target, "ROTATE"); err != nil {
			return err
		}
		if _, err := backupXray(cmd, rt, "pre-rotate"); err != nil {
			return err
		}
		cfg, err := rt.Xray().LoadConfig()
		if err != nil {
			return err
		}
		_, meta, err := cfg.Rotate(xray.RotateTarget(target), xray.SystemGenerator{Runner: rt.Runner, XrayBin: rt.Config.XrayBin})
		if err != nil {
			return err
		}
		if err := rt.Xray().WriteConfig(cfg); err != nil {
			return err
		}
		if err := rt.Xray().TestConfig(cmd.Context()); err != nil {
			return err
		}
		if err := rt.Xray().Restart(cmd.Context()); err != nil {
			return err
		}
		_ = rt.State().SaveState(meta)
		return runHealth(cmd, rt)
	}}
}

func newSwitchCommand(rt *Runtime, apply bool) *cobra.Command {
	return &cobra.Command{Use: "switch xray reality-target <host:port>", Short: "Switch Xray Reality target", Args: cobra.ExactArgs(3), RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireProvider(args[0]); err != nil {
			return err
		}
		if args[1] != "reality-target" {
			return fmt.Errorf("unsupported switch key: %s", args[1])
		}
		target := args[2]
		p := rt.Printer()
		p.Section("Switch Reality Target")
		p.Info("Target: %s", target)
		if !apply {
			return nil
		}
		if err := rt.Prompt().Danger("switch reality target to "+target, "SWITCH"); err != nil {
			return err
		}
		if _, err := backupXray(cmd, rt, "pre-switch"); err != nil {
			return err
		}
		cfg, err := rt.Xray().LoadConfig()
		if err != nil {
			return err
		}
		if err := cfg.SetRealityTarget(target); err != nil {
			return err
		}
		if err := rt.Xray().WriteConfig(cfg); err != nil {
			return err
		}
		if err := rt.Xray().TestConfig(cmd.Context()); err != nil {
			return err
		}
		if err := rt.Xray().Restart(cmd.Context()); err != nil {
			return err
		}
		return runHealth(cmd, rt)
	}}
}

func newSSHCommand(rt *Runtime) *cobra.Command {
	root := &cobra.Command{Use: "ssh", Short: "Audit and harden SSH"}
	var apply bool
	harden := &cobra.Command{Use: "harden", Short: "Plan or apply SSH hardening", RunE: func(cmd *cobra.Command, args []string) error {
		p := rt.Printer()
		system.PlanSSHHarden(cmd.Context(), rt.Runner, p)
		if !apply {
			return nil
		}
		if err := rt.Prompt().Danger("ssh hardening", "SSH"); err != nil {
			return err
		}
		return system.ApplySSHHarden(cmd.Context(), rt.Runner)
	}}
	harden.Flags().BoolVar(&apply, "apply", false, "apply SSH hardening")
	harden.Flags().Bool("plan", true, "show SSH hardening plan")
	root.AddCommand(harden)
	return root
}

func newFirewallCommand(rt *Runtime) *cobra.Command {
	root := &cobra.Command{Use: "firewall", Short: "Audit and harden firewall"}
	var apply bool
	enable := &cobra.Command{Use: "enable", Short: "Plan or enable firewall", RunE: func(cmd *cobra.Command, args []string) error {
		p := rt.Printer()
		system.PlanFirewall(cmd.Context(), rt.Runner, p)
		if !apply {
			return nil
		}
		if err := rt.Prompt().Danger("enable firewall", "FIREWALL"); err != nil {
			return err
		}
		return system.ApplyFirewall(cmd.Context(), rt.Runner, rt.Config.ConfigPath)
	}}
	enable.Flags().BoolVar(&apply, "apply", false, "apply firewall plan")
	enable.Flags().Bool("plan", true, "show firewall plan")
	root.AddCommand(enable)
	return root
}

func newBootCheckCommand(rt *Runtime) *cobra.Command {
	var record bool
	c := &cobra.Command{Use: "boot-check", Short: "Run boot health check", RunE: func(cmd *cobra.Command, args []string) error {
		err := runHealth(cmd, rt)
		if record {
			logDir := filepath.Join(rt.Config.StateDir, "logs")
			_ = os.MkdirAll(logDir, 0700)
			_ = os.WriteFile(filepath.Join(logDir, "boot-check.log"), []byte("boot-check executed\n"), 0600)
		}
		return err
	}}
	c.Flags().BoolVar(&record, "record", false, "record boot-check marker")
	return c
}

func newDocsCommand(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: "docs <dir>", Short: "Generate Markdown command docs", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := os.MkdirAll(args[0], 0755); err != nil {
			return err
		}
		root := NewRootCommand(rt)
		root.DisableAutoGenTag = true
		return doc.GenMarkdownTree(root, args[0])
	}}
}

func newWizardCommand(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: "wizard", Short: "Interactive setup wizard", RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(rt.Out, "Interactive wizard is reserved for the TUI layer. Use proxctl init xray --plan for now.")
		return nil
	}}
}
