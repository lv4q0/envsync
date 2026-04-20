package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of an env file's contents.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Entries   map[string]string `json:"entries"`
}

// TakeSnapshot reads the given env file and returns a Snapshot.
func TakeSnapshot(path string) (*Snapshot, error) {
	entries, err := Parse(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: failed to parse %q: %w", path, err)
	}
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Source:    path,
		Entries:   entries,
	}, nil
}

// SaveSnapshot writes a Snapshot to disk as JSON.
func SaveSnapshot(snap *Snapshot, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("snapshot: cannot create file %q: %w", dest, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encode error: %w", err)
	}
	return nil
}

// LoadSnapshot reads a previously saved Snapshot from a JSON file.
func LoadSnapshot(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: cannot open %q: %w", path, err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshot: decode error: %w", err)
	}
	return &snap, nil
}
