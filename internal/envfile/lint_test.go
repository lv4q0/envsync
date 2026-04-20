package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envsync/internal/envfile"
)

func writeTempEnvForLint(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestLint_CleanFile_NoIssues(t *testing.T) {
	path := writeTempEnvForLint(t, "APP_NAME=myapp\nPORT=8080\nDEBUG=false\n")
	issues, err := envfile.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %+v", len(issues), issues)
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	path := writeTempEnvForLint(t, "APP_NAME=myapp\nAPP_NAME=other\nPORT=8080\n")
	issues, err := envfile.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsLintCode(issues, "DUPLICATE_KEY") {
		t.Errorf("expected DUPLICATE_KEY issue, got: %+v", issues)
	}
}

func TestLint_TrailingWhitespace(t *testing.T) {
	path := writeTempEnvForLint(t, "APP_NAME=myapp  \nPORT=8080\n")
	issues, err := envfile.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsLintCode(issues, "TRAILING_WHITESPACE") {
		t.Errorf("expected TRAILING_WHITESPACE issue, got: %+v", issues)
	}
}

func TestLint_UnquotedSpecialChars(t *testing.T) {
	path := writeTempEnvForLint(t, "APP_NAME=my app name\nPORT=8080\n")
	issues, err := envfile.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsLintCode(issues, "UNQUOTED_SPACES") {
		t.Errorf("expected UNQUOTED_SPACES issue, got: %+v", issues)
	}
}

func TestLint_EmptySecretValue(t *testing.T) {
	path := writeTempEnvForLint(t, "API_SECRET=\nPORT=8080\n")
	issues, err := envfile.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsLintCode(issues, "EMPTY_SECRET") {
		t.Errorf("expected EMPTY_SECRET issue, got: %+v", issues)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	path := writeTempEnvForLint(t, "app_name=myapp\nPORT=8080\n")
	issues, err := envfile.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsLintCode(issues, "LOWERCASE_KEY") {
		t.Errorf("expected LOWERCASE_KEY issue, got: %+v", issues)
	}
}

func TestLint_MissingFile(t *testing.T) {
	_, err := envfile.Lint("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	// duplicate key + trailing whitespace
	path := writeTempEnvForLint(t, "PORT=8080  \nPORT=9090\nAPP=ok\n")
	issues, err := envfile.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d: %+v", len(issues), issues)
	}
}

// containsLintCode checks whether any lint issue has the given code.
func containsLintCode(issues []envfile.LintIssue, code string) bool {
	for _, issue := range issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}
