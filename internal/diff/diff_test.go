package diff_test

import (
	"testing"

	"github.com/user/envsync/internal/diff"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "myapp",
		"DB_PASSWORD": "secret123",
		"PORT": "8080",
	}
}

func TestCompare_Added(t *testing.T) {
	base := baseEnv()
	target := baseEnv()
	target["NEW_KEY"] = "newval"

	entries := diff.Compare(base, target)
	found := findEntry(entries, "NEW_KEY")
	if found == nil || found.Status != diff.StatusAdded {
		t.Errorf("expected NEW_KEY to be added")
	}
}

func TestCompare_Removed(t *testing.T) {
	base := baseEnv()
	target := baseEnv()
	delete(target, "PORT")

	entries := diff.Compare(base, target)
	found := findEntry(entries, "PORT")
	if found == nil || found.Status != diff.StatusRemoved {
		t.Errorf("expected PORT to be removed")
	}
}

func TestCompare_Changed(t *testing.T) {
	base := baseEnv()
	target := baseEnv()
	target["APP_NAME"] = "otherapp"

	entries := diff.Compare(base, target)
	found := findEntry(entries, "APP_NAME")
	if found == nil || found.Status != diff.StatusChanged {
		t.Errorf("expected APP_NAME to be changed")
	}
}

func TestCompare_SecretMaskedInString(t *testing.T) {
	base := baseEnv()
	target := baseEnv()
	target["DB_PASSWORD"] = "newsecret"

	entries := diff.Compare(base, target)
	found := findEntry(entries, "DB_PASSWORD")
	if found == nil {
		t.Fatal("DB_PASSWORD entry not found")
	}
	if !found.Secret {
		t.Errorf("expected DB_PASSWORD to be marked as secret")
	}
	s := found.String()
	if contains(s, "secret123") || contains(s, "newsecret") {
		t.Errorf("secret value should be masked in String(), got: %s", s)
	}
}

func findEntry(entries []diff.Entry, key string) *diff.Entry {
	for i := range entries {
		if entries[i].Key == key {
			return &entries[i]
		}
	}
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
