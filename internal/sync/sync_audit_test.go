package sync

import (
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func TestAuditedSync_RecordsAdditions(t *testing.T) {
	base := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	target := map[string]string{"APP_NAME": "myapp"}

	_, log, err := AuditedSync(base, target, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(log.Entries))
	}
	if log.Entries[0].Action != envfile.ActionAdded || log.Entries[0].Key != "PORT" {
		t.Errorf("unexpected entry: %+v", log.Entries[0])
	}
}

func TestAuditedSync_RecordsChanges(t *testing.T) {
	base := map[string]string{"PORT": "9090"}
	target := map[string]string{"PORT": "8080"}

	_, log, err := AuditedSync(base, target, Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if e.Action != envfile.ActionChanged {
		t.Errorf("expected CHANGED, got %s", e.Action)
	}
	if e.OldValue != "8080" || e.NewValue != "9090" {
		t.Errorf("unexpected values: old=%q new=%q", e.OldValue, e.NewValue)
	}
}

func TestAuditedSync_SecretsMaskedInLog(t *testing.T) {
	base := map[string]string{"DB_SECRET": "newpass"}
	target := map[string]string{"DB_SECRET": "oldpass"}

	_, log, err := AuditedSync(base, target, Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(log.Entries) == 0 {
		t.Fatal("expected at least one audit entry")
	}
	e := log.Entries[0]
	if e.OldValue != "***" || e.NewValue != "***" {
		t.Errorf("secret values should be masked, got old=%q new=%q", e.OldValue, e.NewValue)
	}
}

func TestAuditedSync_NoChanges_EmptyLog(t *testing.T) {
	base := map[string]string{"KEY": "value"}
	target := map[string]string{"KEY": "value"}

	_, log, err := AuditedSync(base, target, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Entries) != 0 {
		t.Errorf("expected no audit entries, got %d", len(log.Entries))
	}
}
