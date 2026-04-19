package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestDiffCmd_ShowsChanges(t *testing.T) {
	base := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\n")
	target := writeTempEnv(t, "APP_ENV=staging\nNEW_KEY=hello\n")

	rootCmd.SetArgs([]string{"diff", base, target})
	out := captureOutput(func() {
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected NEW_KEY in diff output, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in diff output, got: %s", out)
	}
}

func TestSyncCmd_DryRun(t *testing.T) {
	src := writeTempEnv(t, "FEATURE_FLAG=true\nAPI_KEY=secret\n")
	dst := writeTempEnv(t, "EXISTING=yes\n")

	rootCmd.SetArgs([]string{"sync", "--dry-run", src, dst})
	out := captureOutput(func() {
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "dry run") {
		t.Errorf("expected dry run notice, got: %s", out)
	}

	contents, _ := os.ReadFile(dst)
	if strings.Contains(string(contents), "FEATURE_FLAG") {
		t.Error("dry run should not write to destination")
	}
}
