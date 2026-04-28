package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvForClone(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp env: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestClone_BasicCopy(t *testing.T) {
	src := writeTempEnvForClone(t, "APP_NAME=envsync\nPORT=8080\n")
	dest := filepath.Join(t.TempDir(), ".env")

	res, err := Clone(src, dest, CloneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.KeysWritten != 2 {
		t.Errorf("expected 2 keys written, got %d", res.KeysWritten)
	}
	if res.SecretsStripped != 0 {
		t.Errorf("expected 0 secrets stripped, got %d", res.SecretsStripped)
	}
	if _, err := os.Stat(dest); err != nil {
		t.Errorf("destination file not created: %v", err)
	}
}

func TestClone_StripSecrets(t *testing.T) {
	src := writeTempEnvForClone(t, "APP_NAME=envsync\nAPI_SECRET=supersecret\nDB_PASSWORD=hunter2\n")
	dest := filepath.Join(t.TempDir(), ".env")

	res, err := Clone(src, dest, CloneOptions{StripSecrets: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.SecretsStripped != 2 {
		t.Errorf("expected 2 secrets stripped, got %d", res.SecretsStripped)
	}

	data, _ := os.ReadFile(dest)
	if strings.Contains(string(data), "supersecret") || strings.Contains(string(data), "hunter2") {
		t.Errorf("secret values should have been stripped from clone")
	}
}

func TestClone_NoOverwrite_ReturnsError(t *testing.T) {
	src := writeTempEnvForClone(t, "KEY=val\n")
	dest := writeTempEnvForClone(t, "EXISTING=true\n")

	_, err := Clone(src, dest, CloneOptions{Overwrite: false})
	if err == nil {
		t.Fatal("expected error when destination exists and overwrite is false")
	}
}

func TestClone_Overwrite_ReplacesFile(t *testing.T) {
	src := writeTempEnvForClone(t, "KEY=newval\n")
	dest := writeTempEnvForClone(t, "KEY=oldval\n")

	_, err := Clone(src, dest, CloneOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(dest)
	if !strings.Contains(string(data), "newval") {
		t.Errorf("expected destination to contain new value after overwrite")
	}
}

func TestClone_WithProfile_CreatesProfileFile(t *testing.T) {
	src := writeTempEnvForClone(t, "APP=test\n")
	dir := t.TempDir()
	dest := filepath.Join(dir, ".env")

	res, err := Clone(src, dest, CloneOptions{Profile: "staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(dir, ".env.staging")
	if res.Destination != expected {
		t.Errorf("expected destination %q, got %q", expected, res.Destination)
	}
	if _, err := os.Stat(expected); err != nil {
		t.Errorf("profile file not created: %v", err)
	}
}

func TestClone_MissingSource_ReturnsError(t *testing.T) {
	_, err := Clone("/nonexistent/.env", "/tmp/dest.env", CloneOptions{})
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}
