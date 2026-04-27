package envfile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// RollbackEntry represents a single rollback point.
type RollbackEntry struct {
	Timestamp time.Time
	Label     string
	Path      string
}

// SaveRollbackPoint saves the current state of envFile as a rollback snapshot
// under the given rollback directory, tagged with an optional label.
func SaveRollbackPoint(envFile, rollbackDir, label string) (RollbackEntry, error) {
	if err := os.MkdirAll(rollbackDir, 0700); err != nil {
		return RollbackEntry{}, fmt.Errorf("rollback: create dir: %w", err)
	}

	data, err := os.ReadFile(envFile)
	if err != nil {
		return RollbackEntry{}, fmt.Errorf("rollback: read source: %w", err)
	}

	now := time.Now().UTC()
	safe := strings.ReplaceAll(label, " ", "_")
	filename := fmt.Sprintf("%d_%s.env", now.UnixNano(), safe)
	dest := filepath.Join(rollbackDir, filename)

	if err := os.WriteFile(dest, data, 0600); err != nil {
		return RollbackEntry{}, fmt.Errorf("rollback: write snapshot: %w", err)
	}

	return RollbackEntry{Timestamp: now, Label: label, Path: dest}, nil
}

// ListRollbackPoints returns all rollback entries in the directory, newest first.
func ListRollbackPoints(rollbackDir string) ([]RollbackEntry, error) {
	entries, err := os.ReadDir(rollbackDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("rollback: list: %w", err)
	}

	var points []RollbackEntry
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".env") {
			continue
		}
		parts := strings.SplitN(strings.TrimSuffix(e.Name(), ".env"), "_", 2)
		if len(parts) < 2 {
			continue
		}
		var ns int64
		fmt.Sscanf(parts[0], "%d", &ns)
		points = append(points, RollbackEntry{
			Timestamp: time.Unix(0, ns).UTC(),
			Label:     strings.ReplaceAll(parts[1], "_", " "),
			Path:      filepath.Join(rollbackDir, e.Name()),
		})
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].Timestamp.After(points[j].Timestamp)
	})
	return points, nil
}

// RestoreRollbackPoint copies the rollback snapshot back to targetFile.
func RestoreRollbackPoint(entry RollbackEntry, targetFile string) error {
	data, err := os.ReadFile(entry.Path)
	if err != nil {
		return fmt.Errorf("rollback: read point: %w", err)
	}
	if err := os.WriteFile(targetFile, data, 0600); err != nil {
		return fmt.Errorf("rollback: restore: %w", err)
	}
	return nil
}
