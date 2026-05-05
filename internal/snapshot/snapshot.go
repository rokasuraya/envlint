// Package snapshot provides functionality to capture and persist the current
// state of a parsed .env file to disk, enabling later comparison via the
// differ package to detect configuration drift over time.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot holds a point-in-time capture of environment variable key/value pairs.
type Snapshot struct {
	CapturedAt time.Time         `json:"captured_at"`
	Source     string            `json:"source"`
	Vars       map[string]string `json:"vars"`
}

// New creates a new Snapshot from the provided vars map and source path.
func New(source string, vars map[string]string) *Snapshot {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Source:     source,
		Vars:       copy,
	}
}

// Save writes the snapshot as JSON to the given file path.
func Save(s *Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a previously saved snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open %q: %w", path, err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode %q: %w", path, err)
	}
	return &s, nil
}
