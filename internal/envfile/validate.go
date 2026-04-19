package envfile

import (
	"fmt"
	"strings"
)

// ValidationError represents a single validation issue in an env file.
type ValidationError struct {
	Line    int
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d: key %q: %s", e.Line, e.Key, e.Message)
	}
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ValidationResult) Error() string {
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "\n")
}

// Validate checks a parsed env map for common issues such as empty keys,
// keys with invalid characters, and empty values for required (non-secret) keys.
func Validate(entries map[string]string) *ValidationResult {
	result := &ValidationResult{}

	for key, value := range entries {
		if key == "" {
			result.Errors = append(result.Errors, ValidationError{
				Key:     key,
				Message: "empty key is not allowed",
			})
			continue
		}

		if strings.ContainsAny(key, " \t") {
			result.Errors = append(result.Errors, ValidationError{
				Key:     key,
				Message: "key must not contain spaces or tabs",
			})
		}

		if !isValidKeyChars(key) {
			result.Errors = append(result.Errors, ValidationError{
				Key:     key,
				Message: "key contains invalid characters (only A-Z, a-z, 0-9, _ allowed)",
			})
		}

		if value == "" && !IsSecret(key) {
			result.Errors = append(result.Errors, ValidationError{
				Key:     key,
				Message: "non-secret key has empty value",
			})
		}
	}

	return result
}

func isValidKeyChars(key string) bool {
	for _, c := range key {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}
