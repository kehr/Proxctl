package state

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Manager struct {
	StateDir string
}

type BackupInput struct {
	ConfigPath      string
	XrayVersion     string
	ServiceUnitText string
	SSHDEffective   string
	CommandLine     []string
}

type Backup struct {
	ID   string
	Path string
}

type Manifest struct {
	ID                  string    `json:"id"`
	Label               string    `json:"label"`
	CreatedAt           time.Time `json:"created_at"`
	Host                string    `json:"host"`
	ConfigPath          string    `json:"config_path"`
	ConfigSHA256        string    `json:"config_sha256"`
	XrayVersion         string    `json:"xray_version,omitempty"`
	RestoreCommand      string    `json:"restore_command"`
	OriginalCommandLine []string  `json:"original_command_line,omitempty"`
}

func (m Manager) Ensure() error {
	for _, dir := range []string{m.StateDir, filepath.Join(m.StateDir, "backups"), filepath.Join(m.StateDir, "pending")} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}

func (m Manager) Backup(label string, in BackupInput) (Backup, error) {
	if label == "" {
		label = "manual"
	}
	label = sanitizeLabel(label)
	if err := m.Ensure(); err != nil {
		return Backup{}, err
	}
	id := time.Now().UTC().Format("20060102T150405Z") + "-" + label
	dir := filepath.Join(m.StateDir, "backups", id)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return Backup{}, err
	}
	hash := ""
	if in.ConfigPath != "" {
		if err := copyFile(in.ConfigPath, filepath.Join(dir, "config.json"), 0600); err != nil {
			return Backup{}, err
		}
		h, err := fileSHA256(in.ConfigPath)
		if err != nil {
			return Backup{}, err
		}
		hash = h
	}
	if in.ServiceUnitText != "" {
		_ = os.WriteFile(filepath.Join(dir, "xray.service.txt"), []byte(in.ServiceUnitText), 0600)
	}
	if in.SSHDEffective != "" {
		_ = os.WriteFile(filepath.Join(dir, "sshd-effective.txt"), []byte(in.SSHDEffective), 0600)
	}
	host, _ := os.Hostname()
	manifest := Manifest{
		ID:                  id,
		Label:               label,
		CreatedAt:           time.Now().UTC(),
		Host:                host,
		ConfigPath:          in.ConfigPath,
		ConfigSHA256:        hash,
		XrayVersion:         in.XrayVersion,
		RestoreCommand:      "proxctl restore " + id,
		OriginalCommandLine: in.CommandLine,
	}
	b, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return Backup{}, err
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), b, 0600); err != nil {
		return Backup{}, err
	}
	latest := filepath.Join(m.StateDir, "backups", "latest")
	_ = os.Remove(latest)
	_ = os.Symlink(dir, latest)
	return Backup{ID: id, Path: dir}, nil
}

func (m Manager) RestoreConfig(id, targetPath string) (string, error) {
	src, err := m.ConfigBackupPath(id)
	if err != nil {
		return "", err
	}
	if err := copyFile(src, targetPath, 0600); err != nil {
		return "", err
	}
	return src, nil
}

func (m Manager) ConfigBackupPath(id string) (string, error) {
	if id == "" || id == "latest" {
		id = "latest"
	}
	src := filepath.Join(m.StateDir, "backups", id, "config.json")
	if id == "latest" {
		src = filepath.Join(m.StateDir, "backups", "latest", "config.json")
	}
	if _, err := os.Stat(src); err != nil {
		return "", err
	}
	return src, nil
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
		return err
	}
	tmp := dst + ".tmp"
	defer os.Remove(tmp)
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp, mode); err != nil {
		return err
	}
	if err := os.Rename(tmp, dst); err != nil {
		return err
	}
	return nil
}

func sanitizeLabel(label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		return "manual"
	}
	var b strings.Builder
	for _, r := range label {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '.', r == '_', r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}
	out := strings.Trim(b.String(), ".-_")
	if out == "" {
		return "manual"
	}
	return out
}

func (m Manager) LockPath() string {
	return filepath.Join(m.StateDir, "proxctl.lock")
}

func (m Manager) SaveState(v any) error {
	if err := m.Ensure(); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(m.StateDir, "state.json"), b, 0600)
}

func (m Manager) LoadState(v any) error {
	b, err := os.ReadFile(filepath.Join(m.StateDir, "state.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("load state: %w", err)
	}
	return nil
}
