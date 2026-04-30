package xray

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
)

type Config map[string]any

type RotateTarget string

const (
	RotateUUID       RotateTarget = "uuid"
	RotateShortID    RotateTarget = "shortid"
	RotateRealityKey RotateTarget = "reality-key"
	RotateAll        RotateTarget = "all"
)

type RotationMeta struct {
	UUID              string `json:"uuid,omitempty"`
	ShortID           string `json:"short_id,omitempty"`
	RealityPrivateKey string `json:"-"`
	RealityPublicKey  string `json:"reality_public_key,omitempty"`
}

type Generator interface {
	UUID() (string, error)
	ShortID() (string, error)
	RealityKeyPair() (privateKey, publicKey string, err error)
}

type FixedGenerator struct{}

func (FixedGenerator) UUID() (string, error)    { return "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", nil }
func (FixedGenerator) ShortID() (string, error) { return "feedfacecafebeef", nil }
func (FixedGenerator) RealityKeyPair() (string, string, error) {
	return "fixed-private-key", "fixed-public-key", nil
}

func ParseConfig(b []byte) (Config, error) {
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return c, nil
}

func (c Config) Bytes() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

func (c Config) Summary() string {
	var b strings.Builder
	log, _ := c["log"].(map[string]any)
	fmt.Fprintf(&b, "loglevel=%v\n", log["loglevel"])
	for i, ib := range c.inbounds() {
		stream := object(ib["streamSettings"])
		settings := object(ib["settings"])
		clients := array(settings["clients"])
		fmt.Fprintf(&b, "inbound[%d]: protocol=%v listen=%v port=%v network=%v security=%v clients=%d\n",
			i+1, ib["protocol"], ib["listen"], ib["port"], stream["network"], stream["security"], len(clients))
		reality := object(stream["realitySettings"])
		if len(reality) > 0 {
			fmt.Fprintf(&b, "  reality: dest=%v serverNames=%s shortIds=%d privateKey=<redacted>\n",
				reality["dest"], strings.Join(stringArray(reality["serverNames"]), ","), len(array(reality["shortIds"])))
		}
	}
	for i, ob := range objects(c["outbounds"]) {
		fmt.Fprintf(&b, "outbound[%d]: tag=%v protocol=%v\n", i+1, ob["tag"], ob["protocol"])
	}
	routing := object(c["routing"])
	fmt.Fprintf(&b, "routing_rules=%d\n", len(array(routing["rules"])))
	return b.String()
}

func (c Config) ClientProfile(server, publicKey, name string) (ClientProfile, error) {
	ibs := c.inbounds()
	if len(ibs) == 0 {
		return ClientProfile{}, errors.New("no inbound found")
	}
	ib := ibs[0]
	stream := object(ib["streamSettings"])
	settings := object(ib["settings"])
	clients := objects(settings["clients"])
	if len(clients) == 0 {
		return ClientProfile{}, errors.New("no client found")
	}
	reality := object(stream["realitySettings"])
	return ClientProfile{
		Name:        name,
		Server:      server,
		Port:        intNumber(ib["port"], 443),
		UUID:        stringValue(clients[0]["id"]),
		Flow:        stringValue(clients[0]["flow"]),
		Network:     stringDefault(stream["network"], "tcp"),
		Security:    stringDefault(stream["security"], "reality"),
		SNI:         firstString(reality["serverNames"]),
		PublicKey:   publicKey,
		ShortID:     firstString(reality["shortIds"]),
		Fingerprint: stringDefault(reality["fingerprint"], "chrome"),
		SpiderX:     stringDefault(reality["spiderX"], "/"),
	}, nil
}

func (c Config) Rotate(target RotateTarget, gen Generator) (bool, RotationMeta, error) {
	switch target {
	case RotateUUID, RotateShortID, RotateRealityKey, RotateAll:
	default:
		return false, RotationMeta{}, fmt.Errorf("unknown rotate target: %s", target)
	}
	var meta RotationMeta
	changed := false
	for _, ib := range c.inbounds() {
		settings := object(ib["settings"])
		clients := objects(settings["clients"])
		if (target == RotateUUID || target == RotateAll) && len(clients) > 0 {
			u, err := gen.UUID()
			if err != nil {
				return false, meta, err
			}
			clients[0]["id"] = u
			meta.UUID = u
			changed = true
		}
		stream := object(ib["streamSettings"])
		reality := object(stream["realitySettings"])
		if len(reality) > 0 && (target == RotateShortID || target == RotateAll) {
			s, err := gen.ShortID()
			if err != nil {
				return false, meta, err
			}
			reality["shortIds"] = []any{s}
			meta.ShortID = s
			changed = true
		}
		if len(reality) > 0 && (target == RotateRealityKey || target == RotateAll) {
			priv, pub, err := gen.RealityKeyPair()
			if err != nil {
				return false, meta, err
			}
			reality["privateKey"] = priv
			meta.RealityPrivateKey = priv
			meta.RealityPublicKey = pub
			changed = true
		}
	}
	return changed, meta, nil
}

func (c Config) SetRealityTarget(dest string) error {
	host, _, err := net.SplitHostPort(dest)
	if err != nil {
		return errors.New("reality target must be host:port")
	}
	host = strings.Trim(host, "[]")
	for _, ib := range c.inbounds() {
		reality := object(object(ib["streamSettings"])["realitySettings"])
		if len(reality) > 0 {
			reality["dest"] = dest
			reality["serverNames"] = []any{host}
			return nil
		}
	}
	return errors.New("no reality inbound found")
}

type ClientProfile struct {
	Name        string
	Server      string
	Port        int
	UUID        string
	Flow        string
	Network     string
	Security    string
	SNI         string
	PublicKey   string
	ShortID     string
	Fingerprint string
	SpiderX     string
}

func (p ClientProfile) GenericURI() string {
	q := url.Values{}
	q.Set("type", p.Network)
	q.Set("security", p.Security)
	q.Set("encryption", "none")
	q.Set("fp", p.Fingerprint)
	q.Set("sni", p.SNI)
	q.Set("sid", p.ShortID)
	q.Set("spx", p.SpiderX)
	if p.PublicKey != "" {
		q.Set("pbk", p.PublicKey)
	}
	if p.Flow != "" {
		q.Set("flow", p.Flow)
	}
	return fmt.Sprintf("vless://%s@%s:%d?%s#%s", p.UUID, p.Server, p.Port, q.Encode(), url.QueryEscape(p.Name))
}

func (c Config) inbounds() []map[string]any { return objects(c["inbounds"]) }

func object(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}
func array(v any) []any {
	if a, ok := v.([]any); ok {
		return a
	}
	return nil
}
func objects(v any) []map[string]any {
	a := array(v)
	out := make([]map[string]any, 0, len(a))
	for _, item := range a {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}
func stringArray(v any) []string {
	a := array(v)
	out := make([]string, 0, len(a))
	for _, item := range a {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
func firstString(v any) string {
	a := stringArray(v)
	if len(a) > 0 {
		return a[0]
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
func stringValue(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
func stringDefault(v any, def string) string {
	if s := stringValue(v); s != "" {
		return s
	}
	return def
}
func intNumber(v any, def int) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	default:
		return def
	}
}
