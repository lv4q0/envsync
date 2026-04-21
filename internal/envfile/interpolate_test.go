package envfile

import (
	"testing"
)

func TestInterpolate_BasicSubstitution(t *testing.T) {
	env := map[string]string{
		"HOST":     "localhost",
		"PORT":     "5432",
		"DATABASE_URL": "postgres://${HOST}:${PORT}/mydb",
	}

	errs := Interpolate(env)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}

	want := "postgres://localhost:5432/mydb"
	if got := env["DATABASE_URL"]; got != want {
		t.Errorf("DATABASE_URL = %q, want %q", got, want)
	}
}

func TestInterpolate_DollarStyleReference(t *testing.T) {
	env := map[string]string{
		"APP":  "envsync",
		"NAME": "app=$APP",
	}

	errs := Interpolate(env)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}

	want := "app=envsync"
	if got := env["NAME"]; got != want {
		t.Errorf("NAME = %q, want %q", got, want)
	}
}

func TestInterpolate_UndefinedReference(t *testing.T) {
	env := map[string]string{
		"URL": "http://${UNDEFINED_HOST}/path",
	}

	errs := Interpolate(env)
	if len(errs) == 0 {
		t.Fatal("expected error for undefined variable, got none")
	}

	ie, ok := errs[0].(*InterpolateError)
	if !ok {
		t.Fatalf("expected *InterpolateError, got %T", errs[0])
	}
	if ie.Key != "URL" || ie.Ref != "UNDEFINED_HOST" {
		t.Errorf("unexpected error fields: key=%q ref=%q", ie.Key, ie.Ref)
	}
}

func TestInterpolate_NoReferences(t *testing.T) {
	env := map[string]string{
		"PLAIN": "just-a-value",
		"NUM":   "42",
	}

	errs := Interpolate(env)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if env["PLAIN"] != "just-a-value" {
		t.Errorf("PLAIN changed unexpectedly: %q", env["PLAIN"])
	}
}

func TestInterpolate_MultipleErrors(t *testing.T) {
	env := map[string]string{
		"A": "${MISSING_A}",
		"B": "${MISSING_B}",
	}

	errs := Interpolate(env)
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}
}
