package envfile

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// ErrInvalidKey is returned when the encryption key is not 16, 24, or 32 bytes.
var ErrInvalidKey = errors.New("encryption key must be 16, 24, or 32 bytes")

// ErrInvalidCiphertext is returned when decryption fails due to bad input.
var ErrInvalidCiphertext = errors.New("invalid ciphertext")

// Encrypt encrypts a plaintext string using AES-GCM with the provided key.
// The key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
// Returns a base64-encoded ciphertext string.
func Encrypt(plaintext string, key []byte) (string, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", ErrInvalidKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-encoded ciphertext string using AES-GCM.
// The key must match the one used during encryption.
func Decrypt(encoded string, key []byte) (string, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", ErrInvalidKey
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrInvalidCiphertext
	}

	return string(plaintext), nil
}

// EncryptSecrets returns a copy of the env map with all secret values encrypted.
func EncryptSecrets(env map[string]string, key []byte) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		if IsSecret(k) {
			encrypted, err := Encrypt(v, key)
			if err != nil {
				return nil, err
			}
			result[k] = encrypted
		} else {
			result[k] = v
		}
	}
	return result, nil
}

// DecryptSecrets returns a copy of the env map with all secret values decrypted.
func DecryptSecrets(env map[string]string, key []byte) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		if IsSecret(k) {
			decrypted, err := Decrypt(v, key)
			if err != nil {
				return nil, err
			}
			result[k] = decrypted
		} else {
			result[k] = v
		}
	}
	return result, nil
}
