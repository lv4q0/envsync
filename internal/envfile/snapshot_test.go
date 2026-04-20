package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvForSnapshot(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp env: %v", err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestTakeSnapshot_BasicEntries(t *testing.T) {
	path := writeTempEnvForSnapshot(t, "APP_ENV=production\nDB_HOST=localhost\n")
	snap, err := TakeSnapshot(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Entries["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", snap.Entries["APP_ENV"])
	}
	if snap.Source != path {
		t.Errorf("expected source %q, got %q", path, snap.Source)
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestTakeSnapshot_MissingFile(t *testing.T) {
	_, err := TakeSnapshot("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	path := writeTempEnvForSnapshot(t, "KEY1=val1\nKEY2=val2\n")
	snap, err := TakeSnapshot(path)
	if err != nil {
		t.Fatalf("TakeSnapshot: %v", err)
	}

	dest := filepath.Join(t.TempDir(), "snap.json")
	if err := SaveSnapshot(snap, dest); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	loaded, err := LoadSnapshot(dest)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	if loaded.Source != snap.Source {
		t.Errorf("source mismatch: got %q want %q", loaded.Source, snap.Source)
	}
	if loaded.Entries["KEY1"] != "val1" {
		t.Errorf("KEY1 mismatch: got %q", loaded.Entries["KEY1"])
	}
	if loaded.Entries["KEY2"] != "val2" {
		t.Errorf("KEY2 mismatch: got %q", loaded.Entries["KEY2"])
	}
}

func TestLoadSnapshot_InvalidFile(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "*.json")
	_, _ = f.WriteString("not json{{{")
	_ = f.Close()
	_, err := LoadSnapshot(f.Name())
	if err == nil {
		t.Error("expected decode error for invalid JSON")
	}
}
