package envfile

import (
	"strings"
	"testing"
)

func makeKey(b byte) []byte {
	key := make([]byte, 32)
	for i := range key {
		key[i] = b
	}
	return key
}

func TestRotate_RotatesSecretKeys(t *testing.T) {
	oldKey := makeKey(0xAA)
	newKey := makeKey(0xBB)

	plain := "supersecret"
	cipher, err := Encrypt(plain, oldKey)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	env := map[string]string{
		"DB_PASSWORD": cipher,
		"APP_NAME":    "myapp",
	}

	result, rr, err := Rotate(env, oldKey, newKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rr.Rotated) != 1 || rr.Rotated[0] != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD rotated, got %v", rr.Rotated)
	}
	if len(rr.Skipped) != 1 || rr.Skipped[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME skipped, got %v", rr.Skipped)
	}

	// Verify new ciphertext decrypts correctly with newKey.
	decrypted, err := Decrypt(result["DB_PASSWORD"], newKey)
	if err != nil {
		t.Fatalf("decrypt with new key: %v", err)
	}
	if decrypted != plain {
		t.Errorf("expected %q, got %q", plain, decrypted)
	}

	// Non-secret key must be unchanged.
	if result["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should be unchanged")
	}
}

func TestRotate_WrongOldKey_ReturnsError(t *testing.T) {
	oldKey := makeKey(0xAA)
	wrongKey := makeKey(0xCC)
	newKey := makeKey(0xBB)

	cipher, _ := Encrypt("value", oldKey)
	env := map[string]string{"API_SECRET": cipher}

	_, rr, err := Rotate(env, wrongKey, newKey)
	if err == nil {
		t.Fatal("expected error for wrong old key")
	}
	if len(rr.Errors) == 0 {
		t.Error("expected errors in RotateResult")
	}
	if !strings.Contains(rr.Errors[0], "API_SECRET") {
		t.Errorf("error should mention key name, got: %s", rr.Errors[0])
	}
}

func TestRotate_EmptyOldKey_ReturnsError(t *testing.T) {
	_, _, err := Rotate(map[string]string{}, []byte{}, makeKey(0x01))
	if err == nil {
		t.Fatal("expected error for empty oldKey")
	}
}

func TestRotate_EmptyNewKey_ReturnsError(t *testing.T) {
	_, _, err := Rotate(map[string]string{}, makeKey(0x01), []byte{})
	if err == nil {
		t.Fatal("expected error for empty newKey")
	}
}

func TestRotate_NoSecrets_AllSkipped(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	_, rr, err := Rotate(env, makeKey(0x01), makeKey(0x02))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rr.Rotated) != 0 {
		t.Errorf("expected no rotated keys, got %v", rr.Rotated)
	}
	if len(rr.Skipped) != 2 {
		t.Errorf("expected 2 skipped keys, got %d", len(rr.Skipped))
	}
}
