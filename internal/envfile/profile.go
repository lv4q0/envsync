package envfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Profile represents a named environment profile (e.g. "staging", "production").
type Profile struct {
	Name    string
	Entries map[string]string
}

// LoadProfile loads a .env file for the given profile name from the base directory.
// It expects files named like: .env.staging, .env.production, etc.
func LoadProfile(dir, profileName string) (*Profile, error) {
	if profileName == "" {
		return nil, fmt.Errorf("profile name must not be empty")
	}

	filename := fmt.Sprintf(".env.%s", profileName)
	path := filepath.Join(dir, filename)

	entries, err := Parse(path)
	if err != nil {
		return nil, fmt.Errorf("load profile %q: %w", profileName, err)
	}

	return &Profile{
		Name:    profileName,
		Entries: entries,
	}, nil
}

// ListProfiles scans the given directory and returns all profile names
// found as .env.<name> files.
func ListProfiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("list profiles in %q: %w", dir, err)
	}

	var profiles []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, ".env.") {
			profile := strings.TrimPrefix(name, ".env.")
			if profile != "" {
				profiles = append(profiles, profile)
			}
		}
	}
	return profiles, nil
}

// DiffProfiles compares two profiles and returns keys that differ.
// Returns a map of key -> [baseValue, targetValue] for changed/added/removed keys.
func DiffProfiles(base, target *Profile) map[string][2]string {
	diffs := make(map[string][2]string)

	for k, v := range base.Entries {
		if tv, ok := target.Entries[k]; !ok {
			diffs[k] = [2]string{v, ""}
		} else if tv != v {
			diffs[k] = [2]string{v, tv}
		}
	}

	for k, v := range target.Entries {
		if _, ok := base.Entries[k]; !ok {
			diffs[k] = [2]string{"", v}
		}
	}

	return diffs
}
