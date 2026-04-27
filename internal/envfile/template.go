package envfile

import (
	"fmt"
	"os"
	"strings"
)

// TemplateEntry represents a key with an optional default value and required flag.
type TemplateEntry struct {
	Key      string
	Default  string
	Required bool
	Comment  string
}

// Template holds a parsed .env.template file.
type Template struct {
	Entries []TemplateEntry
}

// LoadTemplate parses a .env.template file where entries may look like:
//   KEY=                   # required, no default
//   KEY=default_value      # optional with default
//   # comment line
func LoadTemplate(path string) (*Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("template: read %q: %w", path, err)
	}

	tmpl := &Template{}
	lines := strings.Split(string(data), "\n")
	var lastComment string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			lastComment = ""
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			lastComment = strings.TrimSpace(strings.TrimPrefix(trimmed, "#"))
			continue
		}
		idx := strings.IndexByte(trimmed, '=')
		if idx < 0 {
			lastComment = ""
			continue
		}
		key := strings.TrimSpace(trimmed[:idx])
		defaultVal := strings.TrimSpace(trimmed[idx+1:])
		if key == "" {
			lastComment = ""
			continue
		}
		tmpl.Entries = append(tmpl.Entries, TemplateEntry{
			Key:      key,
			Default:  defaultVal,
			Required: defaultVal == "",
			Comment:  lastComment,
		})
		lastComment = ""
	}
	return tmpl, nil
}

// ApplyTemplate fills missing keys in env from the template defaults.
// Returns a list of keys that are required but still missing after applying defaults.
func ApplyTemplate(env map[string]string, tmpl *Template) (missing []string) {
	for _, entry := range tmpl.Entries {
		if _, ok := env[entry.Key]; !ok {
			if entry.Default != "" {
				env[entry.Key] = entry.Default
			} else {
				missing = append(missing, entry.Key)
			}
		}
	}
	return missing
}
