package cli

import (
	"io"
	"os"

	"github.com/kehr/proxctl/src/internal/command"
	"github.com/kehr/proxctl/src/internal/output"
	"github.com/kehr/proxctl/src/internal/prompt"
	"github.com/kehr/proxctl/src/internal/state"
	"github.com/kehr/proxctl/src/internal/xray"
)

const Version = "0.2.0"

type Runtime struct {
	Out    io.Writer
	Err    io.Writer
	In     io.Reader
	Runner command.Runner
	Config Config
}

type Config struct {
	Provider   string
	ConfigPath string
	StateDir   string
	XrayBin    string
	Service    string
	Yes        bool
	JSON       bool
	NoColor    bool
}

func DefaultConfig() Config {
	return Config{
		Provider:   "xray",
		ConfigPath: "/usr/local/etc/xray/config.json",
		StateDir:   "/var/lib/proxctl",
		XrayBin:    "xray",
		Service:    "xray",
	}
}

func NewRuntime(out, err io.Writer, in io.Reader) *Runtime {
	return &Runtime{
		Out:    out,
		Err:    err,
		In:     in,
		Runner: command.LocalRunner{},
		Config: DefaultConfig(),
	}
}

func (r *Runtime) Printer() *output.Printer {
	return &output.Printer{Out: r.Out, Err: r.Err, JSON: r.Config.JSON, Color: !r.Config.NoColor}
}

func (r *Runtime) Prompt() prompt.Prompter {
	return prompt.Prompter{In: r.In, Out: r.Out, AssumeYes: r.Config.Yes}
}

func (r *Runtime) State() state.Manager {
	return state.Manager{StateDir: r.Config.StateDir}
}

func (r *Runtime) Xray() xray.Provider {
	return xray.Provider{ConfigPath: r.Config.ConfigPath, Service: r.Config.Service, Bin: r.Config.XrayBin, Runner: r.Runner}
}

func (r *Runtime) ApplyEnvFallbacks() {
	if v := os.Getenv("PROXCTL_CONFIG"); v != "" {
		r.Config.ConfigPath = v
	}
	if v := os.Getenv("PROXCTL_STATE_DIR"); v != "" {
		r.Config.StateDir = v
	}
	if v := os.Getenv("PROXCTL_XRAY_BIN"); v != "" {
		r.Config.XrayBin = v
	}
}
