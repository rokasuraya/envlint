package summary

import (
	"time"

	"github.com/user/envlint/internal/audit"
	"github.com/user/envlint/internal/differ"
	"github.com/user/envlint/internal/snapshot"
	"github.com/user/envlint/internal/validator"
)

// Builder constructs a Report incrementally.
type Builder struct {
	report Report
}

// NewBuilder creates a Builder pre-filled with env/schema paths and the
// current timestamp.
func NewBuilder(envFile, schemaFile string) *Builder {
	return &Builder{
		report: Report{
			EnvFile:    envFile,
			SchemaFile: schemaFile,
			RunAt:      time.Now(),
		},
	}
}

// WithProfile sets the active profile name.
func (b *Builder) WithProfile(name string) *Builder {
	b.report.Profile = name
	return b
}

// WithIssues attaches validation issues.
func (b *Builder) WithIssues(issues []validator.Issue) *Builder {
	b.report.Issues = issues
	return b
}

// WithSnapshot computes the diff between current vars and a saved snapshot
// file at snapshotPath. If the file does not exist the diff is skipped.
func (b *Builder) WithSnapshot(vars map[string]string, snapshotPath string) *Builder {
	snap, err := snapshot.Load(snapshotPath)
	if err != nil {
		return b
	}
	b.report.Diff = differ.Diff(snap.Vars, vars)
	return b
}

// WithDiff attaches a pre-computed diff.
func (b *Builder) WithDiff(changes []differ.Change) *Builder {
	b.report.Diff = changes
	return b
}

// WithAuditEntry attaches the audit log entry produced for this run.
func (b *Builder) WithAuditEntry(e *audit.Entry) *Builder {
	b.report.Entry = e
	return b
}

// Build returns the assembled Report.
func (b *Builder) Build() Report {
	return b.report
}
