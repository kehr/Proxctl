package xray

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kehr/proxctl/src/internal/command"
)

type fakeRunner struct {
	stdout string
	stderr string
	code   int
}

func (r fakeRunner) Run(ctx context.Context, name string, args ...string) command.Result {
	return command.Result{Stdout: r.stdout, Stderr: r.stderr, Code: r.code}
}

func (r fakeRunner) LookPath(name string) bool { return true }

func TestWriteConfigPathDoesNotTouchActiveConfig(t *testing.T) {
	tmp := t.TempDir()
	active := filepath.Join(tmp, "config.json")
	if err := os.WriteFile(active, []byte(`{"log":{"loglevel":"warning"}}`), 0600); err != nil {
		t.Fatal(err)
	}
	provider := Provider{ConfigPath: active}
	cfg, err := ParseConfig([]byte(sampleConfig))
	if err != nil {
		t.Fatal(err)
	}
	pending := filepath.Join(tmp, "pending", "config.json")
	if err := provider.WriteConfigPath(cfg, pending); err != nil {
		t.Fatal(err)
	}
	activeBytes, err := os.ReadFile(active)
	if err != nil {
		t.Fatal(err)
	}
	if string(activeBytes) != `{"log":{"loglevel":"warning"}}` {
		t.Fatalf("active config changed before install: %s", activeBytes)
	}
	if err := provider.InstallConfigFromPath(pending); err != nil {
		t.Fatal(err)
	}
	activeBytes, err = os.ReadFile(active)
	if err != nil {
		t.Fatal(err)
	}
	if string(activeBytes) == `{"log":{"loglevel":"warning"}}` {
		t.Fatal("active config was not installed")
	}
}

func TestPortOwnedByServiceParsesSSOutput(t *testing.T) {
	out := `LISTEN 0 4096 0.0.0.0:22 0.0.0.0:* users:(("sshd",pid=1,fd=3))
LISTEN 0 4096 [::]:443 [::]:* users:(("xray",pid=2,fd=3))`
	provider := Provider{Service: "xray", Runner: fakeRunner{stdout: out}}
	if !provider.PortOwnedByService(context.Background(), 443) {
		t.Fatal("expected xray to own 443")
	}
	if provider.PortOwnedByService(context.Background(), 22) {
		t.Fatal("did not expect xray to own 22")
	}
}

func TestRealityPublicKeyDerivesFromPrivateKey(t *testing.T) {
	runner := fakeRunner{stdout: "PrivateKey: server-private-key\nPassword (PublicKey): derived-public-key\nHash32: ignored\n"}
	pub, err := RealityPublicKey(context.Background(), runner, "xray", "server-private-key")
	if err != nil {
		t.Fatal(err)
	}
	if pub != "derived-public-key" {
		t.Fatalf("unexpected public key: %s", pub)
	}
}
