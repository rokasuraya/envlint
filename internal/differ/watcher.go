package differ

import (
	"fmt"
	"io"
	"time"

	"github.com/envlint/envlint/internal/snapshot"
)

// WatchResult holds the outcome of comparing the current env vars
// against a previously saved snapshot.
type WatchResult struct {
	SnapshotTime time.Time
	Diff         []Change
	HasChanges   bool
}

// CompareWithSnapshot loads a snapshot from path and diffs it against
// the provided current vars map. Returns a WatchResult summarising any
// changes, or an error if the snapshot cannot be read.
func CompareWithSnapshot(snapshotPath string, current map[string]string) (*WatchResult, error) {
	snap, err := snapshot.Load(snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("load snapshot %q: %w", snapshotPath, err)
	}

	changes := Diff(snap.Vars, current)

	return &WatchResult{
		SnapshotTime: snap.CreatedAt,
		Diff:         changes,
		HasChanges:   HasChanges(changes),
	}, nil
}

// WriteWatchReport writes a human-readable diff report to w.
// If there are no changes it prints a short confirmation line.
func WriteWatchReport(w io.Writer, result *WatchResult) {
	if !result.HasChanges {
		fmt.Fprintf(w, "No changes since snapshot taken at %s\n",
			result.SnapshotTime.Format(time.RFC3339))
		return
	}

	fmt.Fprintf(w, "Changes since snapshot taken at %s:\n",
		result.SnapshotTime.Format(time.RFC3339))
	WriteText(w, result.Diff)
}
