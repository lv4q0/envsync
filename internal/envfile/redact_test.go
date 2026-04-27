package envfile

import (
	"strings"
	"testing"
)

func findEntry(entries []RedactedEntry, key string) (RedactedEntry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return RedactedEntry{}, false
}

func TestRedact_MaskMode_SecretIsHidden(t *testing.T) {
	entries := map[string]string{
		"API_SECRET": "supersecret",
		"APP_NAME":   "myapp",
	}
	result := Redact(entries, RedactMask)

	secret, ok := findEntry(result, "API_SECRET")
	if !ok {
		t.Fatal("expected API_SECRET in result")
	}
	if secret.Value != "********" {
		t.Errorf("expected masked value, got %q", secret.Value)
	}
	if !secret.Redacted {
		t.Error("expected Redacted=true for secret key")
	}

	plain, ok := findEntry(result, "APP_NAME")
	if !ok {
		t.Fatal("expected APP_NAME in result")
	}
	if plain.Value != "myapp" {
		t.Errorf("expected plain value, got %q", plain.Value)
	}
	if plain.Redacted {
		t.Error("expected Redacted=false for non-secret key")
	}
}

func TestRedact_PartialMode_ShowsFirstAndLast(t *testing.T) {
	entries := map[string]string{
		"DB_PASSWORD": "abcdef",
	}
	result := Redact(entries, RedactPartial)

	e, ok := findEntry(result, "DB_PASSWORD")
	if !ok {
		t.Fatal("expected DB_PASSWORD in result")
	}
	if !strings.HasPrefix(e.Value, "a") || !strings.HasSuffix(e.Value, "f") {
		t.Errorf("expected partial redaction like a****f, got %q", e.Value)
	}
}

func TestRedact_HashMode_ShowsLength(t *testing.T) {
	entries := map[string]string{
		"AUTH_TOKEN": "tok123",
	}
	result := Redact(entries, RedactHash)

	e, ok := findEntry(result, "AUTH_TOKEN")
	if !ok {
		t.Fatal("expected AUTH_TOKEN in result")
	}
	if e.Value != "[redacted:6]" {
		t.Errorf("expected [redacted:6], got %q", e.Value)
	}
}

func TestRedact_EmptySecretValue_ReturnsEmpty(t *testing.T) {
	entries := map[string]string{
		"API_KEY": "",
	}
	result := Redact(entries, RedactMask)

	e, ok := findEntry(result, "API_KEY")
	if !ok {
		t.Fatal("expected API_KEY in result")
	}
	if e.Value != "" {
		t.Errorf("expected empty string for empty secret, got %q", e.Value)
	}
}

func TestRedact_PartialMode_ShortValue(t *testing.T) {
	entries := map[string]string{
		"DB_PASSWORD": "ab",
	}
	result := Redact(entries, RedactPartial)

	e, ok := findEntry(result, "DB_PASSWORD")
	if !ok {
		t.Fatal("expected DB_PASSWORD in result")
	}
	if e.Value != "**" {
		t.Errorf("expected '**' for short value, got %q", e.Value)
	}
}
