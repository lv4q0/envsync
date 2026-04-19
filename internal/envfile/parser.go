package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair in an .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
}

// EnvFile holds all parsed entries from an .env file.
type EnvFile struct {
	Path    string
	Entries []Entry
	Index   map[string]*Entry
}

// Parse reads and parses an .env file at the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	env := &EnvFile{
		Path:  path,
		Index: make(map[string]*Entry),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, `"'`)
		e := Entry{Key: key, Value: val}
		env.Entries = append(env.Entries, e)
		env.Index[key] = &env.Entries[len(env.Entries)-1]
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan %s: %w", path, err)
	}
	return env, nil
}
