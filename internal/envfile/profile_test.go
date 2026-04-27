package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeProfileEnv(t *testing.T, dir, profileName, content string) {
	t.Helper()
	path := filepath.Join(dir, ".env."+profileName)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeProfileEnv: %v", err)
	}
}

func TestLoadProfile_BasicEntries(t *testing.T) {
	dir := t.TempDir()
	writeProfileEnv(t, dir, "staging", "APP_ENV=staging\nDB_HOST=staging-db\n")

	p, err := LoadProfile(dir, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "staging" {
		t.Errorf("expected name 'staging', got %q", p.Name)
	}
	if p.Entries["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", p.Entries["APP_ENV"])
	}
}

func TestLoadProfile_EmptyName_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadProfile(dir, "")
	if err == nil {
		t.Fatal("expected error for empty profile name")
	}
}

func TestLoadProfile_MissingFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadProfile(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing profile file")
	}
}

func TestListProfiles_FindsProfiles(t *testing.T) {
	dir := t.TempDir()
	writeProfileEnv(t, dir, "staging", "KEY=val\n")
	writeProfileEnv(t, dir, "production", "KEY=val\n")

	profiles, err := ListProfiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(profiles))
	}
}

func TestListProfiles_EmptyDir_ReturnsNone(t *testing.T) {
	dir := t.TempDir()
	profiles, err := ListProfiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("expected 0 profiles, got %d", len(profiles))
	}
}

func TestDiffProfiles_DetectsDifferences(t *testing.T) {
	base := &Profile{Name: "staging", Entries: map[string]string{
		"APP_ENV": "staging",
		"DB_HOST": "staging-db",
		"REMOVED": "gone",
	}}
	target := &Profile{Name: "production", Entries: map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "staging-db",
		"ADDED":   "new",
	}}

	diffs := DiffProfiles(base, target)

	if _, ok := diffs["APP_ENV"]; !ok {
		t.Error("expected APP_ENV to be in diffs")
	}
	if _, ok := diffs["DB_HOST"]; ok {
		t.Error("expected DB_HOST to NOT be in diffs (same value)")
	}
	if _, ok := diffs["REMOVED"]; !ok {
		t.Error("expected REMOVED to be in diffs")
	}
	if _, ok := diffs["ADDED"]; !ok {
		t.Error("expected ADDED to be in diffs")
	}
}
