package envfile

import (
	"strings"
	"testing"
)

func TestDiffSnapshots_Added(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	current := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	diff := DiffSnapshots(base, current)

	if len(diff.Added) != 1 || diff.Added["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY to be added, got %v", diff.Added)
	}
	if len(diff.Removed) != 0 || len(diff.Changed) != 0 {
		t.Errorf("unexpected removed or changed entries")
	}
}

func TestDiffSnapshots_Removed(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD_KEY": "gone"}
	current := map[string]string{"FOO": "bar"}

	diff := DiffSnapshots(base, current)

	if len(diff.Removed) != 1 || diff.Removed["OLD_KEY"] != "gone" {
		t.Errorf("expected OLD_KEY to be removed, got %v", diff.Removed)
	}
}

func TestDiffSnapshots_Changed(t *testing.T) {
	base := map[string]string{"FOO": "old"}
	current := map[string]string{"FOO": "new"}

	diff := DiffSnapshots(base, current)

	if len(diff.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(diff.Changed))
	}
	if diff.Changed["FOO"].Old != "old" || diff.Changed["FOO"].New != "new" {
		t.Errorf("unexpected change: %+v", diff.Changed["FOO"])
	}
}

func TestDiffSnapshots_NoChanges(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	current := map[string]string{"FOO": "bar"}

	diff := DiffSnapshots(base, current)

	if len(diff.Added) != 0 || len(diff.Removed) != 0 || len(diff.Changed) != 0 {
		t.Errorf("expected no changes")
	}
}

func TestFormatSnapshotDiff_MasksSecrets(t *testing.T) {
	base := map[string]string{"DB_PASSWORD": "old_secret"}
	current := map[string]string{"DB_PASSWORD": "new_secret", "API_KEY": "abc"}

	diff := DiffSnapshots(base, current)
	output := FormatSnapshotDiff(diff, true)

	if strings.Contains(output, "old_secret") || strings.Contains(output, "new_secret") || strings.Contains(output, "abc") {
		t.Errorf("expected secrets to be masked, got:\n%s", output)
	}
	if !strings.Contains(output, "***") {
		t.Errorf("expected masked output, got:\n%s", output)
	}
}

func TestFormatSnapshotDiff_NoChangesMessage(t *testing.T) {
	diff := SnapshotDiff{
		Added:   map[string]string{},
		Removed: map[string]string{},
		Changed: map[string]SnapshotChange{},
	}
	output := FormatSnapshotDiff(diff, false)
	if !strings.Contains(output, "No changes") {
		t.Errorf("expected no-changes message, got: %s", output)
	}
}
