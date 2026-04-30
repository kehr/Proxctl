package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kehr/proxctl/src/internal/client"
	"github.com/kehr/proxctl/src/internal/command"
	"github.com/kehr/proxctl/src/internal/output"
	"github.com/kehr/proxctl/src/internal/prompt"
	"github.com/kehr/proxctl/src/internal/state"
	"github.com/kehr/proxctl/src/internal/system"
	"github.com/kehr/proxctl/src/internal/xray"
)

const Version = "0.1.0"

type App struct {
	Out     io.Writer
	Err     io.Writer
	In      io.Reader
	Runner  command.Runner
	Options Options
}

type Options struct {
	ConfigPath string
	StateDir   string
	XrayBin    string
	Service    string
	Yes        bool
	JSON       bool
	NoColor    bool
}

func New(out, err io.Writer, in io.Reader) *App {
	return &App{
		Out:    out,
		Err:    err,
		In:     in,
		Runner: command.LocalRunner{},
		Options: Options{
			ConfigPath: "/usr/local/etc/xray/config.json",
			StateDir:   "/var/lib/proxctl",
			XrayBin:    "xray",
			Service:    "xray",
		},
	}
}

func (a *App) Run(ctx context.Context, args []string) error {
	args = a.parseGlobal(args)
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		a.usage()
		return nil
	}
	p := &output.Printer{Out: a.Out, Err: a.Err, JSON: a.Options.JSON, Color: !a.Options.NoColor}
	cmd := args[0]
	args = args[1:]

	switch cmd {
	case "version":
		fmt.Fprintln(a.Out, "proxctl "+Version)
	case "status":
		return a.status(ctx, p)
	case "audit":
		return a.audit(ctx, p, args)
	case "health":
		return a.health(ctx, p)
	case "doctor":
		return a.doctor(ctx, p)
	case "config":
		return a.config(ctx, p, args)
	case "client":
		return a.client(ctx, args)
	case "backup":
		label := "manual"
		if len(args) > 0 {
			label = args[0]
		}
		return a.backup(ctx, p, label)
	case "restore":
		id := "latest"
		if len(args) > 0 {
			id = args[0]
		}
		return a.restore(ctx, p, id)
	case "adopt":
		return a.adopt(ctx, p)
	case "install":
		return a.install(ctx, p, args)
	case "init":
		return a.initXray(ctx, p, args)
	case "plan":
		return a.plan(ctx, p, args)
	case "apply":
		return a.apply(ctx, p, args)
	case "ssh":
		return a.ssh(ctx, p, args)
	case "firewall":
		return a.firewall(ctx, p, args)
	case "boot-check":
		return a.bootCheck(ctx, p, args)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
	return nil
}

func (a *App) parseGlobal(args []string) []string {
	out := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--yes", "-y":
			a.Options.Yes = true
		case "--json":
			a.Options.JSON = true
		case "--no-color":
			a.Options.NoColor = true
		case "--config":
			if i+1 < len(args) {
				i++
				a.Options.ConfigPath = args[i]
			}
		case "--state-dir":
			if i+1 < len(args) {
				i++
				a.Options.StateDir = args[i]
			}
		case "--xray-bin":
			if i+1 < len(args) {
				i++
				a.Options.XrayBin = args[i]
			}
		default:
			out = append(out, args[i:]...)
			return out
		}
	}
	return out
}

func (a *App) usage() {
	fmt.Fprint(a.Out, `proxctl - lightweight proxy deployment and operations CLI

Usage:
  proxctl status
  proxctl audit [--skip-updates]
  proxctl health
  proxctl doctor
  proxctl config summary
  proxctl client list
  proxctl client export <format> --provider xray --server <host> --public-key <key> [--name name]
  proxctl backup [label]
  proxctl restore [id|latest]
  proxctl adopt xray
  proxctl install xray --plan|--apply
  proxctl init xray --plan|--apply [--force]
  proxctl plan rotate xray <uuid|shortid|reality-key|all>
  proxctl apply rotate xray <uuid|shortid|reality-key|all>
  proxctl plan switch xray reality-target <host:port>
  proxctl apply switch xray reality-target <host:port>
  proxctl ssh harden --plan|--apply
  proxctl firewall enable --plan|--apply
  proxctl boot-check [--record]

Global:
  --config <path>      Xray config path
  --state-dir <path>   State directory
  --xray-bin <path>    Xray binary
  --yes, -y            Accept low-risk defaults
  --json               JSON output where supported
  --no-color           Disable color
`)
}

