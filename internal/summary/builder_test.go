package summary_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlint/internal/audit"
	"github.com/user/envlint/internal/differ"
	"github.com/user/envlint/internal/summary"
	"github.com/user/envlint/internal/validator"
)

func TestBuilder_BasicFields(t *testing.T) {
	r := summary.NewBuilder(".env", "schema.yaml").Build()
	if r.EnvFile != ".env" {
		t.Errorf("EnvFile = %q, want .env", r.EnvFile)
	}
	if r.SchemaFile != "schema.yaml" {
		t.Errorf("SchemaFile = %q, want schema.yaml", r.SchemaFile)
	}
	if r.RunAt.IsZero() {
		t.Error("RunAt should not be zero")
	}
}

func TestBuilder_WithProfile(t *testing.T) {
	r := summary.NewBuilder(".env", "s.yaml").WithProfile("staging").Build()
	if r.Profile != "staging" {
		t.Errorf("Profile = %q, want staging", r.Profile)
	}
}

func TestBuilder_WithIssues(t *testing.T) {
	issues := []validator.Issue{
		{Key: "X", Message: "missing", Severity: validator.SeverityError},
	}
	r := summary.NewBuilder(".env", "s.yaml").WithIssues(issues).Build()
	if len(r.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(r.Issues))
	}
}

func TestBuilder_WithDiff(t *testing.T) {
	changes := []differ.Change{{Key: "FOO", Type: differ.Added}}
	r := summary.NewBuilder(".env", "s.yaml").WithDiff(changes).Build()
	if len(r.Diff) != 1 {
		t.Fatalf("expected 1 diff change, got %d", len(r.Diff))
	}
}

func TestBuilder_WithAuditEntry(t *testing.T) {
	e := &audit.Entry{ID: "xyz", Passed: false}
	r := summary.NewBuilder(".env", "s.yaml").WithAuditEntry(e).Build()
	if r.Entry == nil || r.Entry.ID != "xyz" {
		t.Error("expected audit entry to be attached")
	}
}

func TestBuilder_WithSnapshot_MissingFile(t *testing.T) {
	// Should not panic or error when snapshot file is absent.
	r := summary.NewBuilder(".env", "s.yaml").
		WithSnapshot(map[string]string{"A": "1"}, "/nonexistent/snap.json").
		Build()
	if len(r.Diff) != 0 {
		t.Errorf("expected no diff when snapshot missing, got %d", len(r.Diff))
	}
}

func TestBuilder_WithSnapshot_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	// Write a minimal snapshot manually.
	data, _ := json.Marshal(map[string]interface{}{
		"vars": map[string]string{"A": "old"},
	})
	_ = os.WriteFile(path, data, 0o644)

	r := summary.NewBuilder(".env", "s.yaml").
		WithSnapshot(map[string]string{"A": "new"}, path).
		Build()
	if len(r.Diff) == 0 {
		t.Error("expected diff changes when value changed")
	}
}
