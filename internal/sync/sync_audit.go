package sync

import (
	"github.com/user/envsync/internal/envfile"
)

// AuditedSync wraps Sync and records all changes to an AuditLog.
// It returns the audit log alongside the sync result.
func AuditedSync(base, target map[string]string, opts Options) (map[string]string, *envfile.AuditLog, error) {
	log := &envfile.AuditLog{}

	result, err := Sync(base, target, opts)
	if err != nil {
		return nil, log, err
	}

	// Record additions and changes by comparing result to target.
	for k, newVal := range result {
		oldVal, existed := target[k]
		if !existed {
			log.Record(envfile.ActionAdded, k, "", newVal)
		} else if oldVal != newVal {
			log.Record(envfile.ActionChanged, k, oldVal, newVal)
		}
	}

	// Record keys present in base but absent from result (removals are not
	// performed by Sync, but we log keys that were in target and dropped).
	for k, oldVal := range target {
		if _, ok := result[k]; !ok {
			log.Record(envfile.ActionRemoved, k, oldVal, "")
		}
	}

	return result, log, nil
}
