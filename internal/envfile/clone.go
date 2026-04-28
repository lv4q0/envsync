package envfile

import (
	"fmt"
	"os"
	"path/filepath"
)

// CloneOptions controls how an env file is cloned to a target path.
type CloneOptions struct {
	// Overwrite replaces the target file if it already exists.
	Overwrite bool
	// StripSecrets replaces secret values with empty strings in the clone.
	StripSecrets bool
	// Profile, if non-empty, writes the clone as a named profile variant.
	Profile string
}

// CloneResult summarises what happened during a clone operation.
type CloneResult struct {
	Destination string
	KeysWritten int
	SecretsStripped int
}

// Clone reads the env file at src, optionally strips secret values, and writes
// the result to dest according to opts.
func Clone(src, dest string, opts CloneOptions) (*CloneResult, error) {
	if src == "" {
		return nil, fmt.Errorf("clone: source path must not be empty")
	}
	if dest == "" {
		return nil, fmt.Errorf("clone: destination path must not be empty")
	}

	if opts.Profile != "" {
		dir := filepath.Dir(dest)
		base := filepath.Base(dest)
		dest = filepath.Join(dir, fmt.Sprintf("%s.%s", base, opts.Profile))
	}

	if !opts.Overwrite {
		if _, err := os.Stat(dest); err == nil {
			return nil, fmt.Errorf("clone: destination %q already exists (use overwrite to replace)", dest)
		}
	}

	entries, err := Parse(src)
	if err != nil {
		return nil, fmt.Errorf("clone: parse source: %w", err)
	}

	result := &CloneResult{Destination: dest}
	stripped := 0

	if opts.StripSecrets {
		for i, e := range entries {
			if IsSecret(e.Key) {
				entries[i].Value = ""
				stripped++
			}
		}
	}

	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}

	if err := Serialize(m, dest, "dotenv"); err != nil {
		return nil, fmt.Errorf("clone: write destination: %w", err)
	}

	result.KeysWritten = len(entries)
	result.SecretsStripped = stripped
	return result, nil
}
