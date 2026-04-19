package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestSync_AddsNewKeys(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=envsync\nNEW_KEY=hello\n")
	dst := writeTempEnv(t, "APP_NAME=envsync\n")

	res, err := Sync(src, dst, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 1 || res.Applied[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY applied, got %v", res.Applied)
	}

	parsed, _ := envfile.Parse(dst)
	if parsed["NEW_KEY"] != "hello" {
		t.Errorf("NEW_KEY not written to destination")
	}
}

func TestSync_DryRunDoesNotWrite(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=envsync\nNEW_KEY=hello\n")
	dst := writeTempEnv(t, "APP_NAME=envsync\n")

	_, err := Sync(src, dst, Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsed, _ := envfile.Parse(dst)
	if _, ok := parsed["NEW_KEY"]; ok {
		t.Error("dry run should not write NEW_KEY to destination")
	}
}

func TestSync_OverwriteChanged(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=newname\n")
	dst := writeTempEnv(t, "APP_NAME=oldname\n")

	res, err := Sync(src, dst, Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied change, got %d", len(res.Applied))
	}

	parsed, _ := envfile.Parse(dst)
	if parsed["APP_NAME"] != "newname" {
		t.Errorf("APP_NAME should be overwritten, got %q", parsed["APP_NAME"])
	}
}

func TestSync_SkipsChangedWithoutOverwrite(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=newname\n")
	dst := writeTempEnv(t, "APP_NAME=oldname\n")

	res, err := Sync(src, dst, Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped key, got %d", len(res.Skipped))
	}

	parsed, _ := envfile.Parse(dst)
	if parsed["APP_NAME"] != "oldname" {
		t.Errorf("APP_NAME should remain oldname, got %q", parsed["APP_NAME"])
	}
}

func TestSync_MissingDestination(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=envsync\n")
	dst := filepath.Join(t.TempDir(), "missing.env")

	res, err := Sync(src, dst, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied key for new destination, got %d", len(res.Applied))
	}
}
