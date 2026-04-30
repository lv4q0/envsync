package envfile

import (
	"strings"
	"testing"
)

func TestCompareEnvMaps_Added(t *testing.T) {
	base := map[string]string{"A": "1"}
	target := map[string]string{"A": "1", "B": "2"}
	d := CompareEnvMaps(base, target)
	if _, ok := d.Added["B"]; !ok {
		t.Error("expected B to be added")
	}
	if len(d.Removed) != 0 || len(d.Changed) != 0 {
		t.Error("unexpected removed or changed entries")
	}
}

func TestCompareEnvMaps_Removed(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	target := map[string]string{"A": "1"}
	d := CompareEnvMaps(base, target)
	if _, ok := d.Removed["B"]; !ok {
		t.Error("expected B to be removed")
	}
}

func TestCompareEnvMaps_Changed(t *testing.T) {
	base := map[string]string{"A": "old"}
	target := map[string]string{"A": "new"}
	d := CompareEnvMaps(base, target)
	pair, ok := d.Changed["A"]
	if !ok {
		t.Fatal("expected A to be changed")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("unexpected change values: %v", pair)
	}
}

func TestCompareEnvMaps_Same(t *testing.T) {
	base := map[string]string{"A": "1"}
	target := map[string]string{"A": "1"}
	d := CompareEnvMaps(base, target)
	if _, ok := d.Same["A"]; !ok {
		t.Error("expected A to be same")
	}
	if len(d.Added)+len(d.Removed)+len(d.Changed) != 0 {
		t.Error("unexpected diff entries")
	}
}

func TestFormatEnvDiff_MasksSecrets(t *testing.T) {
	base := map[string]string{"SECRET_KEY": "old_secret"}
	target := map[string]string{"SECRET_KEY": "new_secret", "APP_NAME": "myapp"}
	d := CompareEnvMaps(base, target)
	out := FormatEnvDiff(d)
	if strings.Contains(out, "old_secret") || strings.Contains(out, "new_secret") {
		t.Error("secret values should be masked in output")
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Error("non-secret added key should appear unmasked")
	}
}

func TestFormatEnvDiff_ShowsSymbols(t *testing.T) {
	base := map[string]string{"OLD": "x"}
	target := map[string]string{"NEW": "y"}
	d := CompareEnvMaps(base, target)
	out := FormatEnvDiff(d)
	if !strings.Contains(out, "+ NEW=y") {
		t.Errorf("expected added line, got: %s", out)
	}
	if !strings.Contains(out, "- OLD=x") {
		t.Errorf("expected removed line, got: %s", out)
	}
}
