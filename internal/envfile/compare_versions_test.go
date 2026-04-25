package envfile

import (
	"strings"
	"testing"
)

func TestCompareVersions_Added(t *testing.T) {
	base := map[string]string{"HOST": "localhost"}
	target := map[string]string{"HOST": "localhost", "PORT": "8080"}

	diffs := CompareVersions(base, target)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Key != "PORT" || diffs[0].Change != "added" || diffs[0].To != "8080" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestCompareVersions_Removed(t *testing.T) {
	base := map[string]string{"HOST": "localhost", "PORT": "8080"}
	target := map[string]string{"HOST": "localhost"}

	diffs := CompareVersions(base, target)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Key != "PORT" || diffs[0].Change != "removed" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestCompareVersions_Changed(t *testing.T) {
	base := map[string]string{"HOST": "localhost"}
	target := map[string]string{"HOST": "prod.example.com"}

	diffs := CompareVersions(base, target)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	d := diffs[0]
	if d.Change != "changed" || d.From != "localhost" || d.To != "prod.example.com" {
		t.Errorf("unexpected diff: %+v", d)
	}
}

func TestCompareVersions_NoChanges(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	diffs := CompareVersions(env, env)
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d", len(diffs))
	}
}

func TestCompareVersions_MasksSecrets(t *testing.T) {
	base := map[string]string{"API_SECRET": "old-secret"}
	target := map[string]string{"API_SECRET": "new-secret"}

	diffs := CompareVersions(base, target)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].From == "old-secret" || diffs[0].To == "new-secret" {
		t.Errorf("expected secret values to be masked, got from=%s to=%s", diffs[0].From, diffs[0].To)
	}
}

func TestFormatVersionDiff_Output(t *testing.T) {
	diffs := []VersionDiff{
		{Key: "PORT", From: "", To: "8080", Change: "added"},
		{Key: "HOST", From: "localhost", To: "prod", Change: "changed"},
		{Key: "OLD", From: "val", To: "", Change: "removed"},
	}
	out := FormatVersionDiff(diffs)
	if !strings.Contains(out, "+ PORT=8080") {
		t.Errorf("expected added line, got: %s", out)
	}
	if !strings.Contains(out, "~ HOST: localhost -> prod") {
		t.Errorf("expected changed line, got: %s", out)
	}
	if !strings.Contains(out, "- OLD=val") {
		t.Errorf("expected removed line, got: %s", out)
	}
}

func TestFormatVersionDiff_NoDiffs(t *testing.T) {
	out := FormatVersionDiff([]VersionDiff{})
	if out != "No differences found." {
		t.Errorf("expected no differences message, got: %s", out)
	}
}