func (a *App) provider() xray.Provider {
	return xray.Provider{
		ConfigPath: a.Options.ConfigPath,
		Service:    a.Options.Service,
		Bin:        a.Options.XrayBin,
		Runner:     a.Runner,
	}
}

func (a *App) state() state.Manager { return state.Manager{StateDir: a.Options.StateDir} }
func (a *App) prompt() prompt.Prompter {
	return prompt.Prompter{In: a.In, Out: a.Out, AssumeYes: a.Options.Yes}
}

func (a *App) status(ctx context.Context, p *output.Printer) error {
	p.Section("Status")
	sys := system.Collect(ctx, a.Runner)
	p.Info("Host: %s", sys.Hostname)
	p.Info("Uptime: %s", sys.Uptime)
	cfg, err := a.provider().LoadConfig()
	if err != nil {
		p.Warn("Xray config not readable: %v", err)
		return nil
	}
	fmt.Fprint(a.Out, cfg.Summary())
	return nil
}

func (a *App) audit(ctx context.Context, p *output.Printer, args []string) error {
	skipUpdates := false
	for _, arg := range args {
		if arg == "--skip-updates" {
			skipUpdates = true
		}
	}
	p.Section("System")
	sys := system.Collect(ctx, a.Runner)
	p.Info("Host: %s", sys.Hostname)
	p.Info("OS: %s", sys.OS)
	p.Info("Kernel: %s", sys.Kernel)
	p.Info("Memory: %s", strings.TrimSpace(sys.Memory))
	p.Info("Disk: %s", strings.TrimSpace(sys.Disk))
	p.Section("Xray")
	if err := a.health(ctx, p); err != nil {
		p.Warn("%v", err)
	}
	p.Section("SSH")
	system.AuditSSH(ctx, a.Runner, p)
	p.Section("Firewall")
	system.AuditFirewall(ctx, a.Runner, p)
	p.Section("Updates")
	if skipUpdates {
		p.Info("Skipped package update check.")
	} else {
		system.CheckUpdates(ctx, a.Runner, p)
	}
	return nil
}

func (a *App) health(ctx context.Context, p *output.Printer) error {
	prov := a.provider()
	if _, err := os.Stat(prov.ConfigPath); err != nil {
		p.Fail("Config missing: %s", prov.ConfigPath)
		return err
	}
	p.Pass("Config exists: %s", prov.ConfigPath)
	if err := prov.TestConfig(ctx); err != nil {
		p.Fail("Xray config syntax failed: %v", err)
	} else {
		p.Pass("Xray config syntax OK.")
	}
	if prov.ServiceActive(ctx) {
		p.Pass("Service is active: %s", prov.Service)
	} else {
		p.Fail("Service is not active: %s", prov.Service)
	}
	ports, _ := prov.Ports()
	for _, port := range ports {
		if prov.PortOwnedByService(ctx, port) {
			p.Pass("Port %d is listened by %s.", port, prov.Service)
		} else {
			p.Fail("Port %d is not confirmed as %s.", port, prov.Service)
		}
	}
	return nil
}

func (a *App) doctor(ctx context.Context, p *output.Printer) error {
	if err := a.health(ctx, p); err != nil {
		return err
	}
	p.Section("Recent Xray Journal")
	res := a.Runner.Run(ctx, "journalctl", "-u", a.Options.Service, "-n", "80", "--no-pager")
	fmt.Fprint(a.Out, res.Stdout)
	fmt.Fprint(a.Err, res.Stderr)
	return nil
}

func (a *App) config(ctx context.Context, p *output.Printer, args []string) error {
	if len(args) == 0 || args[0] != "summary" {
		return errors.New("usage: proxctl config summary")
	}
	cfg, err := a.provider().LoadConfig()
	if err != nil {
		return err
	}
	fmt.Fprint(a.Out, cfg.Summary())
	return nil
}

