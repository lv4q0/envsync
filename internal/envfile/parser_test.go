package envfile

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParse_BasicEntries(t *testing.T) {
	path := writeTempEnv(t, `
# comment
DB_HOST=localhost
DB_PORT=5432
SECRET_KEY="supersecret"
`)
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(env.Entries))
	}
	if env.Index["DB_HOST"].Value != "localhost" {
		t.Errorf("expected localhost, got %s", env.Index["DB_HOST"].Value)
	}
	if env.Index["SECRET_KEY"].Value != "supersecret" {
		t.Errorf("expected supersecret (unquoted), got %s", env.Index["SECRET_KEY"].Value)
	}
}

func TestParse_EmptyFile(t *testing.T) {
	path := writeTempEnv(t, "")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(env.Entries))
	}
}

func TestParse_MissingFile(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
