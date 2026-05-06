// Package summary produces a human-readable run summary that aggregates
// validation results, diff changes, and audit metadata into a single report.
package summary

import (
	"fmt"
	"io"
	"time"

	"github.com/user/envlint/internal/audit"
	"github.com/user/envlint/internal/differ"
	"github.com/user/envlint/internal/validator"
)

// Report holds all data needed to render a summary.
type Report struct {
	EnvFile    string
	SchemaFile string
	Profile    string
	RunAt      time.Time
	Issues     []validator.Issue
	Diff       []differ.Change
	Entry      *audit.Entry
}

// Write renders the summary to w in plain text.
func Write(w io.Writer, r Report) error {
	fmt.Fprintf(w, "envlint summary — %s\n", r.RunAt.Format(time.RFC3339))
	fmt.Fprintf(w, "  env    : %s\n", r.EnvFile)
	fmt.Fprintf(w, "  schema : %s\n", r.SchemaFile)
	if r.Profile != "" {
		fmt.Fprintf(w, "  profile: %s\n", r.Profile)
	}
	fmt.Fprintln(w)

	errors := validator.FilterBySeverity(r.Issues, validator.SeverityError)
	warnings := validator.FilterBySeverity(r.Issues, validator.SeverityWarning)

	fmt.Fprintf(w, "validation : %d error(s), %d warning(s)\n", len(errors), len(warnings))
	for _, iss := range r.Issues {
		marker := "WARN"
		if iss.Severity == validator.SeverityError {
			marker = "ERR "
		}
		fmt.Fprintf(w, "  [%s] %s: %s\n", marker, iss.Key, iss.Message)
	}

	if len(r.Diff) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "drift      : %s\n", differ.Summary(r.Diff))
	}

	if r.Entry != nil {
		passed := "no"
		if r.Entry.Passed {
			passed = "yes"
		}
		fmt.Fprintln(w)
		fmt.Fprintf(w, "audit      : run #%s passed=%s\n", r.Entry.ID, passed)
	}

	return nil
}

// Passed returns true when the report contains no error-level issues.
func (r Report) Passed() bool {
	return !validator.HasErrors(r.Issues)
}
