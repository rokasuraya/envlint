package differ_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/envlint/envlint/internal/differ"
	"github.com/envlint/envlint/internal/snapshot"
)

// writeSnapshot is a test helper that persists a snapshot to a temp file
// and returns the file path.
func writeSnapshot(t *testing.T, vars map[string]string) string {
	t.Helper()
	snap := snapshot.New("test", vars)
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	if err := snapshot.Save(snap, path); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	return path
}

func TestCompareWithSnapshot_NoChanges(t *testing.T) {
	vars := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	path := writeSnapshot(t, vars)

	result, err := differ.CompareWithSnapshot(path, vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasChanges {
		t.Error("expected no changes")
	}
	if len(result.Diff) != 0 {
		t.Errorf("expected empty diff, got %d entries", len(result.Diff))
	}
}

func TestCompareWithSnapshot_DetectsChanges(t *testing.T) {
	old := map[string]string{"APP_ENV": "staging", "PORT": "8080"}
	current := map[string]string{"APP_ENV": "production", "DB_URL": "postgres://localhost/db"}
	path := writeSnapshot(t, old)

	result, err := differ.CompareWithSnapshot(path, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasChanges {
		t.Error("expected changes to be detected")
	}
	if len(result.Diff) == 0 {
		t.Error("expected non-empty diff")
	}
}

func TestCompareWithSnapshot_MissingFile(t *testing.T) {
	_, err := differ.CompareWithSnapshot("/nonexistent/snap.json", map[string]string{})
	if err == nil {
		t.Error("expected error for missing snapshot file")
	}
}

func TestWriteWatchReport_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	result := &differ.WatchResult{
		SnapshotTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		HasChanges:   false,
	}
	differ.WriteWatchReport(&buf, result)
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected 'No changes' in output, got: %s", buf.String())
	}
}

func TestWriteWatchReport_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	old := map[string]string{"KEY": "old"}
	new := map[string]string{"KEY": "new"}
	result := &differ.WatchResult{
		SnapshotTime: time.Now(),
		Diff:         differ.Diff(old, new),
		HasChanges:   true,
	}
	differ.WriteWatchReport(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "Changes since") {
		t.Errorf("expected 'Changes since' in output, got: %s", out)
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected KEY in diff output, got: %s", out)
	}
}

// Ensure snapshot JSON round-trips the CreatedAt field so watcher can display it.
func TestCompareWithSnapshot_SnapshotTimePreserved(t *testing.T) {
	vars := map[string]string{"X": "1"}
	path := writeSnapshot(t, vars)

	// Read raw JSON to confirm CreatedAt is stored.
	data, _ := os.ReadFile(path)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := raw["created_at"]; !ok {
		t.Error("snapshot JSON missing created_at field")
	}
}
