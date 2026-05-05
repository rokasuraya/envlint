// Package audit provides structured logging of envlint validation runs,
// recording which files were checked, when, and what issues were found.
package audit

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/envlint/internal/validator"
)

// Entry represents a single audit record for one validation run.
type Entry struct {
	Timestamp  time.Time        `json:"timestamp"`
	EnvFile    string           `json:"env_file"`
	SchemaFile string           `json:"schema_file"`
	Issues     []validator.Issue `json:"issues"`
	ErrorCount int              `json:"error_count"`
	WarnCount  int              `json:"warn_count"`
	Passed     bool             `json:"passed"`
}

// New builds an Entry from the provided validation context.
func New(envFile, schemaFile string, issues []validator.Issue) Entry {
	var errors, warns int
	for _, iss := range issues {
		switch iss.Severity {
		case validator.SeverityError:
			errors++
		case validator.SeverityWarning:
			warns++
		}
	}
	return Entry{
		Timestamp:  time.Now().UTC(),
		EnvFile:    envFile,
		SchemaFile: schemaFile,
		Issues:     issues,
		ErrorCount: errors,
		WarnCount:  warns,
		Passed:     errors == 0,
	}
}

// Append serialises entry as a JSON line and appends it to the file at path.
// The file is created if it does not exist.
func Append(path string, entry Entry) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	return enc.Encode(entry)
}

// ReadAll reads every JSON-line entry from the file at path.
func ReadAll(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
