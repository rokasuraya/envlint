package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envlint/internal/snapshot"
)

func TestNew_CopiesVars(t *testing.T) {
	orig := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshot.New(".env", orig)

	// Mutate original — snapshot must be unaffected.
	orig["FOO"] = "mutated"

	if s.Vars["FOO"] != "bar" {
		t.Errorf("expected snapshot FOO=bar, got %q", s.Vars["FOO"])
	}
	if s.Source != ".env" {
		t.Errorf("unexpected source %q", s.Source)
	}
	if s.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	vars := map[string]string{
		"APP_ENV":  "production",
		"PORT":     "8080",
		"DB_HOST":  "localhost",
	}

	s := snapshot.New(".env.production", vars)
	s.CapturedAt = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	tmp := filepath.Join(t.TempDir(), "snap.json")

	if err := snapshot.Save(s, tmp); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(tmp)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Source != s.Source {
		t.Errorf("source mismatch: got %q, want %q", loaded.Source, s.Source)
	}
	if !loaded.CapturedAt.Equal(s.CapturedAt) {
		t.Errorf("time mismatch: got %v, want %v", loaded.CapturedAt, s.CapturedAt)
	}
	for k, v := range vars {
		if loaded.Vars[k] != v {
			t.Errorf("var %s: got %q, want %q", k, loaded.Vars[k], v)
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	s := snapshot.New(".env", map[string]string{"K": "V"})
	err := snapshot.Save(s, filepath.Join(t.TempDir(), "no", "such", "dir", "snap.json"))
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestLoad_MalformedJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(tmp, []byte("{not valid json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := snapshot.Load(tmp)
	if err == nil {
		t.Error("expected error for malformed JSON, got nil")
	}
}
