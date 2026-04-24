package envfile

import (
	"errors"
	"fmt"
)

// RotateResult holds the outcome of a key rotation operation.
type RotateResult struct {
	Rotated []string
	Skipped []string
	Errors  []string
}

// Rotate re-encrypts all secret values in env using oldKey and newKey.
// Non-secret keys are left unchanged. Returns a RotateResult summarising
// which keys were rotated, skipped, or failed.
func Rotate(env map[string]string, oldKey, newKey []byte) (map[string]string, RotateResult, error) {
	if len(oldKey) == 0 {
		return nil, RotateResult{}, errors.New("rotate: oldKey must not be empty")
	}
	if len(newKey) == 0 {
		return nil, RotateResult{}, errors.New("rotate: newKey must not be empty")
	}

	result := make(map[string]string, len(env))
	var rr RotateResult

	for k, v := range env {
		if !IsSecret(k) {
			result[k] = v
			rr.Skipped = append(rr.Skipped, k)
			continue
		}

		// Decrypt with old key.
		plaintext, err := Decrypt(v, oldKey)
		if err != nil {
			rr.Errors = append(rr.Errors, fmt.Sprintf("%s: decrypt failed: %v", k, err))
			result[k] = v // preserve original on error
			continue
		}

		// Re-encrypt with new key.
		ciphertext, err := Encrypt(plaintext, newKey)
		if err != nil {
			rr.Errors = append(rr.Errors, fmt.Sprintf("%s: encrypt failed: %v", k, err))
			result[k] = v
			continue
		}

		result[k] = ciphertext
		rr.Rotated = append(rr.Rotated, k)
	}

	if len(rr.Errors) > 0 {
		return result, rr, fmt.Errorf("rotate: %d key(s) failed rotation", len(rr.Errors))
	}
	return result, rr, nil
}
