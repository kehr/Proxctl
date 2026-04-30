package xray

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kehr/proxctl/src/internal/command"
)

type Provider struct {
	ConfigPath string
	Service    string
	Bin        string
	Runner     command.Runner
}

func (p Provider) LoadConfig() (Config, error) {
	b, err := os.ReadFile(p.ConfigPath)
	if err != nil {
		return nil, err
	}
	return ParseConfig(b)
}

func (p Provider) WriteConfig(cfg Config) error {
	b, err := cfg.Bytes()
	if err != nil {
		return err
	}
	return p.WriteConfigBytes(p.ConfigPath, b)
}

func (p Provider) WriteConfigPath(cfg Config, path string) error {
	b, err := cfg.Bytes()
	if err != nil {
		return err
	}
	return p.WriteConfigBytes(path, b)
}

func (p Provider) WriteConfigBytes(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	tmp := path + ".proxctl.tmp"
	if err := os.WriteFile(tmp, b, 0600); err != nil {
		return err
	}
	defer os.Remove(tmp)
	if err := os.Chmod(tmp, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (p Provider) InstallConfigFromPath(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return p.WriteConfigBytes(p.ConfigPath, b)
}

func (p Provider) TestConfig(ctx context.Context) error {
	return p.TestConfigPath(ctx, p.ConfigPath)
}

func (p Provider) TestConfigPath(ctx context.Context, path string) error {
	bin := p.Bin
	if bin == "" {
		bin = "xray"
	}
	res := p.Runner.Run(ctx, bin, "run", "-test", "-config", path)
	if res.Code != 0 {
		return fmt.Errorf("%s%s", res.Stdout, res.Stderr)
	}
	return nil
}

func (p Provider) ServiceActive(ctx context.Context) bool {
	service := p.Service
	if service == "" {
		service = "xray"
	}
	return p.Runner.Run(ctx, "systemctl", "is-active", "--quiet", service).Code == 0
}

func (p Provider) Restart(ctx context.Context) error {
	service := p.Service
	if service == "" {
		service = "xray"
	}
	res := p.Runner.Run(ctx, "systemctl", "restart", service)
	if res.Code != 0 {
		return fmt.Errorf("restart %s failed: %s", service, res.Stderr)
	}
	return nil
}

func (p Provider) Version(ctx context.Context) string {
	bin := p.Bin
	if bin == "" {
		bin = "xray"
	}
	return command.FirstLine(p.Runner.Run(ctx, bin, "version").Stdout)
}

func (p Provider) ServiceUnit(ctx context.Context) string {
	service := p.Service
	if service == "" {
		service = "xray"
	}
	return p.Runner.Run(ctx, "systemctl", "cat", service).Stdout
}

func (p Provider) Ports() ([]int, error) {
	cfg, err := p.LoadConfig()
	if err != nil {
		return nil, err
	}
	var ports []int
	for _, ib := range cfg.inbounds() {
		ports = append(ports, intNumber(ib["port"], 0))
	}
	return ports, nil
}

func (p Provider) PortOwnedByService(ctx context.Context, port int) bool {
	if port == 0 {
		return false
	}
	res := p.Runner.Run(ctx, "ss", "-H", "-ltnp")
	if res.Code != 0 || res.Stdout == "" {
		return false
	}
	service := p.Service
	if service == "" {
		service = "xray"
	}
	portSuffix := ":" + strconv.Itoa(port)
	for _, line := range strings.Split(res.Stdout, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		localAddr := fields[3]
		if strings.HasSuffix(localAddr, portSuffix) && strings.Contains(line, service) {
			return true
		}
	}
	return false
}

func (p Provider) InitDefault(ctx context.Context) error {
	gen := SystemGenerator{Runner: p.Runner, XrayBin: p.Bin}
	uuid, err := gen.UUID()
	if err != nil {
		return err
	}
	shortID, err := gen.ShortID()
	if err != nil {
		return err
	}
	priv, _, err := gen.RealityKeyPair()
	if err != nil {
		return err
	}
	cfg := Config{
		"log": map[string]any{"loglevel": "warning"},
		"inbounds": []any{map[string]any{
			"listen":   "0.0.0.0",
			"port":     float64(443),
			"protocol": "vless",
			"settings": map[string]any{
				"clients":    []any{map[string]any{"id": uuid, "flow": "xtls-rprx-vision"}},
				"decryption": "none",
			},
			"streamSettings": map[string]any{
				"network":  "tcp",
				"security": "reality",
				"realitySettings": map[string]any{
					"dest":        "www.microsoft.com:443",
					"serverNames": []any{"www.microsoft.com"},
					"privateKey":  priv,
					"shortIds":    []any{shortID},
				},
			},
		}},
		"outbounds": []any{
			map[string]any{"protocol": "freedom", "tag": "direct"},
			map[string]any{"protocol": "blackhole", "tag": "block"},
		},
	}
	b, _ := json.Marshal(cfg)
	parsed, _ := ParseConfig(b)
	return p.WriteConfig(parsed)
}
