package formatter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlint/internal/formatter"
	"github.com/user/envlint/internal/validator"
)

var sampleIssues = []validator.Issue{
	{Key: "DB_HOST", Severity: "error", Message: "required variable is missing"},
	{Key: "PORT", Severity: "error", Message: "value must be an integer"},
}

func TestWrite_TextNoIssues(t *testing.T) {
	var buf bytes.Buffer
	if err := formatter.Write(&buf, nil, formatter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "valid") {
		t.Errorf("expected success message, got: %s", buf.String())
	}
}

func TestWrite_TextWithIssues(t *testing.T) {
	var buf bytes.Buffer
	if err := formatter.Write(&buf, sampleIssues, formatter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected ERROR label, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
}

func TestWrite_JSONValid(t *testing.T) {
	var buf bytes.Buffer
	if err := formatter.Write(&buf, nil, formatter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var report formatter.Report
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !report.Valid {
		t.Error("expected valid=true when no issues")
	}
}

func TestWrite_JSONWithIssues(t *testing.T) {
	var buf bytes.Buffer
	if err := formatter.Write(&buf, sampleIssues, formatter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var report formatter.Report
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if report.Valid {
		t.Error("expected valid=false when issues present")
	}
	if len(report.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(report.Results))
	}
}

func TestWrite_GitHubActions(t *testing.T) {
	var buf bytes.Buffer
	if err := formatter.Write(&buf, sampleIssues, formatter.FormatGitHubActions); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "::error") {
		t.Errorf("expected GitHub Actions error annotation, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
}

func TestWrite_UnknownFormatFallsBackToText(t *testing.T) {
	var buf bytes.Buffer
	if err := formatter.Write(&buf, sampleIssues, formatter.Format("unknown")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Error("expected text fallback output")
	}
}
