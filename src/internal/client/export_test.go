package client

import (
	"strings"
	"testing"

	"github.com/kehr/proxctl/src/internal/xray"
)

func TestExportGenericURIAndSurge(t *testing.T) {
	profile := xray.ClientProfile{
		Name:        "node",
		Server:      "example.com",
		Port:        443,
		UUID:        "11111111-2222-3333-4444-555555555555",
		Flow:        "xtls-rprx-vision",
		Network:     "tcp",
		Security:    "reality",
		SNI:         "www.microsoft.com",
		PublicKey:   "public-key",
		ShortID:     "abcdef",
		Fingerprint: "chrome",
		SpiderX:     "/",
	}

	uri, err := Export("generic-uri", profile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(uri, "vless://11111111-2222-3333-4444-555555555555@example.com:443") {
		t.Fatalf("bad uri: %s", uri)
	}
	if !strings.Contains(uri, "pbk=public-key") {
		t.Fatalf("missing public key: %s", uri)
	}

	surge, err := Export("surge", profile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(surge, "node = vless") || !strings.Contains(surge, "reality-public-key=public-key") {
		t.Fatalf("bad surge export: %s", surge)
	}
}
