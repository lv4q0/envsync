package envfile

import (
	"fmt"
	"sort"
)

// SnapshotDiff represents the difference between two snapshots.
type SnapshotDiff struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string]SnapshotChange
}

// SnapshotChange holds the before/after values for a changed key.
type SnapshotChange struct {
	Old string
	New string
}

// DiffSnapshots compares two snapshots and returns the differences.
func DiffSnapshots(base, current map[string]string) SnapshotDiff {
	diff := SnapshotDiff{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]SnapshotChange),
	}

	for k, v := range current {
		if old, ok := base[k]; !ok {
			diff.Added[k] = v
		} else if old != v {
			diff.Changed[k] = SnapshotChange{Old: old, New: v}
		}
	}

	for k, v := range base {
		if _, ok := current[k]; !ok {
			diff.Removed[k] = v
		}
	}

	return diff
}

// FormatSnapshotDiff returns a human-readable string of the snapshot diff.
func FormatSnapshotDiff(diff SnapshotDiff, maskSecrets bool) string {
	if len(diff.Added) == 0 && len(diff.Removed) == 0 && len(diff.Changed) == 0 {
		return "No changes detected between snapshots.\n"
	}

	result := ""

	keys := sortedKeys(diff.Added)
	for _, k := range keys {
		v := diff.Added[k]
		if maskSecrets && IsSecret(k) {
			v = "***"
		}
		result += fmt.Sprintf("+ %s=%s\n", k, v)
	}

	keys = sortedKeys(diff.Removed)
	for _, k := range keys {
		v := diff.Removed[k]
		if maskSecrets && IsSecret(k) {
			v = "***"
		}
		result += fmt.Sprintf("- %s=%s\n", k, v)
	}

	keys = sortedKeys(diff.Changed)
	for _, k := range keys {
		old := diff.Changed[k].Old
		new := diff.Changed[k].New
		if maskSecrets && IsSecret(k) {
			old = "***"
			new = "***"
		}
		result += fmt.Sprintf("~ %s: %s -> %s\n", k, old, new)
	}

	return result
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
