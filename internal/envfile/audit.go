package envfile

import (
	"fmt"
	"time"
)

// AuditAction represents the type of change recorded in an audit log entry.
type AuditAction string

const (
	ActionAdded   AuditAction = "ADDED"
	ActionRemoved AuditAction = "REMOVED"
	ActionChanged AuditAction = "CHANGED"
	ActionSynced  AuditAction = "SYNCED"
)

// AuditEntry represents a single recorded change to an environment variable.
type AuditEntry struct {
	Timestamp time.Time
	Action    AuditAction
	Key       string
	OldValue  string
	NewValue  string
	Secret    bool
}

// AuditLog holds a collection of audit entries.
type AuditLog struct {
	Entries []AuditEntry
}

// Record appends a new entry to the audit log.
func (a *AuditLog) Record(action AuditAction, key, oldVal, newVal string) {
	secret := IsSecret(key)
	if secret {
		if oldVal != "" {
			oldVal = "***"
		}
		if newVal != "" {
			newVal = "***"
		}
	}
	a.Entries = append(a.Entries, AuditEntry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Key:       key,
		OldValue:  oldVal,
		NewValue:  newVal,
		Secret:    secret,
	})
}

// Filter returns a new AuditLog containing only entries that match the given action.
func (a *AuditLog) Filter(action AuditAction) *AuditLog {
	filtered := &AuditLog{}
	for _, e := range a.Entries {
		if e.Action == action {
			filtered.Entries = append(filtered.Entries, e)
		}
	}
	return filtered
}

// Summary returns a human-readable summary of all audit entries.
func (a *AuditLog) Summary() string {
	if len(a.Entries) == 0 {
		return "No changes recorded.\n"
	}
	out := fmt.Sprintf("Audit log (%d entries):\n", len(a.Entries))
	for _, e := range a.Entries {
		ts := e.Timestamp.Format(time.RFC3339)
		switch e.Action {
		case ActionAdded:
			out += fmt.Sprintf("  [%s] %s %s = %q\n", ts, e.Action, e.Key, e.NewValue)
		case ActionRemoved:
			out += fmt.Sprintf("  [%s] %s %s (was %q)\n", ts, e.Action, e.Key, e.OldValue)
		case ActionChanged, ActionSynced:
			out += fmt.Sprintf("  [%s] %s %s: %q -> %q\n", ts, e.Action, e.Key, e.OldValue, e.NewValue)
		}
	}
	return out
}
