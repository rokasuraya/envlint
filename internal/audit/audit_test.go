package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envlint/internal/audit"
	"github.com/user/envlint/internal/validator"
)

func sampleIssues() []validator.Issue {
	return []validator.Issue{
		{Key: "PORT", Message: "must be int", Severity: validator.SeverityError},
		{Key: "DEBUG", Message: "optional missing", Severity: validator.SeverityWarning},
	}
}

func TestNew_FieldsPopulated(t *testing.T) {
	entry := audit.New(".env", "schema.yaml", sampleIssues())

	if entry.EnvFile != ".env" {
		t.Errorf("EnvFile = %q, want .env", entry.EnvFile)
	}
	if entry.SchemaFile != "schema.yaml" {
		t.Errorf("SchemaFile = %q, want schema.yaml", entry.SchemaFile)
	}
	if entry.ErrorCount != 1 {
		t.Errorf("ErrorCount = %d, want 1", entry.ErrorCount)
	}
	if entry.WarnCount != 1 {
		t.Errorf("WarnCount = %d, want 1", entry.WarnCount)
	}
	if entry.Passed {
		t.Error("Passed should be false when errors > 0")
	}
	if entry.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

func TestNew_PassedWhenNoErrors(t *testing.T) {
	issue := validator.Issue{Key: "X", Message: "hint", Severity: validator.SeverityWarning}
	entry := audit.New(".env", "s.yaml", []validator.Issue{issue})
	if !entry.Passed {
		t.Error("Passed should be true when only warnings")
	}
}

func TestAppendAndReadAll_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	e1 := audit.New(".env", "schema.yaml", sampleIssues())
	e2 := audit.New(".env.prod", "schema.yaml", nil)

	if err := audit.Append(path, e1); err != nil {
		t.Fatalf("Append e1: %v", err)
	}
	if err := audit.Append(path, e2); err != nil {
		t.Fatalf("Append e2: %v", err)
	}

	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("len = %d, want 2", len(entries))
	}
	if entries[0].EnvFile != ".env" {
		t.Errorf("entry[0].EnvFile = %q", entries[0].EnvFile)
	}
	if entries[1].Passed != true {
		t.Error("entry[1] should have Passed=true")
	}
}

func TestReadAll_MissingFile(t *testing.T) {
	_, err := audit.ReadAll("/no/such/audit.log")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestAppend_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new.log")

	entry := audit.New(".env", "s.yaml", nil)
	entry.Timestamp = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	if err := audit.Append(path, entry); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
