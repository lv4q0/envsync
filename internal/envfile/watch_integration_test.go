package envfile_test

import (
	"os"
	"testing"
	"time"

	"envsync/internal/envfile"
)

func TestWatch_MultipleChanges(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/multi.env"
	if err := os.WriteFile(path, []byte("A=1\n"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	done := make(chan struct{})
	defer close(done)

	ch := envfile.Watch(path, 20*time.Millisecond, done)

	changes := 0
	for i := 0; i < 3; i++ {
		time.Sleep(40 * time.Millisecond)
		content := []byte("A=" + string(rune('1'+i+1)) + "\n")
		if err := os.WriteFile(path, content, 0644); err != nil {
			t.Fatalf("write iteration %d: %v", i, err)
		}
		select {
		case ev := <-ch:
			if ev.Err != nil {
				t.Fatalf("unexpected error: %v", ev.Err)
			}
			if ev.Changed {
				changes++
			}
		case <-time.After(300 * time.Millisecond):
			t.Errorf("timed out waiting for change %d", i+1)
		}
	}

	if changes < 2 {
		t.Errorf("expected at least 2 change events, got %d", changes)
	}
}

func TestWatch_StopsOnDone(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/stop.env"
	if err := os.WriteFile(path, []byte("KEY=val\n"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	done := make(chan struct{})
	ch := envfile.Watch(path, 20*time.Millisecond, done)

	close(done)

	select {
	case _, ok := <-ch:
		if ok {
			// drain any pending event; channel should close soon
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("channel did not close after done signal")
	}
}
