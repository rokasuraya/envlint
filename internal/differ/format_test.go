package differ_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlint/internal/differ"
)

func TestWriteText_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	differ.WriteText(&buf, nil)
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-differences message, got: %q", buf.String())
	}
}

func TestWriteText_Added(t *testing.T) {
	changes := []differ.Change{
		{Key: "NEW_VAR", Kind: differ.Added, NewValue: "hello"},
	}
	var buf bytes.Buffer
	differ.WriteText(&buf, changes)
	out := buf.String()
	if !strings.Contains(out, "+ ") || !strings.Contains(out, "NEW_VAR") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestWriteText_Removed(t *testing.T) {
	changes := []differ.Change{
		{Key: "OLD_VAR", Kind: differ.Removed, OldValue: "bye"},
	}
	var buf bytes.Buffer
	differ.WriteText(&buf, changes)
	out := buf.String()
	if !strings.Contains(out, "- ") || !strings.Contains(out, "OLD_VAR") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestWriteText_Changed(t *testing.T) {
	changes := []differ.Change{
		{Key: "FOO", Kind: differ.Changed, OldValue: "a", NewValue: "b"},
	}
	var buf bytes.Buffer
	differ.WriteText(&buf, changes)
	out := buf.String()
	if !strings.Contains(out, "~ ") || !strings.Contains(out, "a -> b") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestSummary(t *testing.T) {
	changes := []differ.Change{
		{Kind: differ.Added},
		{Kind: differ.Added},
		{Kind: differ.Removed},
		{Kind: differ.Changed},
	}
	s := differ.Summary(changes)
	if s != "2 added, 1 removed, 1 changed" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestSummary_Empty(t *testing.T) {
	s := differ.Summary(nil)
	if s != "0 added, 0 removed, 0 changed" {
		t.Errorf("unexpected summary: %q", s)
	}
}
