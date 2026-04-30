package system

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/kehr/proxctl/src/internal/command"
)

func TestProxyPortsReadsUniqueInboundPorts(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.json")
	if err := os.WriteFile(path, []byte(`{"inbounds":[{"port":443},{"port":8443},{"port":443}]}`), 0600); err != nil {
		t.Fatal(err)
	}
	ports, err := proxyPorts(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []int{443, 8443}
	if !reflect.DeepEqual(ports, want) {
		t.Fatalf("ports=%v want=%v", ports, want)
	}
}

type recordingRunner struct {
	active bool
	calls  []string
}

func (r *recordingRunner) Run(ctx context.Context, name string, args ...string) command.Result {
	call := name + " " + strings.Join(args, " ")
	r.calls = append(r.calls, call)
	if call == "systemctl is-active --quiet firewalld" && !r.active {
		return command.Result{Code: 3}
	}
	return command.Result{}
}

func (r *recordingRunner) LookPath(name string) bool { return true }

func TestApplyFirewallUsesOfflineCommandWhenFirewalldInactive(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.json")
	if err := os.WriteFile(path, []byte(`{"inbounds":[{"port":443}]}`), 0600); err != nil {
		t.Fatal(err)
	}
	runner := &recordingRunner{}
	if err := ApplyFirewall(context.Background(), runner, path); err != nil {
		t.Fatal(err)
	}
	got := strings.Join(runner.calls, "\n")
	if !strings.Contains(got, "firewall-offline-cmd --add-service=ssh") {
		t.Fatalf("offline service rule missing:\n%s", got)
	}
	if strings.Contains(got, "firewall-offline-cmd --permanent") {
		t.Fatalf("offline command must not use --permanent:\n%s", got)
	}
	if strings.Contains(got, "firewall-cmd --reload") {
		t.Fatalf("inactive firewalld should not reload before startup:\n%s", got)
	}
}
