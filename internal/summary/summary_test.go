package summary_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/envlint/internal/audit"
	"github.com/user/envlint/internal/differ"
	"github.com/user/envlint/internal/summary"
	"github.com/user/envlint/internal/validator"
)

func baseReport() summary.Report {
	return summary.Report{
		EnvFile:    ".env",
		SchemaFile: "schema.yaml",
		RunAt:      time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestWrite_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	r := baseReport()
	if err := summary.Write(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "0 error(s)") {
		t.Errorf("expected 0 errors in output, got:\n%s", out)
	}
}

func TestWrite_WithErrors(t *testing.T) {
	var buf bytes.Buffer
	r := baseReport()
	r.Issues = []validator.Issue{
		{Key: "DB_URL", Message: "missing required variable", Severity: validator.SeverityError},
	}
	_ = summary.Write(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "ERR ") {
		t.Errorf("expected ERR marker, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected DB_URL in output, got:\n%s", out)
	}
}

func TestWrite_WithDiff(t *testing.T) {
	var buf bytes.Buffer
	r := baseReport()
	r.Diff = []differ.Change{
		{Key: "NEW_VAR", Type: differ.Added},
	}
	_ = summary.Write(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "drift") {
		t.Errorf("expected drift section, got:\n%s", out)
	}
}

func TestWrite_WithAuditEntry(t *testing.T) {
	var buf bytes.Buffer
	r := baseReport()
	r.Entry = &audit.Entry{ID: "abc123", Passed: true}
	_ = summary.Write(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "abc123") {
		t.Errorf("expected audit ID in output, got:\n%s", out)
	}
	if !strings.Contains(out, "passed=yes") {
		t.Errorf("expected passed=yes, got:\n%s", out)
	}
}

func TestReport_Passed(t *testing.T) {
	r := baseReport()
	if !r.Passed() {
		t.Error("empty report should pass")
	}
	r.Issues = []validator.Issue{
		{Severity: validator.SeverityError, Message: "bad"},
	}
	if r.Passed() {
		t.Error("report with errors should not pass")
	}
}

func TestWrite_ProfileLine(t *testing.T) {
	var buf bytes.Buffer
	r := baseReport()
	r.Profile = "production"
	_ = summary.Write(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "production") {
		t.Errorf("expected profile name in output, got:\n%s", out)
	}
}
