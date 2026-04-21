package cmd

import (
	"os"
	"strings"
	"testing"
)

func writeInterpolateEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "interpolate-*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestInterpolateCmd_ResolvesReferences(t *testing.T) {
	path := writeInterpolateEnv(t, "HOST=localhost\nPORT=5432\nURL=postgres://${HOST}:${PORT}/db\n")

	out, err := captureOutput(func() {
		rootCmd.SetArgs([]string{"interpolate", path})
		_ = rootCmd.Execute()
	})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out, "postgres://localhost:5432/db") {
		t.Errorf("expected resolved URL in output, got:\n%s", out)
	}
}

func TestInterpolateCmd_UndefinedVar_ReturnsError(t *testing.T) {
	path := writeInterpolateEnv(t, "URL=http://${GHOST_HOST}/api\n")

	var exitErr error
	_, _ = captureOutput(func() {
		rootCmd.SetArgs([]string{"interpolate", path})
		exitErr = rootCmd.Execute()
	})

	if exitErr == nil {
		t.Error("expected error for undefined variable reference, got nil")
	}
}

func TestInterpolateCmd_JSONFormat(t *testing.T) {
	path := writeInterpolateEnv(t, "BASE=hello\nGREET=${BASE}_world\n")

	out, err := captureOutput(func() {
		rootCmd.SetArgs([]string{"interpolate", "--format", "json", path})
		_ = rootCmd.Execute()
	})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out, "hello_world") {
		t.Errorf("expected resolved value in JSON output, got:\n%s", out)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON output, got:\n%s", out)
	}
}
