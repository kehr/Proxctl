package xray

import (
	"strings"
	"testing"
)

const sampleConfig = `{
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

func TestSummaryRedactsRealityPrivateKey(t *testing.T) {
	cfg, err := ParseConfig([]byte(sampleConfig))
	if err != nil {
		t.Fatal(err)
	}
	summary := cfg.Summary()
	if !strings.Contains(summary, "protocol=vless") {
		t.Fatalf("summary missing protocol: %s", summary)
	}
	if !strings.Contains(summary, "security=reality") {
		t.Fatalf("summary missing security: %s", summary)
	}
	if strings.Contains(summary, "PRIVATE_SHOULD_NOT_LEAK") {
		t.Fatalf("summary leaked private key: %s", summary)
	}
}

func TestRotateAllChangesCredentialsAndKeepsPrivateKeyOutOfClient(t *testing.T) {
	cfg, err := ParseConfig([]byte(sampleConfig))
	if err != nil {
		t.Fatal(err)
	}
	changed, meta, err := cfg.Rotate(RotateAll, FixedGenerator{})
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("expected config change")
	}
	client, err := cfg.ClientProfile("example.com", meta.RealityPublicKey, "test-node")
	if err != nil {
		t.Fatal(err)
	}
	if client.UUID != "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee" {
		t.Fatalf("uuid not rotated: %s", client.UUID)
	}
	if client.PublicKey != "fixed-public-key" {
		t.Fatalf("public key not set: %s", client.PublicKey)
	}
	if strings.Contains(client.GenericURI(), "fixed-private-key") {
		t.Fatalf("client URI leaked private key: %s", client.GenericURI())
	}
}
