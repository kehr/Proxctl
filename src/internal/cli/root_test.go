package cli

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kehr/proxctl/src/internal/command"
)

func TestRootHelpIncludesCompletionAndDocs(t *testing.T) {
	var out, err bytes.Buffer
	rt := NewRuntime(&out, &err, strings.NewReader(""))
	if e := Execute(context.Background(), []string{"--help"}, rt); e != nil {
		t.Fatal(e)
	}
	text := out.String()
	for _, want := range []string{"completion", "docs", "client", "apply"} {
		if !strings.Contains(text, want) {
			t.Fatalf("help missing %q:\n%s", want, text)
		}
	}
}

func TestClientExportDefaultsServerAndPublicKey(t *testing.T) {
	ipServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("203.0.113.10\n"))
	}))
	defer ipServer.Close()
	oldPublicIPURLs := publicIPURLs
	publicIPURLs = []string{ipServer.URL}
	defer func() { publicIPURLs = oldPublicIPURLs }()

	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	if err := os.WriteFile(configPath, []byte(strings.ReplaceAll(sampleXrayConfig, "PRIVATE_SHOULD_NOT_LEAK", "server-private-key")), 0600); err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	rt := NewRuntime(&out, &errOut, strings.NewReader(""))
	rt.Runner = fakeRunner{results: map[string]command.Result{
		"xray\x00x25519\x00-i\x00server-private-key": {
			Stdout: "Private key: server-private-key\nPublic key: derived-public-key\n",
		},
	}}

	err := Execute(context.Background(), []string{
		"--xray-config", configPath,
		"client", "export", "generic-uri",
		"--name", "test-node",
	}, rt)
	if err != nil {
		t.Fatal(err)
	}
	text := out.String()
	if !strings.Contains(text, "@203.0.113.10:443") {
		t.Fatalf("export missing detected public IP:\n%s", text)
	}
	if !strings.Contains(text, "pbk=derived-public-key") {
		t.Fatalf("export missing derived public key:\n%s", text)
	}
}

const sampleXrayConfig = `{
  "log": {"loglevel": "warning"},
  "inbounds": [{
    "listen": "0.0.0.0",
    "port": 443,
    "protocol": "vless",
    "settings": {
      "clients": [{"id": "11111111-2222-3333-4444-555555555555", "flow": "xtls-rprx-vision"}],
      "decryption": "none"
    },
    "streamSettings": {
      "network": "tcp",
      "security": "reality",
      "realitySettings": {
        "dest": "www.microsoft.com:443",
        "serverNames": ["www.microsoft.com"],
        "privateKey": "PRIVATE_SHOULD_NOT_LEAK",
        "shortIds": ["abcdef1234567890"]
      }
    }
  }],
  "outbounds": [{"protocol": "freedom", "tag": "direct"}]
}`

type fakeRunner struct {
	results map[string]command.Result
}

func (r fakeRunner) Run(ctx context.Context, name string, args ...string) command.Result {
	key := strings.Join(append([]string{name}, args...), "\x00")
	if res, ok := r.results[key]; ok {
		return res
	}
	return command.Result{Code: 127, Stderr: "unexpected command: " + key}
}

func (r fakeRunner) LookPath(name string) bool {
	return true
}
