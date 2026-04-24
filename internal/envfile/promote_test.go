package envfile

import (
	"testing"
)

func TestPromote_AddsNewKeys(t *testing.T) {
	src := map[string]string{"NEW_KEY": "hello", "ANOTHER": "world"}
	dst := map[string]string{"EXISTING": "yes"}

	out, res, err := Promote(src, dst, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 2 {
		t.Errorf("expected 2 added, got %d", res.Added)
	}
	if out["NEW_KEY"] != "hello" || out["ANOTHER"] != "world" {
		t.Error("new keys not present in output")
	}
	if out["EXISTING"] != "yes" {
		t.Error("existing key should be preserved")
	}
}

func TestPromote_SkipsConflictWithoutOverwrite(t *testing.T) {
	src := map[string]string{"KEY": "new_val"}
	dst := map[string]string{"KEY": "old_val"}

	out, res, err := Promote(src, dst, PromoteOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
	if out["KEY"] != "old_val" {
		t.Errorf("expected old_val to be kept, got %s", out["KEY"])
	}
}

func TestPromote_OverwriteChangedKeys(t *testing.T) {
	src := map[string]string{"KEY": "new_val"}
	dst := map[string]string{"KEY": "old_val"}

	out, res, err := Promote(src, dst, PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Overwrite != 1 {
		t.Errorf("expected 1 overwrite, got %d", res.Overwrite)
	}
	if out["KEY"] != "new_val" {
		t.Errorf("expected new_val, got %s", out["KEY"])
	}
}

func TestPromote_OnlyKeys_FiltersKeys(t *testing.T) {
	src := map[string]string{"ALLOWED": "yes", "BLOCKED": "no"}
	dst := map[string]string{}

	out, res, err := Promote(src, dst, PromoteOptions{OnlyKeys: []string{"ALLOWED"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 1 {
		t.Errorf("expected 1 added, got %d", res.Added)
	}
	if _, ok := out["BLOCKED"]; ok {
		t.Error("BLOCKED key should not have been promoted")
	}
	if out["ALLOWED"] != "yes" {
		t.Error("ALLOWED key should be present")
	}
}

func TestPromote_NilSource_ReturnsError(t *testing.T) {
	_, _, err := Promote(nil, map[string]string{}, PromoteOptions{})
	if err == nil {
		t.Error("expected error for nil source, got nil")
	}
}

func TestPromote_SameValue_CountedAsSkipped(t *testing.T) {
	src := map[string]string{"KEY": "same"}
	dst := map[string]string{"KEY": "same"}

	_, res, err := Promote(src, dst, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped for identical value, got %d", res.Skipped)
	}
}