func (a *App) client(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("usage: proxctl client list|export")
	}
	switch args[0] {
	case "list":
		for _, f := range client.Formats() {
			fmt.Fprintln(a.Out, f)
		}
		return nil
	case "export":
		fs := flag.NewFlagSet("client export", flag.ContinueOnError)
		fs.SetOutput(a.Err)
		provider := fs.String("provider", "xray", "provider")
		server := fs.String("server", "", "server address")
		pub := fs.String("public-key", "", "Reality public key")
		name := fs.String("name", "proxctl-node", "profile name")
		if len(args) < 2 {
			return errors.New("missing export format")
		}
		format := args[1]
		if err := fs.Parse(args[2:]); err != nil {
			return err
		}
		if *provider != "xray" {
			return fmt.Errorf("unsupported provider: %s", *provider)
		}
		if *server == "" {
			return errors.New("--server is required")
		}
		cfg, err := a.provider().LoadConfig()
		if err != nil {
			return err
		}
		profile, err := cfg.ClientProfile(*server, *pub, *name)
		if err != nil {
			return err
		}
		text, err := client.Export(format, profile)
		if err != nil {
			return err
		}
		fmt.Fprintln(a.Out, text)
		return nil
	default:
		return fmt.Errorf("unknown client command: %s", args[0])
	}
}

func (a *App) backup(ctx context.Context, p *output.Printer, label string) error {
	prov := a.provider()
	b, err := a.state().Backup(label, state.BackupInput{
		ConfigPath:      prov.ConfigPath,
		XrayVersion:     prov.Version(ctx),
		ServiceUnitText: prov.ServiceUnit(ctx),
		SSHDEffective:   system.SSHDEffective(ctx, a.Runner),
		CommandLine:     os.Args,
	})
	if err != nil {
		return err
	}
	p.Pass("Backup created: %s", b.Path)
	return nil
}

func (a *App) restore(ctx context.Context, p *output.Printer, id string) error {
	if err := a.prompt().Danger("restore xray config from backup "+id, "RESTORE"); err != nil {
		return err
	}
	src, err := a.state().RestoreConfig(id, a.Options.ConfigPath)
	if err != nil {
		return err
	}
	if err := a.provider().TestConfig(ctx); err != nil {
		return err
	}
	if err := a.provider().Restart(ctx); err != nil {
		return err
	}
	p.Pass("Restored %s to %s", src, a.Options.ConfigPath)
	return a.health(ctx, p)
}

func (a *App) adopt(ctx context.Context, p *output.Printer) error {
	if len(os.Args) > 0 && len(os.Args) > 2 && os.Args[2] != "xray" {
	}
	cfg, err := a.provider().LoadConfig()
	if err != nil {
		return err
	}
	fmt.Fprint(a.Out, cfg.Summary())
	return a.backup(ctx, p, "baseline")
}

func (a *App) install(ctx context.Context, p *output.Printer, args []string) error {
	if len(args) == 0 || args[0] != "xray" {
		return errors.New("usage: proxctl install xray --plan|--apply")
	}
	apply := contains(args, "--apply")
	p.Section("Install Xray")
	p.Info("Provider: xray")
	p.Info("Method: official XTLS install script")
	p.Info("Will not overwrite existing config unless init --force is used.")
	if !apply {
		return nil
	}
	if err := a.prompt().Danger("install xray using remote official script", "INSTALL"); err != nil {
		return err
	}
	if a.Runner.Run(ctx, "bash", "-lc", `curl -L https://github.com/XTLS/Xray-install/raw/main/install-release.sh | bash -s -- install`).Code != 0 {
		return errors.New("xray install script failed")
	}
	return nil
}

func (a *App) initXray(ctx context.Context, p *output.Printer, args []string) error {
	if len(args) == 0 || args[0] != "xray" {
		return errors.New("usage: proxctl init xray --plan|--apply")
	}
	force := contains(args, "--force")
	apply := contains(args, "--apply")
	p.Section("Init Xray")
	p.Info("Default protocol: VLESS + TCP + Reality")
	p.Info("Default port: 443")
	p.Info("Default log: warning, access log off")
	if !apply {
		return nil
	}
	if _, err := os.Stat(a.Options.ConfigPath); err == nil && !force {
		return errors.New("config exists; use adopt xray or init xray --apply --force")
	}
	if err := a.prompt().Danger("initialize xray config", "INIT"); err != nil {
		return err
	}
	return a.provider().InitDefault(ctx)
}

func (a *App) plan(ctx context.Context, p *output.Printer, args []string) error {
	if len(args) < 3 {
		return errors.New("usage: proxctl plan rotate xray <target> | plan switch xray reality-target <host:port>")
	}
	switch args[0] {
	case "rotate":
		return a.rotate(ctx, p, args[2], false)
	case "switch":
		if len(args) < 4 || args[2] != "reality-target" {
			return errors.New("usage: proxctl plan switch xray reality-target <host:port>")
		}
		return a.switchReality(ctx, p, args[3], false)
	default:
		return fmt.Errorf("unknown plan subject: %s", args[0])
	}
}

