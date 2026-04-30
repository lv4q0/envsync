package envfile

import (
	"testing"
)

func baseRenameEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestRename_BasicRename(t *testing.T) {
	entries, result, err := Rename(baseRenameEntries(), "PORT", "HTTP_PORT", RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Renamed {
		t.Fatalf("expected Renamed=true, reason: %s", result.Reason)
	}
	for _, e := range entries {
		if e.Key == "PORT" {
			t.Error("old key PORT should no longer exist")
		}
		if e.Key == "HTTP_PORT" && e.Value != "8080" {
			t.Errorf("expected value 8080, got %s", e.Value)
		}
	}
}

func TestRename_OldKeyNotFound(t *testing.T) {
	entries, result, err := Rename(baseRenameEntries(), "MISSING", "NEW_KEY", RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Renamed {
		t.Error("expected Renamed=false for missing key")
	}
	if len(entries) != 3 {
		t.Errorf("entries should be unchanged, got %d", len(entries))
	}
}

func TestRename_NewKeyExistsNoOverwrite(t *testing.T) {
	_, result, err := Rename(baseRenameEntries(), "PORT", "APP_NAME", RenameOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Renamed {
		t.Error("expected Renamed=false when new key exists and overwrite is off")
	}
}

func TestRename_NewKeyExistsWithOverwrite(t *testing.T) {
	entries, result, err := Rename(baseRenameEntries(), "PORT", "APP_NAME", RenameOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Renamed {
		t.Fatalf("expected Renamed=true, reason: %s", result.Reason)
	}
	// Only one APP_NAME should remain, with PORT's original value
	count := 0
	for _, e := range entries {
		if e.Key == "APP_NAME" {
			count++
			if e.Value != "8080" {
				t.Errorf("expected overwritten value 8080, got %s", e.Value)
			}
		}
	}
	if count != 1 {
		t.Errorf("expected exactly one APP_NAME, got %d", count)
	}
}

func TestRename_InvalidNewKey(t *testing.T) {
	_, _, err := Rename(baseRenameEntries(), "PORT", "INVALID KEY!", RenameOptions{})
	if err == nil {
		t.Error("expected error for invalid new key")
	}
}

func TestRename_EmptyNewKey(t *testing.T) {
	_, _, err := Rename(baseRenameEntries(), "PORT", "", RenameOptions{})
	if err == nil {
		t.Error("expected error for empty new key")
	}
}

func TestRename_PreservesOrder(t *testing.T) {
	entries, _, _ := Rename(baseRenameEntries(), "DB_PASSWORD", "DB_PASS", RenameOptions{})
	expected := []string{"APP_NAME", "DB_PASS", "PORT"}
	for i, e := range entries {
		if e.Key != expected[i] {
			t.Errorf("position %d: expected %s, got %s", i, expected[i], e.Key)
		}
	}
}
