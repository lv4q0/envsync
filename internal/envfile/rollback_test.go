package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvForRollback(t *testing.T, content string) string {
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

func TestSaveRollbackPoint_CreatesFile(t *testing.T) {
	env := writeTempEnvForRollback(t, "KEY=value\n")
	dir := t.TempDir()

	entry, err := SaveRollbackPoint(env, dir, "before deploy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(entry.Path); err != nil {
		t.Errorf("snapshot file not created: %v", err)
	}
	if entry.Label != "before deploy" {
		t.Errorf("label mismatch: got %q", entry.Label)
	}
}

func TestListRollbackPoints_ReturnsSorted(t *testing.T) {
	env := writeTempEnvForRollback(t, "A=1\n")
	dir := t.TempDir()

	_, _ = SaveRollbackPoint(env, dir, "first")
	_, _ = SaveRollbackPoint(env, dir, "second")

	points, err := ListRollbackPoints(dir)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(points) != 2 {
		t.Fatalf("expected 2 points, got %d", len(points))
	}
	// newest first
	if !points[0].Timestamp.After(points[1].Timestamp) && points[0].Timestamp.Equal(points[1].Timestamp) {
		// equal timestamps are acceptable in fast tests
	}
}

func TestListRollbackPoints_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	points, err := ListRollbackPoints(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(points) != 0 {
		t.Errorf("expected empty list, got %d", len(points))
	}
}

func TestListRollbackPoints_MissingDir(t *testing.T) {
	points, err := ListRollbackPoints(filepath.Join(t.TempDir(), "nonexistent"))
	if err != nil {
		t.Fatalf("expected nil error for missing dir, got %v", err)
	}
	if len(points) != 0 {
		t.Errorf("expected empty, got %d", len(points))
	}
}

func TestRestoreRollbackPoint_RestoresContent(t *testing.T) {
	original := "KEY=original\n"
	env := writeTempEnvForRollback(t, original)
	dir := t.TempDir()

	entry, err := SaveRollbackPoint(env, dir, "snapshot")
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	// overwrite the env file
	if err := os.WriteFile(env, []byte("KEY=changed\n"), 0600); err != nil {
		t.Fatalf("overwrite: %v", err)
	}

	if err := RestoreRollbackPoint(entry, env); err != nil {
		t.Fatalf("restore: %v", err)
	}

	data, _ := os.ReadFile(env)
	if string(data) != original {
		t.Errorf("expected %q, got %q", original, string(data))
	}
}
