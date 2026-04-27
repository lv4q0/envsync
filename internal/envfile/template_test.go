package envfile

import (
	"os"
	"testing"
)

func writeTempTemplate(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env.template")
	if err != nil {
		t.Fatalf("create temp template: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp template: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadTemplate_BasicEntries(t *testing.T) {
	path := writeTempTemplate(t, "# database url\nDB_URL=\nAPP_ENV=development\n")
	tmpl, err := LoadTemplate(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tmpl.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(tmpl.Entries))
	}
	if tmpl.Entries[0].Key != "DB_URL" || !tmpl.Entries[0].Required {
		t.Errorf("DB_URL should be required with no default")
	}
	if tmpl.Entries[0].Comment != "database url" {
		t.Errorf("expected comment 'database url', got %q", tmpl.Entries[0].Comment)
	}
	if tmpl.Entries[1].Key != "APP_ENV" || tmpl.Entries[1].Default != "development" {
		t.Errorf("APP_ENV should have default 'development'")
	}
	if tmpl.Entries[1].Required {
		t.Errorf("APP_ENV should not be required")
	}
}

func TestLoadTemplate_MissingFile(t *testing.T) {
	_, err := LoadTemplate("/nonexistent/path.env.template")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestApplyTemplate_FillsDefaults(t *testing.T) {
	path := writeTempTemplate(t, "APP_ENV=production\nLOG_LEVEL=info\n")
	tmpl, _ := LoadTemplate(path)
	env := map[string]string{"APP_ENV": "staging"}
	missing := ApplyTemplate(env, tmpl)
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got %v", missing)
	}
	if env["APP_ENV"] != "staging" {
		t.Errorf("existing key should not be overwritten")
	}
	if env["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info from template default")
	}
}

func TestApplyTemplate_ReportsMissingRequired(t *testing.T) {
	path := writeTempTemplate(t, "SECRET_KEY=\nDB_URL=\nAPP_ENV=development\n")
	tmpl, _ := LoadTemplate(path)
	env := map[string]string{}
	missing := ApplyTemplate(env, tmpl)
	if len(missing) != 2 {
		t.Errorf("expected 2 missing required keys, got %v", missing)
	}
	if env["APP_ENV"] != "development" {
		t.Errorf("expected APP_ENV filled from default")
	}
}

func TestApplyTemplate_EmptyTemplate(t *testing.T) {
	path := writeTempTemplate(t, "# just a comment\n")
	tmpl, err := LoadTemplate(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := map[string]string{"KEY": "val"}
	missing := ApplyTemplate(env, tmpl)
	if len(missing) != 0 {
		t.Errorf("expected no missing keys for empty template")
	}
}
