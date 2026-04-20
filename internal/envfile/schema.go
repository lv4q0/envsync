package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

// SchemaRule defines a rule for a key in the schema.
type SchemaRule struct {
	Required bool   `json:"required"`
	Pattern  string `json:"pattern"`
	Secret   bool   `json:"secret"`
}

// Schema maps key names to their rules.
type Schema map[string]SchemaRule

// SchemaViolation describes a single schema violation.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("key %q: %s", v.Key, v.Message)
}

// LoadSchema reads a JSON schema file from disk.
func LoadSchema(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: read %q: %w", path, err)
	}
	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("schema: parse %q: %w", path, err)
	}
	return s, nil
}

// ValidateAgainstSchema checks env entries against the schema rules.
// It returns a slice of violations (empty means all good).
func ValidateAgainstSchema(entries map[string]string, schema Schema) []SchemaViolation {
	var violations []SchemaViolation

	for key, rule := range schema {
		val, present := entries[key]

		if rule.Required && !present {
			violations = append(violations, SchemaViolation{Key: key, Message: "required key is missing"})
			continue
		}

		if present && rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				violations = append(violations, SchemaViolation{Key: key, Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err)})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, SchemaViolation{Key: key, Message: fmt.Sprintf("value does not match pattern %q", rule.Pattern)})
			}
		}
	}

	return violations
}
