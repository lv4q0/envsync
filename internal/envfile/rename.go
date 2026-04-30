package envfile

import (
	"fmt"
	"strings"
)

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	OldKey  string
	NewKey  string
	Renamed bool
	Reason  string
}

// RenameOptions controls the behaviour of Rename.
type RenameOptions struct {
	// Overwrite allows the new key to replace an existing entry.
	Overwrite bool
}

// Rename renames oldKey to newKey inside entries, returning the updated slice
// and a RenameResult describing what happened.
//
// Rules:
//   - oldKey must exist.
//   - newKey must be a valid identifier (same rules as Validate).
//   - If newKey already exists and Overwrite is false, the rename is refused.
//   - The renamed entry keeps its original value and position.
func Rename(entries []Entry, oldKey, newKey string, opts RenameOptions) ([]Entry, RenameResult, error) {
	newKey = strings.TrimSpace(newKey)
	if newKey == "" {
		return nil, RenameResult{}, fmt.Errorf("new key must not be empty")
	}
	if !isValidKeyChars(newKey) {
		return nil, RenameResult{}, fmt.Errorf("new key %q contains invalid characters", newKey)
	}

	oldIdx := -1
	newIdx := -1
	for i, e := range entries {
		if e.Key == oldKey {
			oldIdx = i
		}
		if e.Key == newKey {
			newIdx = i
		}
	}

	if oldIdx == -1 {
		return entries, RenameResult{
			OldKey:  oldKey,
			NewKey:  newKey,
			Renamed: false,
			Reason:  fmt.Sprintf("key %q not found", oldKey),
		}, nil
	}

	if newIdx != -1 && !opts.Overwrite {
		return entries, RenameResult{
			OldKey:  oldKey,
			NewKey:  newKey,
			Renamed: false,
			Reason:  fmt.Sprintf("key %q already exists; use overwrite to replace", newKey),
		}, nil
	}

	updated := make([]Entry, 0, len(entries))
	for i, e := range entries {
		if i == newIdx {
			// drop the old occupant of newKey when overwriting
			continue
		}
		if i == oldIdx {
			e.Key = newKey
		}
		updated = append(updated, e)
	}

	return updated, RenameResult{
		OldKey:  oldKey,
		NewKey:  newKey,
		Renamed: true,
	}, nil
}
