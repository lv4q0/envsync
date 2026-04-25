package envfile

import (
	"os"
	"testing"
	"time"
)

func writeTempEnvForWatch(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestFileHash_Consistent(t *testing.T) {
	path := writeTempEnvForWatch(t, "KEY=value\n")
	h1, err := FileHash(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h2, err := FileHash(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h1 != h2 {
		t.Errorf("expected consistent hash, got %s vs %s", h1, h2)
	}
}

func TestFileHash_ChangesOnEdit(t *testing.T) {
	path := writeTempEnvForWatch(t, "KEY=value\n")
	h1, _ := FileHash(path)

	if err := os.WriteFile(path, []byte("KEY=changed\n"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	h2, _ := FileHash(path)
	if h1 == h2 {
		t.Error("expected hash to change after file edit")
	}
}

func TestFileHash_MissingFile(t *testing.T) {
	_, err := FileHash("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestWatch_DetectsChange(t *testing.T) {
	path := writeTempEnvForWatch(t, "KEY=initial\n")
	done := make(chan struct{})
	defer close(done)

	ch := Watch(path, 20*time.Millisecond, done)

	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte("KEY=updated\n"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.Err != nil {
			t.Fatalf("unexpected error: %v", ev.Err)
		}
		if !ev.Changed {
			t.Error("expected Changed=true")
		}
		if ev.Path != path {
			t.Errorf("expected path %s, got %s", path, ev.Path)
		}
	case <-time.After(300 * time.Millisecond):
		t.Error("timed out waiting for watch event")
	}
}

func TestWatch_NoEventWhenUnchanged(t *testing.T) {
	path := writeTempEnvForWatch(t, "KEY=stable\n")
	done := make(chan struct{})
	defer close(done)

	ch := Watch(path, 20*time.Millisecond, done)

	select {
	case ev := <-ch:
		if ev.Changed {
			t.Error("unexpected change event for unmodified file")
		}
	case <-time.After(100 * time.Millisecond):
		// expected: no event
	}
}
