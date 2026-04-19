package sync

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envsync/internal/diff"
	"github.com/user/envsync/internal/envfile"
)

// Options controls sync behaviour.
type Options struct {
	DryRun    bool
	Overwrite bool
}

// Result holds the outcome of a sync operation.
type Result struct {
	Applied []string
	Skipped []string
}

// Sync applies changes from source to destination .env file based on the diff.
// When DryRun is true the destination file is never written.
func Sync(srcPath, dstPath string, opts Options) (*Result, error) {
	src, err := envfile.Parse(srcPath)
	if err != nil {
		return nil, fmt.Errorf("parse source: %w", err)
	}

	dst, err := envfile.Parse(dstPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("parse destination: %w", err)
	}

	changes := diff.Compare(src, dst)

	result := &Result{}

	// Build merged map starting from destination.
	merged := make(map[string]string, len(dst))
	for k, v := range dst {
		merged[k] = v
	}

	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			merged[c.Key] = src[c.Key]
			result.Applied = append(result.Applied, c.Key)
		case diff.Changed:
			if opts.Overwrite {
				merged[c.Key] = src[c.Key]
				result.Applied = append(result.Applied, c.Key)
			} else {
				result.Skipped = append(result.Skipped, c.Key)
			}
		case diff.Removed:
			// Keys present in dst but not src are left untouched.
			result.Skipped = append(result.Skipped, c.Key)
		}
	}

	if opts.DryRun {
		return result, nil
	}

	if err := write(dstPath, merged); err != nil {
		return nil, fmt.Errorf("write destination: %w", err)
	}

	return result, nil
}

// write serialises the map to a .env file.
func write(path string, entries map[string]string) error {
	var sb strings.Builder
	for k, v := range entries {
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}
	return os.WriteFile(path, []byte(sb.String()), 0o600)
}
