package system

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kehr/proxctl/src/internal/command"
	"github.com/kehr/proxctl/src/internal/output"
)

type Snapshot struct {
	Hostname string
	OS       string
	Kernel   string
	Uptime   string
	Memory   string
	Disk     string
}

func Collect(ctx context.Context, r command.Runner) Snapshot {
	host, _ := os.Hostname()
	return Snapshot{
		Hostname: host,
		OS:       osRelease(),
		Kernel:   command.FirstLine(r.Run(ctx, "uname", "-r").Stdout),
		Uptime:   strings.TrimSpace(r.Run(ctx, "uptime").Stdout),
		Memory:   r.Run(ctx, "free", "-h").Stdout,
		Disk:     r.Run(ctx, "df", "-hT", "/", "/boot").Stdout,
	}
}

func osRelease() string {
	b, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "unknown"
	}
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`)
		}
	}
	return "unknown"
}

func SSHDEffective(ctx context.Context, r command.Runner) string {
	return r.Run(ctx, "sshd", "-T").Stdout
}

func AuditSSH(ctx context.Context, r command.Runner, p *output.Printer) {
	eff := SSHDEffective(ctx, r)
	if eff == "" {
		p.Warn("Cannot read sshd effective config.")
		return
	}
	keys := []string{"port", "passwordauthentication", "permitrootlogin", "pubkeyauthentication", "kbdinteractiveauthentication", "x11forwarding", "maxauthtries"}
	for _, key := range keys {
		for _, line := range strings.Split(eff, "\n") {
			if strings.HasPrefix(line, key+" ") {
				p.Info(line)
			}
		}
	}
	if strings.Contains(eff, "passwordauthentication no\n") {
		p.Pass("PasswordAuthentication is disabled.")
	} else {
		p.Warn("PasswordAuthentication should be disabled after confirming key login.")
	}
	if strings.Contains(eff, "permitrootlogin prohibit-password\n") || strings.Contains(eff, "permitrootlogin without-password\n") || strings.Contains(eff, "permitrootlogin no\n") {
		p.Pass("PermitRootLogin is restricted.")
	} else {
		p.Warn("PermitRootLogin is permissive.")
	}
	if strings.Contains(eff, "x11forwarding no\n") {
		p.Pass("X11Forwarding is disabled.")
	} else {
		p.Warn("X11Forwarding should be disabled.")
	}
}

func PlanSSHHarden(ctx context.Context, r command.Runner, p *output.Printer) {
	AuditSSH(ctx, r, p)
	p.Info("Plan: write /etc/ssh/sshd_config.d/90-proxctl-hardening.conf")
	p.Info("Plan: run sshd -t, reload sshd, require a new SSH session confirmation")
}

func ApplySSHHarden(ctx context.Context, r command.Runner) error {
	content := "PasswordAuthentication no\nPermitRootLogin prohibit-password\nKbdInteractiveAuthentication no\nX11Forwarding no\nMaxAuthTries 3\n"
	if err := os.WriteFile("/etc/ssh/sshd_config.d/90-proxctl-hardening.conf", []byte(content), 0600); err != nil {
		return err
	}
	if res := r.Run(ctx, "sshd", "-t"); res.Code != 0 {
		return fmt.Errorf("sshd -t failed: %s", res.Stderr)
	}
	if res := r.Run(ctx, "systemctl", "reload", "sshd"); res.Code != 0 {
		return fmt.Errorf("reload sshd failed: %s", res.Stderr)
	}
	return nil
}

func AuditFirewall(ctx context.Context, r command.Runner, p *output.Printer) {
	if r.LookPath("firewall-cmd") {
		p.Info("firewalld: %s", strings.TrimSpace(r.Run(ctx, "systemctl", "is-active", "firewalld").Stdout+r.Run(ctx, "systemctl", "is-active", "firewalld").Stderr))
		out := r.Run(ctx, "firewall-cmd", "--list-all").Stdout
		if out != "" {
			fmt.Fprint(p.Out, out)
		}
	} else {
		p.Warn("firewall-cmd not found.")
	}
	if r.LookPath("nft") {
		out := r.Run(ctx, "nft", "list", "ruleset").Stdout
		if out != "" {
			fmt.Fprint(p.Out, out)
		}
	}
}

func PlanFirewall(ctx context.Context, r command.Runner, p *output.Printer) {
	AuditFirewall(ctx, r, p)
	p.Info("Plan: allow SSH and proxy ports before enabling firewalld.")
	p.Info("Plan: use rollback timer before committing firewall changes.")
}

func ApplyFirewall(ctx context.Context, r command.Runner, configPath string) error {
	if res := r.Run(ctx, "firewall-cmd", "--permanent", "--add-service=ssh"); res.Code != 0 {
		return fmt.Errorf("allow ssh failed: %s", res.Stderr)
	}
	if res := r.Run(ctx, "firewall-cmd", "--permanent", "--add-port=443/tcp"); res.Code != 0 {
		return fmt.Errorf("allow 443 failed: %s", res.Stderr)
	}
	if res := r.Run(ctx, "systemctl", "enable", "--now", "firewalld"); res.Code != 0 {
		return fmt.Errorf("enable firewalld failed: %s", res.Stderr)
	}
	if res := r.Run(ctx, "firewall-cmd", "--reload"); res.Code != 0 {
		return fmt.Errorf("reload firewalld failed: %s", res.Stderr)
	}
	return nil
}

func CheckUpdates(ctx context.Context, r command.Runner, p *output.Printer) {
	if !r.LookPath("dnf") {
		p.Warn("dnf not found.")
		return
	}
	res := r.Run(ctx, "dnf", "check-update")
	if res.Code == 0 {
		p.Pass("No package updates available.")
		return
	}
	if res.Code == 100 {
		p.Warn("Package updates are available.")
		fmt.Fprint(p.Out, firstLines(res.Stdout, 30))
		return
	}
	p.Warn("dnf check-update returned %d", res.Code)
}

func firstLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, "\n")
}