func (a *App) apply(ctx context.Context, p *output.Printer, args []string) error {
	if len(args) < 3 {
		return errors.New("usage: proxctl apply rotate xray <target> | apply switch xray reality-target <host:port>")
	}
	switch args[0] {
	case "rotate":
		return a.rotate(ctx, p, args[2], true)
	case "switch":
		if len(args) < 4 || args[2] != "reality-target" {
			return errors.New("usage: proxctl apply switch xray reality-target <host:port>")
		}
		return a.switchReality(ctx, p, args[3], true)
	default:
		return fmt.Errorf("unknown apply subject: %s", args[0])
	}
}

func (a *App) rotate(ctx context.Context, p *output.Printer, target string, apply bool) error {
	p.Section("Rotate Xray Credentials")
	p.Info("Target: %s", target)
	p.Info("Sequence: backup -> render -> xray test -> atomic replace -> restart -> healthcheck -> client export")
	if !apply {
		return nil
	}
	if err := a.prompt().Danger("rotate xray "+target, "ROTATE"); err != nil {
		return err
	}
	if err := a.backup(ctx, p, "pre-rotate"); err != nil {
		return err
	}
	cfg, err := a.provider().LoadConfig()
	if err != nil {
		return err
	}
	_, meta, err := cfg.Rotate(xray.RotateTarget(target), xray.SystemGenerator{Runner: a.Runner, XrayBin: a.Options.XrayBin})
	if err != nil {
		return err
	}
	if err := a.provider().WriteConfig(cfg); err != nil {
		return err
	}
	if err := a.provider().TestConfig(ctx); err != nil {
		return err
	}
	if err := a.provider().Restart(ctx); err != nil {
		return err
	}
	_ = a.state().SaveState(meta)
	return a.health(ctx, p)
}

func (a *App) switchReality(ctx context.Context, p *output.Printer, target string, apply bool) error {
	p.Section("Switch Reality Target")
	p.Info("Target: %s", target)
	if !apply {
		return nil
	}
	if err := a.prompt().Danger("switch reality target to "+target, "SWITCH"); err != nil {
		return err
	}
	if err := a.backup(ctx, p, "pre-switch"); err != nil {
		return err
	}
	cfg, err := a.provider().LoadConfig()
	if err != nil {
		return err
	}
	if err := cfg.SetRealityTarget(target); err != nil {
		return err
	}
	if err := a.provider().WriteConfig(cfg); err != nil {
		return err
	}
	if err := a.provider().TestConfig(ctx); err != nil {
		return err
	}
	if err := a.provider().Restart(ctx); err != nil {
		return err
	}
	return a.health(ctx, p)
}

func (a *App) ssh(ctx context.Context, p *output.Printer, args []string) error {
	if len(args) == 0 || args[0] != "harden" {
		system.AuditSSH(ctx, a.Runner, p)
		return nil
	}
	apply := contains(args, "--apply")
	system.PlanSSHHarden(ctx, a.Runner, p)
	if !apply {
		return nil
	}
	if err := a.prompt().Danger("ssh hardening", "SSH"); err != nil {
		return err
	}
	return system.ApplySSHHarden(ctx, a.Runner)
}

func (a *App) firewall(ctx context.Context, p *output.Printer, args []string) error {
	apply := contains(args, "--apply")
	system.PlanFirewall(ctx, a.Runner, p)
	if !apply {
		return nil
	}
	if err := a.prompt().Danger("enable firewall", "FIREWALL"); err != nil {
		return err
	}
	return system.ApplyFirewall(ctx, a.Runner, a.Options.ConfigPath)
}

func (a *App) bootCheck(ctx context.Context, p *output.Printer, args []string) error {
	err := a.health(ctx, p)
	if contains(args, "--record") {
		logDir := filepath.Join(a.Options.StateDir, "logs")
		_ = os.MkdirAll(logDir, 0700)
		_ = os.WriteFile(filepath.Join(logDir, "boot-check.log"), []byte("boot-check executed\n"), 0600)
	}
	return err
}

func contains(args []string, needle string) bool {
	for _, a := range args {
		if a == needle {
			return true
		}
	}
	return false
}
