package envfile

import (
	"testing"
)

var testKey16 = []byte("0123456789abcdef") // 16 bytes
var testKey32 = []byte("0123456789abcdef0123456789abcdef") // 32 bytes

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plaintext := "super-secret-value"
	encrypted, err := Encrypt(plaintext, testKey16)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	if encrypted == plaintext {
		t.Error("encrypted text should differ from plaintext")
	}
	decrypted, err := Decrypt(encrypted, testKey16)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncrypt_InvalidKeyLength(t *testing.T) {
	_, err := Encrypt("value", []byte("shortkey"))
	if err != ErrInvalidKey {
		t.Errorf("expected ErrInvalidKey, got %v", err)
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := Decrypt("not-valid-base64!!!", testKey16)
	if err != ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	encrypted, err := Encrypt("my-secret", testKey16)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	_, err = Decrypt(encrypted, testKey32)
	if err != ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext with wrong key, got %v", err)
	}
}

func TestEncryptSecrets_OnlyEncryptsSecretKeys(t *testing.T) {
	env := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "hunter2",
		"API_SECRET":  "topsecret",
	}
	result, err := EncryptSecrets(env, testKey16)
	if err != nil {
		t.Fatalf("EncryptSecrets failed: %v", err)
	}
	if result["APP_NAME"] != "myapp" {
		t.Errorf("non-secret key should be unchanged, got %q", result["APP_NAME"])
	}
	if result["DB_PASSWORD"] == "hunter2" {
		t.Error("DB_PASSWORD should be encrypted")
	}
	if result["API_SECRET"] == "topsecret" {
		t.Error("API_SECRET should be encrypted")
	}
}

func TestDecryptSecrets_RoundTrip(t *testing.T) {
	env := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "hunter2",
	}
	encrypted, err := EncryptSecrets(env, testKey16)
	if err != nil {
		t.Fatalf("EncryptSecrets failed: %v", err)
	}
	decrypted, err := DecryptSecrets(encrypted, testKey16)
	if err != nil {
		t.Fatalf("DecryptSecrets failed: %v", err)
	}
	for k, v := range env {
		if decrypted[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, decrypted[k])
		}
	}
}
