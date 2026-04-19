package envfile

import "strings"

// secretPatterns are substrings that indicate a value should be masked.
var secretPatterns = []string{
	"secret",
	"password",
	"passwd",
	"token",
	"api_key",
	"apikey",
	"private",
	"auth",
	"credential",
}

const maskValue = "***"

// IsSecret returns true if the key looks like it holds a sensitive value.
func IsSecret(key string) bool {
	lower := strings.ToLower(key)
	for _, pattern := range secretPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}

// MaskEntry returns the entry value, replacing it with '***' if the key is
// considered secret.
func MaskEntry(e Entry) string {
	if IsSecret(e.Key) {
		return maskValue
	}
	return e.Value
}

// MaskedEntries returns a copy of all entries with secret values masked.
func MaskedEntries(env *EnvFile) []Entry {
	out := make([]Entry, len(env.Entries))
	for i, e := range env.Entries {
		out[i] = Entry{
			Key:   e.Key,
			Value: MaskEntry(e),
		}
	}
	return out
}
