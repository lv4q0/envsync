package envfile

import (
	"testing"
)

func TestValidate_ValidEntries(t *testing.T) {
	entries := map[string]string{
		"APP_NAME": "envsync",
		"PORT":     "8080",
		"API_KEY":  "supersecret",
	}
	result := Validate(entries)
	if result.HasErrors() {
		t.Errorf("expected no errors, got: %s", result.Error())
	}
}

func TestValidate_EmptyValueNonSecret(t *testing.T) {
	entries := map[string]string{
		"APP_NAME": "",
	}
	result := Validate(entries)
	if !result.HasErrors() {
		t.Error("expected error for empty non-secret value")
	}
	if len(result.Errors) != 1 || result.Errors[0].Key != "APP_NAME" {
		t.Errorf("unexpected errors: %v", result.Errors)
	}
}

func TestValidate_EmptyValueSecretAllowed(t *testing.T) {
	entries := map[string]string{
		"DB_PASSWORD": "",
	}
	result := Validate(entries)
	if result.HasErrors() {
		t.Errorf("expected no errors for empty secret value, got: %s", result.Error())
	}
}

func TestValidate_InvalidKeyChars(t *testing.T) {
	entries := map[string]string{
		"APP-NAME": "value",
	}
	result := Validate(entries)
	if !result.HasErrors() {
		t.Error("expected error for key with hyphen")
	}
}

func TestValidate_KeyWithSpaces(t *testing.T) {
	entries := map[string]string{
		"APP NAME": "value",
	}
	result := Validate(entries)
	if !result.HasErrors() {
		t.Error("expected error for key with space")
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	entries := map[string]string{
		"APP-NAME": "",
		"DB_HOST":  "",
	}
	result := Validate(entries)
	// APP-NAME: invalid chars + empty non-secret value; DB_HOST: empty non-secret value
	if len(result.Errors) < 2 {
		t.Errorf("expected at least 2 errors, got %d: %s", len(result.Errors), result.Error())
	}
}
