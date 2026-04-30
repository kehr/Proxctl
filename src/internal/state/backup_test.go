package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBackupWritesManifestAndConfig(t *testing.T) {
	tmp := t.TempDir()
	config := filepath.Join(tmp, "config.json")
	if err := os.WriteFile(config, []byte(`{"log":{"loglevel":"warning"}}`), 0600); err != nil {
		t.Fatal(err)
	}
	manager := Manager{StateDir: filepath.Join(tmp, "state")}
	backup, err := manager.Backup("baseline", BackupInput{ConfigPath: config, XrayVersion: "Xray test"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(backup.Path, "config.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(backup.Path, "manifest.json")); err != nil {
		t.Fatal(err)
	}
}

func TestBackupSanitizesLabel(t *testing.T) {
	tmp := t.TempDir()
	config := filepath.Join(tmp, "config.json")
	if err := os.WriteFile(config, []byte(`{"log":{"loglevel":"warning"}}`), 0600); err != nil {
		t.Fatal(err)
	}
	manager := Manager{StateDir: filepath.Join(tmp, "state")}
	backup, err := manager.Backup("../bad label", BackupInput{ConfigPath: config})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(backup.Path) != filepath.Join(tmp, "state", "backups") {
		t.Fatalf("backup escaped backups dir: %s", backup.Path)
	}
	if filepath.Base(backup.Path) == "../bad label" {
		t.Fatalf("label was not sanitized: %s", backup.Path)
	}
}
