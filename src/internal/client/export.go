package client

import (
	"fmt"
	"strings"

	"github.com/kehr/proxctl/src/internal/xray"
)

func Formats() []string {
	return []string{"generic-uri", "shadowrocket", "surge", "stash", "mihomo", "sing-box", "v2rayn", "v2rayng"}
}

func Export(format string, p xray.ClientProfile) (string, error) {
	switch strings.ToLower(format) {
	case "generic-uri", "shadowrocket", "v2rayn", "v2rayng":
		return p.GenericURI(), nil
	case "surge":
		parts := []string{
			fmt.Sprintf("%s = vless", p.Name),
			fmt.Sprintf("server=%s", p.Server),
			fmt.Sprintf("port=%d", p.Port),
			fmt.Sprintf("username=%s", p.UUID),
			"tls=true",
			"reality=true",
			fmt.Sprintf("sni=%s", p.SNI),
			fmt.Sprintf("reality-public-key=%s", p.PublicKey),
			fmt.Sprintf("reality-short-id=%s", p.ShortID),
			fmt.Sprintf("client-fingerprint=%s", p.Fingerprint),
		}
		if p.Flow != "" {
			parts = append(parts, "flow="+p.Flow)
		}
		return strings.Join(parts, ", "), nil
	case "stash", "mihomo":
		return fmt.Sprintf("- name: %s\n  type: vless\n  server: %s\n  port: %d\n  uuid: %s\n  network: %s\n  tls: true\n  reality-opts:\n    public-key: %s\n    short-id: %s\n  servername: %s\n  client-fingerprint: %s\n", p.Name, p.Server, p.Port, p.UUID, p.Network, p.PublicKey, p.ShortID, p.SNI, p.Fingerprint), nil
	case "sing-box":
		return fmt.Sprintf(`{
  "type": "vless",
  "tag": "%s",
  "server": "%s",
  "server_port": %d,
  "uuid": "%s",
  "flow": "%s",
  "tls": {
    "enabled": true,
    "server_name": "%s",
    "utls": {"enabled": true, "fingerprint": "%s"},
    "reality": {"enabled": true, "public_key": "%s", "short_id": "%s"}
  }
}`, p.Name, p.Server, p.Port, p.UUID, p.Flow, p.SNI, p.Fingerprint, p.PublicKey, p.ShortID), nil
	default:
		return "", fmt.Errorf("unsupported client format: %s", format)
	}
}
