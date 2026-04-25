package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// VersionDiff represents a single key difference between two versioned env maps.
type VersionDiff struct {
	Key    string
	From   string
	To     string
	Change string // "added", "removed", "changed"
}

// CompareVersions compares two parsed env maps (e.g. from two snapshots or branches)
// and returns a list of VersionDiff entries describing what changed.
func CompareVersions(base, target map[string]string) []VersionDiff {
	diffs := []VersionDiff{}

	for key, baseVal := range base {
		if targetVal, ok := target[key]; !ok {
			diffs = append(diffs, VersionDiff{
				Key:    key,
				From:   maskIfSecret(key, baseVal),
				To:     "",
				Change: "removed",
			})
		} else if baseVal != targetVal {
			diffs = append(diffs, VersionDiff{
				Key:    key,
				From:   maskIfSecret(key, baseVal),
				To:     maskIfSecret(key, targetVal),
				Change: "changed",
			})
		}
	}

	for key, targetVal := range target {
		if _, ok := base[key]; !ok {
			diffs = append(diffs, VersionDiff{
				Key:    key,
				From:   "",
				To:     maskIfSecret(key, targetVal),
				Change: "added",
			})
		}
	}

	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].Key < diffs[j].Key
	})

	return diffs
}

// FormatVersionDiff formats a slice of VersionDiff into a human-readable string.
func FormatVersionDiff(diffs []VersionDiff) string {
	if len(diffs) == 0 {
		return "No differences found."
	}

	var sb strings.Builder
	for _, d := range diffs {
		switch d.Change {
		case "added":
			sb.WriteString(fmt.Sprintf("+ %s=%s\n", d.Key, d.To))
		case "removed":
			sb.WriteString(fmt.Sprintf("- %s=%s\n", d.Key, d.From))
		case "changed":
			sb.WriteString(fmt.Sprintf("~ %s: %s -> %s\n", d.Key, d.From, d.To))
		}
	}
	return sb.String()
}

func maskIfSecret(key, value string) string {
	if IsSecret(key) {
		return MaskEntry(value)
	}
	return value
}
