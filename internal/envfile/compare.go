package envfile

import (
	"fmt"
	"sort"
)

// EnvDiff represents the result of comparing two env maps.
type EnvDiff struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // [old, new]
	Same    map[string]string
}

// CompareEnvMaps compares two parsed env maps and returns a structured diff.
func CompareEnvMaps(base, target map[string]string) EnvDiff {
	diff := EnvDiff{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
		Same:    make(map[string]string),
	}

	for k, tv := range target {
		if bv, ok := base[k]; !ok {
			diff.Added[k] = tv
		} else if bv != tv {
			diff.Changed[k] = [2]string{bv, tv}
		} else {
			diff.Same[k] = tv
		}
	}

	for k, bv := range base {
		if _, ok := target[k]; !ok {
			diff.Removed[k] = bv
		}
	}

	return diff
}

// FormatEnvDiff formats an EnvDiff into a human-readable string,
// masking secret values.
func FormatEnvDiff(d EnvDiff) string {
	out := ""

	keys := sortedMapKeys(d.Added)
	for _, k := range keys {
		v := d.Added[k]
		if IsSecret(k) {
			v = "***"
		}
		out += fmt.Sprintf("+ %s=%s\n", k, v)
	}

	keys = sortedMapKeys(d.Removed)
	for _, k := range keys {
		v := d.Removed[k]
		if IsSecret(k) {
			v = "***"
		}
		out += fmt.Sprintf("- %s=%s\n", k, v)
	}

	keys = sortedMapKeys(d.Changed)
	for _, k := range keys {
		old, nw := d.Changed[k][0], d.Changed[k][1]
		if IsSecret(k) {
			old, nw = "***", "***"
		}
		out += fmt.Sprintf("~ %s: %s -> %s\n", k, old, nw)
	}

	return out
}

func sortedMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
