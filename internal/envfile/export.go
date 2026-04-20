package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// ExportFormat specifies the output format for env data.
type ExportFormat string

const (
	FormatDotEnv ExportFormat = "dotenv"
	FormatExport ExportFormat = "export"
	FormatJSON   ExportFormat = "json"
)

// Serialize converts an env map to a string in the given format.
// Keys are sorted alphabetically for deterministic output.
func Serialize(env map[string]string, format ExportFormat) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder

	switch format {
	case FormatExport:
		for _, k := range keys {
			fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
		}
	case FormatJSON:
		sb.WriteString("{\n")
		for i, k := range keys {
			comma := ","
			if i == len(keys)-1 {
				comma = ""
			}
			fmt.Fprintf(&sb, "  %q: %q%s\n", k, env[k], comma)
		}
		sb.WriteString("}\n")
	default: // dotenv
		for _, k := range keys {
			fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
		}
	}

	return sb.String()
}
