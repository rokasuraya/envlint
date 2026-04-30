package differ_test

import (
	"testing"

	"github.com/user/envlint/internal/differ"
)

func TestDiff_NoChanges(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	changes := differ.Diff(a, b)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestDiff_Added(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "NEW_KEY": "value"}
	changes := differ.Diff(a, b)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != differ.Added || changes[0].Key != "NEW_KEY" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiff_Removed(t *testing.T) {
	a := map[string]string{"FOO": "bar", "OLD_KEY": "val"}
	b := map[string]string{"FOO": "bar"}
	changes := differ.Diff(a, b)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != differ.Removed || changes[0].Key != "OLD_KEY" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiff_Changed(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	changes := differ.Diff(a, b)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Kind != differ.Changed || c.OldValue != "old" || c.NewValue != "new" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_Mixed(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"A": "9", "C": "3"}
	changes := differ.Diff(a, b)
	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(changes))
	}
}

func TestHasChanges(t *testing.T) {
	if differ.HasChanges(nil) {
		t.Error("expected false for nil slice")
	}
	if !differ.HasChanges([]differ.Change{{Key: "X", Kind: differ.Added}}) {
		t.Error("expected true for non-empty slice")
	}
}
