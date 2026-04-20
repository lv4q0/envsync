package envfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTempSchema(t *testing.T, schema Schema) string {
	t.Helper()
	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("marshal schema: %v", err)
	}
	dir := t.TempDir()
	p := filepath.Join(dir, ".env.schema.json")
	if err := os.WriteFile(p, data, 0644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	return p
}

func TestLoadSchema_Valid(t *testing.T) {
	path := writeTempSchema(t, Schema{
		"APP_ENV": {Required: true},
	})
	s, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s["APP_ENV"]; !ok {
		t.Error("expected APP_ENV in schema")
	}
}

func TestLoadSchema_MissingFile(t *testing.T) {
	_, err := LoadSchema("/nonexistent/schema.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidateAgainstSchema_RequiredMissing(t *testing.T) {
	schema := Schema{"DB_URL": {Required: true}}
	violations := ValidateAgainstSchema(map[string]string{}, schema)
	if len(violations) != 1 || violations[0].Key != "DB_URL" {
		t.Errorf("expected 1 violation for DB_URL, got %v", violations)
	}
}

func TestValidateAgainstSchema_PatternMatch(t *testing.T) {
	schema := Schema{"PORT": {Pattern: `^\d+$`}}
	violations := ValidateAgainstSchema(map[string]string{"PORT": "8080"}, schema)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestValidateAgainstSchema_PatternFail(t *testing.T) {
	schema := Schema{"PORT": {Pattern: `^\d+$`}}
	violations := ValidateAgainstSchema(map[string]string{"PORT": "abc"}, schema)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %v", violations)
	}
}

func TestValidateAgainstSchema_OptionalAbsent(t *testing.T) {
	schema := Schema{"OPTIONAL_KEY": {Required: false, Pattern: `^yes|no$`}}
	violations := ValidateAgainstSchema(map[string]string{}, schema)
	if len(violations) != 0 {
		t.Errorf("expected no violations for absent optional key, got %v", violations)
	}
}

func TestValidateAgainstSchema_NoViolations(t *testing.T) {
	schema := Schema{
		"APP_ENV": {Required: true, Pattern: `^(dev|staging|prod)$`},
		"SECRET_KEY": {Required: true, Secret: true},
	}
	entries := map[string]string{"APP_ENV": "prod", "SECRET_KEY": "abc123"}
	violations := ValidateAgainstSchema(entries, schema)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}
