package diff

import (
	"fmt"

	"github.com/user/envsync/internal/envfile"
)

// Status represents the diff status of a key.
type Status string

const (
	StatusAdded   Status = "added"
	StatusRemoved Status = "removed"
	StatusChanged Status = "changed"
	StatusSame    Status = "same"
)

// Entry represents a single diff result for a key.
type Entry struct {
	Key      string
	Status   Status
	BaseVal  string
	TargetVal string
	Secret   bool
}

// String returns a human-readable representation of the diff entry.
func (e Entry) String() string {
	base := e.BaseVal
	target := e.TargetVal
	if e.Secret {
		base = envfile.MaskEntry(base)
		target = envfile.MaskEntry(target)
	}
	switch e.Status {
	case StatusAdded:
		return fmt.Sprintf("+ %s=%s", e.Key, target)
	case StatusRemoved:
		return fmt.Sprintf("- %s=%s", e.Key, base)
	case StatusChanged:
		return fmt.Sprintf("~ %s: %s -> %s", e.Key, base, target)
	default:
		return fmt.Sprintf("  %s=%s", e.Key, base)
	}
}

// Compare computes the diff between base and target env maps.
func Compare(base, target map[string]string) []Entry {
	var entries []Entry

	for k, bv := range base {
		if tv, ok := target[k]; ok {
			if bv == tv {
				entries = append(entries, Entry{Key: k, Status: StatusSame, BaseVal: bv, TargetVal: tv, Secret: envfile.IsSecret(k)})
			} else {
				entries = append(entries, Entry{Key: k, Status: StatusChanged, BaseVal: bv, TargetVal: tv, Secret: envfile.IsSecret(k)})
			}
		} else {
			entries = append(entries, Entry{Key: k, Status: StatusRemoved, BaseVal: bv, Secret: envfile.IsSecret(k)})
		}
	}

	for k, tv := range target {
		if _, ok := base[k]; !ok {
			entries = append(entries, Entry{Key: k, Status: StatusAdded, TargetVal: tv, Secret: envfile.IsSecret(k)})
		}
	}

	return entries
}
