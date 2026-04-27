package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"envsync/internal/envfile"
)

func writeRollbackEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	return p
}

func runRollbackCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"rollback"}, args...))
	err := rootCmd.Execute()
	rootCmd.SetArgs(nil)
	return buf.String(), err
}

func TestRollbackCmd_SavesPoint(t *testing.T) {
	dir := t.TempDir()
	env := writeRollbackEnv(t, dir, ".env", "KEY=val\n")
	rollDir := filepath.Join(dir, "rolls")

	out, err := runRollbackCmd(t,
		"--file", env,
		"--dir", rollDir,
		"--label", "test save",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Rollback point saved") {
		t.Errorf("expected save confirmation, got: %q", out)
	}
}

func TestRollbackCmd_ListPoints(t *testing.T) {
	dir := t.TempDir()
	env := writeRollbackEnv(t, dir, ".env", "KEY=val\n")
	rollDir := filepath.Join(dir, "rolls")

	_, _ = envfile.SaveRollbackPoint(env, rollDir, "alpha")
	_, _ = envfile.SaveRollbackPoint(env, rollDir, "beta")

	out, err := runRollbackCmd(t,
		"--file", env,
		"--dir", rollDir,
		"--list",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "beta") {
		t.Errorf("expected labels in list output, got: %q", out)
	}
}

func TestRollbackCmd_RestoresContent(t *testing.T) {
	dir := t.TempDir()
	env := writeRollbackEnv(t, dir, ".env", "KEY=original\n")
	rollDir := filepath.Join(dir, "rolls")

	entry, err := envfile.SaveRollbackPoint(env, rollDir, "pre")
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	// mutate the file
	_ = os.WriteFile(env, []byte("KEY=changed\n"), 0600)

	out, err := runRollbackCmd(t,
		"--file", env,
		"--dir", rollDir,
		"--restore", entry.Path,
	)
	if err != nil {
		t.Fatalf("restore error: %v", err)
	}
	if !strings.Contains(out, "Restored") {
		t.Errorf("expected restore confirmation, got: %q", out)
	}
	data, _ := os.ReadFile(env)
	if string(data) != "KEY=original\n" {
		t.Errorf("expected original content, got: %q", string(data))
	}
}

var _ = (*cobra.Command)(nil) // ensure cobra import used
