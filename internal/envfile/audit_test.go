package envfile

import (
	"strings"
	"testing"
)

func TestAuditLog_RecordAndSummary_NoEntries(t *testing.T) {
	log := &AuditLog{}
	summary := log.Summary()
	if summary != "No changes recorded.\n" {
		t.Errorf("expected empty message, got %q", summary)
	}
}

func TestAuditLog_RecordAdded(t *testing.T) {
	log := &AuditLog{}
	log.Record(ActionAdded, "APP_NAME", "", "myapp")
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if e.Action != ActionAdded || e.Key != "APP_NAME" || e.NewValue != "myapp" {
		t.Errorf("unexpected entry: %+v", e)
	}
	summary := log.Summary()
	if !strings.Contains(summary, "ADDED") || !strings.Contains(summary, "APP_NAME") {
		t.Errorf("summary missing expected content: %s", summary)
	}
}

func TestAuditLog_RecordChanged(t *testing.T) {
	log := &AuditLog{}
	log.Record(ActionChanged, "PORT", "3000", "4000")
	e := log.Entries[0]
	if e.OldValue != "3000" || e.NewValue != "4000" {
		t.Errorf("unexpected values: old=%q new=%q", e.OldValue, e.NewValue)
	}
	summary := log.Summary()
	if !strings.Contains(summary, "3000") || !strings.Contains(summary, "4000") {
		t.Errorf("summary missing values: %s", summary)
	}
}

func TestAuditLog_SecretValuesMasked(t *testing.T) {
	log := &AuditLog{}
	log.Record(ActionChanged, "DB_PASSWORD", "hunter2", "s3cr3t")
	e := log.Entries[0]
	if e.Secret != true {
		t.Error("expected Secret=true for DB_PASSWORD")
	}
	if e.OldValue != "***" || e.NewValue != "***" {
		t.Errorf("expected masked values, got old=%q new=%q", e.OldValue, e.NewValue)
	}
	summary := log.Summary()
	if strings.Contains(summary, "hunter2") || strings.Contains(summary, "s3cr3t") {
		t.Error("secret values should not appear in summary")
	}
}

func TestAuditLog_RecordRemoved(t *testing.T) {
	log := &AuditLog{}
	log.Record(ActionRemoved, "OLD_KEY", "oldval", "")
	summary := log.Summary()
	if !strings.Contains(summary, "REMOVED") || !strings.Contains(summary, "OLD_KEY") {
		t.Errorf("summary missing expected content: %s", summary)
	}
}

func TestAuditLog_MultipleEntries(t *testing.T) {
	log := &AuditLog{}
	log.Record(ActionAdded, "KEY1", "", "v1")
	log.Record(ActionRemoved, "KEY2", "v2", "")
	log.Record(ActionSynced, "KEY3", "old", "new")
	if len(log.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(log.Entries))
	}
	summary := log.Summary()
	if !strings.Contains(summary, "3 entries") {
		t.Errorf("summary should mention entry count: %s", summary)
	}
}
