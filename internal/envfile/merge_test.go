package envfile

import (
	"testing"
)

func TestMerge_AddsNewKeys(t *testing.T) {
	base := map[string]string{"A": "1"}
	overlay := map[string]string{"B": "2"}
	r := Merge(base, overlay, StrategyKeepBase)
	if r.Merged["B"] != "2" {
		t.Errorf("expected B=2, got %s", r.Merged["B"])
	}
	if len(r.Added) != 1 || r.Added[0] != "B" {
		t.Errorf("expected Added=[B], got %v", r.Added)
	}
}

func TestMerge_KeepBase_SkipsConflict(t *testing.T) {
	base := map[string]string{"A": "old"}
	overlay := map[string]string{"A": "new"}
	r := Merge(base, overlay, StrategyKeepBase)
	if r.Merged["A"] != "old" {
		t.Errorf("expected A=old, got %s", r.Merged["A"])
	}
	if len(r.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(r.Skipped))
	}
}

func TestMerge_Overwrite_UpdatesConflict(t *testing.T) {
	base := map[string]string{"A": "old"}
	overlay := map[string]string{"A": "new"}
	r := Merge(base, overlay, StrategyOverwrite)
	if r.Merged["A"] != "new" {
		t.Errorf("expected A=new, got %s", r.Merged["A"])
	}
	if len(r.Updated) != 1 || r.Updated[0] != "A" {
		t.Errorf("expected Updated=[A], got %v", r.Updated)
	}
}

func TestMerge_SameValue_NoCounts(t *testing.T) {
	base := map[string]string{"A": "same"}
	overlay := map[string]string{"A": "same"}
	r := Merge(base, overlay, StrategyOverwrite)
	if len(r.Updated) != 0 || len(r.Skipped) != 0 || len(r.Added) != 0 {
		t.Errorf("expected no changes, got added=%v updated=%v skipped=%v", r.Added, r.Updated, r.Skipped)
	}
}

func TestMerge_PreservesBaseKeys(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	overlay := map[string]string{"C": "3"}
	r := Merge(base, overlay, StrategyKeepBase)
	if r.Merged["A"] != "1" || r.Merged["B"] != "2" {
		t.Errorf("base keys not preserved: %v", r.Merged)
	}
}
