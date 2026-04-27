package envfile

import (
	"fmt"
	"strings"
)

// RedactedEntry represents a key-value pair with optional redaction applied.
type RedactedEntry struct {
	Key      string
	Value    string
	Redacted bool
}

// RedactMode controls how secret values are redacted.
type RedactMode int

const (
	// RedactMask replaces the value with asterisks.
	RedactMask RedactMode = iota
	// RedactPartial shows the first and last character with asterisks in between.
	RedactPartial
	// RedactHash replaces the value with a fixed-length hash hint.
	RedactHash
)

// Redact applies redaction to all entries in the map based on the given mode.
// Non-secret keys are returned as-is.
func Redact(entries map[string]string, mode RedactMode) []RedactedEntry {
	result := make([]RedactedEntry, 0, len(entries))
	for k, v := range entries {
		if IsSecret(k) {
			result = append(result, RedactedEntry{
				Key:      k,
				Value:    applyRedaction(v, mode),
				Redacted: true,
			})
		} else {
			result = append(result, RedactedEntry{
				Key:      k,
				Value:    v,
				Redacted: false,
			})
		}
	}
	return result
}

// applyRedaction returns a redacted form of the value based on the mode.
func applyRedaction(value string, mode RedactMode) string {
	if value == "" {
		return ""
	}
	switch mode {
	case RedactPartial:
		if len(value) <= 2 {
			return strings.Repeat("*", len(value))
		}
		return string(value[0]) + strings.Repeat("*", len(value)-2) + string(value[len(value)-1])
	case RedactHash:
		return fmt.Sprintf("[redacted:%d]", len(value))
	default:
		return "********"
	}
}
