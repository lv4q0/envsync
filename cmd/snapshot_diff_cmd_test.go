package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSnapshotFile(t *testing.T, data map[string]string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func writeEnvFileForSnapshotDiff(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestSnapshotDiffCmd_NoChanges(t *testing.T) {
	snap := writeSnapshotFile(t, map[string]string{"FOO": "bar"})
	env := writeEnvFileForSnapshotDiff(t, "FOO=bar\n")

	rootCmd.SetArgs([]string{"snapshot-diff", snap, env, "--mask=false"})
	out := captureOutput(t, func() {
		_ = rootCmd.Execute()
	})

	if !strings.Contains(out, "No changes") {
		t.Errorf("expected no-changes message, got: %s", out)
	}
}

func TestSnapshotDiffCmd_ShowsAdded(t *testing.T) {
	snap := writeSnapshotFile(t, map[string]string{"FOO": "bar"})
	env := writeEnvFileForSnapshotDiff(t, "FOO=bar\nNEW_VAR=hello\n")

	rootCmd.SetArgs([]string{"snapshot-diff", snap, env, "--mask=false"})
	out := captureOutput(t, func() {
		_ = rootCmd.Execute()
	})

	if !strings.Contains(out, "+ NEW_VAR") {
		t.Errorf("expected added key in output, got: %s", out)
	}
}

func TestSnapshotDiffCmd_MasksSecrets(t *testing.T) {
	snap := writeSnapshotFile(t, map[string]string{"API_SECRET": "oldval"})
	env := writeEnvFileForSnapshotDiff(t, "API_SECRET=newval\n")

	rootCmd.SetArgs([]string{"snapshot-diff", snap, env, "--mask=true"})
	out := captureOutput(t, func() {
		_ = rootCmd.Execute()
	})

	if strings.Contains(out, "oldval") || strings.Contains(out, "newval") {
		t.Errorf("expected secrets masked, got: %s", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected masked placeholder, got: %s", out)
	}
}
