package envfile

import (
	"encoding/json"
	"fmt"
	"os"
)

// TakeSnapshot creates a key-value map from a parsed env file.
func TakeSnapshot(entries []Entry) map[string]string {
	snap := make(map[string]string, len(entries))
	for _, e := range entries {
		snap[e.Key] = e.Value
	}
	return snap
}

// SaveSnapshot serializes a snapshot map to a JSON file.
func SaveSnapshot(path string, snap map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating snapshot file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("encoding snapshot: %w", err)
	}
	return nil
}

// LoadSnapshot reads a snapshot JSON file and returns the key-value map.
func LoadSnapshot(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening snapshot file: %w", err)
	}
	defer f.Close()

	var snap map[string]string
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("decoding snapshot: %w", err)
	}
	return snap, nil
}
