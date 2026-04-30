package xray

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/kehr/proxctl/src/internal/command"
)

type SystemGenerator struct {
	Runner  command.Runner
	XrayBin string
}

func (g SystemGenerator) UUID() (string, error) {
	if g.Runner != nil {
		bin := g.XrayBin
		if bin == "" {
			bin = "xray"
		}
		res := g.Runner.Run(context.Background(), bin, "uuid")
		if res.Code == 0 && strings.TrimSpace(res.Stdout) != "" {
			return strings.TrimSpace(res.Stdout), nil
		}
	}
	return randomUUID()
}

func (g SystemGenerator) ShortID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (g SystemGenerator) RealityKeyPair() (string, string, error) {
	if g.Runner == nil {
		return "", "", fmt.Errorf("xray x25519 is required to generate reality keypair")
	}
	bin := g.XrayBin
	if bin == "" {
		bin = "xray"
	}
	res := g.Runner.Run(context.Background(), bin, "x25519")
	if res.Code != 0 {
		return "", "", fmt.Errorf("xray x25519 failed: %s", res.Stderr)
	}
	var priv, pub string
	for _, line := range strings.Split(res.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "private key:") {
			priv = strings.TrimSpace(strings.TrimPrefix(line, "Private key:"))
		}
		if strings.HasPrefix(strings.ToLower(line), "public key:") {
			pub = strings.TrimSpace(strings.TrimPrefix(line, "Public key:"))
		}
	}
	if priv == "" || pub == "" {
		return "", "", fmt.Errorf("could not parse x25519 output")
	}
	return priv, pub, nil
}

func randomUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
