package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeCmpEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeCmpEnv: %v", err)
	}
	return p
}

func TestCompareCmd_NoChanges(t *testing.T) {
	dir := t.TempDir()
	base := writeCmpEnv(t, dir, "base.env", "APP=hello\nPORT=8080\n")
	target := writeCmpEnv(t, dir, "target.env", "APP=hello\nPORT=8080\n")

	out, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"compare", base, target})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected no-diff message, got: %s", out)
	}
}

func TestCompareCmd_ShowsAdded(t *testing.T) {
	dir := t.TempDir()
	base := writeCmpEnv(t, dir, "base.env", "APP=hello\n")
	target := writeCmpEnv(t, dir, "target.env", "APP=hello\nNEW_KEY=world\n")

	out, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"compare", base, target})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "+ NEW_KEY=world") {
		t.Errorf("expected added key in output, got: %s", out)
	}
}

func TestCompareCmd_MasksSecrets(t *testing.T) {
	dir := t.TempDir()
	base := writeCmpEnv(t, dir, "base.env", "SECRET_TOKEN=old\n")
	target := writeCmpEnv(t, dir, "target.env", "SECRET_TOKEN=new\n")

	out, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"compare", base, target})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "old") || strings.Contains(out, "new") {
		t.Errorf("secret values should be masked, got: %s", out)
	}
	if !strings.Contains(out, "SECRET_TOKEN") {
		t.Errorf("key name should still appear in output")
	}
}
