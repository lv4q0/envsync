package envfile

import (
	"fmt"
	"strings"
)

// LintSeverity represents the severity level of a lint warning.
type LintSeverity string

const (
	LintWarn  LintSeverity = "WARN"
	LintError LintSeverity = "ERROR"
	LintInfo  LintSeverity = "INFO"
)

// LintIssue describes a single linting problem found in an env file.
type LintIssue struct {
	Line     int
	Key      string
	Message  string
	Severity LintSeverity
}

func (i LintIssue) String() string {
	if i.Line > 0 {
		return fmt.Sprintf("[%s] line %d (%s): %s", i.Severity, i.Line, i.Key, i.Message)
	}
	return fmt.Sprintf("[%s] (%s): %s", i.Severity, i.Key, i.Message)
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Issues []LintIssue
}

// HasErrors returns true if any ERROR-level issues exist.
func (r *LintResult) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == LintError {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary line.
func (r *LintResult) Summary() string {
	if len(r.Issues) == 0 {
		return "No lint issues found."
	}
	errCount := 0
	warnCount := 0
	for _, issue := range r.Issues {
		switch issue.Severity {
		case LintError:
			errCount++
		case LintWarn:
			warnCount++
		}
	}
	return fmt.Sprintf("%d error(s), %d warning(s) found.", errCount, warnCount)
}

// Lint analyses a parsed env map and raw lines for common style and correctness issues.
// It checks for:
//   - keys that are not uppercase
//   - values containing unquoted whitespace
//   - duplicate keys (requires raw lines)
//   - secret keys with empty values
//   - very long values (>500 chars) which may indicate accidental pastes
func Lint(entries map[string]string, rawLines []string) *LintResult {
	result := &LintResult{}

	seen := make(map[string]int) // key -> first line number

	for lineNum, raw := range rawLines {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		idx := strings.IndexByte(trimmed, '=')
		if idx < 0 {
			continue
		}

		key := strings.TrimSpace(trimmed[:idx])
		value := trimmed[idx+1:]

		// Check for duplicate keys
		if first, exists := seen[key]; exists {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum + 1,
				Key:      key,
				Message:  fmt.Sprintf("duplicate key (first defined on line %d)", first),
				Severity: LintError,
			})
		} else {
			seen[key] = lineNum + 1
		}

		// Check for non-uppercase keys
		if key != strings.ToUpper(key) {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum + 1,
				Key:      key,
				Message:  "key should be uppercase",
				Severity: LintWarn,
			})
		}

		// Check for unquoted leading/trailing whitespace in value
		if !strings.HasPrefix(value, `"`) && !strings.HasPrefix(value, "'") {
			if value != strings.TrimSpace(value) {
				result.Issues = append(result.Issues, LintIssue{
					Line:     lineNum + 1,
					Key:      key,
					Message:  "value has unquoted leading or trailing whitespace",
					Severity: LintWarn,
				})
			}
		}
	}

	// Check entries from parsed map for semantic issues
	for key, value := range entries {
		// Secret keys must not have empty values
		if IsSecret(key) && strings.TrimSpace(value) == "" {
			result.Issues = append(result.Issues, LintIssue{
				Key:      key,
				Message:  "secret key has an empty value",
				Severity: LintError,
			})
		}

		// Warn on suspiciously long values
		if len(value) > 500 {
			result.Issues = append(result.Issues, LintIssue{
				Key:      key,
				Message:  fmt.Sprintf("value is unusually long (%d chars); verify this is intentional", len(value)),
				Severity: LintWarn,
			})
		}
	}

	return result
}
