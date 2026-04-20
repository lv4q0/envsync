package envfile

import (
	"strings"
	"testing"
)

func TestSerialize_DotEnv(t *testing.T) {
	env := map[string]string{"B": "2", "A": "1"}
	out := Serialize(env, FormatDotEnv)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "A=1" {
		t.Errorf("expected A=1, got %s", lines[0])
	}
	if lines[1] != "B=2" {
		t.Errorf("expected B=2, got %s", lines[1])
	}
}

func TestSerialize_Export(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	out := Serialize(env, FormatExport)
	if !strings.Contains(out, "export FOO=") {
		t.Errorf("expected export prefix, got: %s", out)
	}
}

func TestSerialize_JSON(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	out := Serialize(env, FormatJSON)
	if !strings.Contains(out, "{") || !strings.Contains(out, "}") {
		t.Errorf("expected JSON braces, got: %s", out)
	}
	if !strings.Contains(out, `"KEY"`) {
		t.Errorf("expected KEY in JSON, got: %s", out)
	}
}

func TestSerialize_EmptyMap(t *testing.T) {
	out := Serialize(map[string]string{}, FormatDotEnv)
	if strings.TrimSpace(out) != "" {
		t.Errorf("expected empty output, got: %q", out)
	}
}
