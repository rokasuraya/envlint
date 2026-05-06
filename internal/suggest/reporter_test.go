package suggest

import (
	"bytes"
	"strings"
	"testing"
)

var sampleSchema = []string{
	"DATABASE_URL",
	"DATABASE_HOST",
	"APP_PORT",
	"APP_ENV",
	"LOG_LEVEL",
}

func TestWriteHints_NoUnknownKeys(t *testing.T) {
	var buf bytes.Buffer
	WriteHints(&buf, nil, sampleSchema, 3)
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestWriteHints_WithSuggestion(t *testing.T) {
	var buf bytes.Buffer
	WriteHints(&buf, []string{"DATABASE_UR"}, sampleSchema, 3)
	out := buf.String()
	if !strings.Contains(out, "DATABASE_URL") {
		t.Errorf("expected suggestion DATABASE_URL in output, got: %s", out)
	}
	if !strings.Contains(out, "did you mean") {
		t.Errorf("expected 'did you mean' hint, got: %s", out)
	}
}

func TestWriteHints_NoMatch(t *testing.T) {
	var buf bytes.Buffer
	WriteHints(&buf, []string{"ZZZZZ_UNKNOWN"}, sampleSchema, 3)
	out := buf.String()
	if !strings.Contains(out, "no similar keys found") {
		t.Errorf("expected 'no similar keys found', got: %s", out)
	}
}

func TestWriteHints_MultipleUnknown(t *testing.T) {
	var buf bytes.Buffer
	WriteHints(&buf, []string{"APP_PORTT", "LOG_LEVL"}, sampleSchema, 2)
	out := buf.String()
	if !strings.Contains(out, "APP_PORT") {
		t.Errorf("expected APP_PORT suggestion, got: %s", out)
	}
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Errorf("expected LOG_LEVEL suggestion, got: %s", out)
	}
}

func TestUnknownKeys_AllKnown(t *testing.T) {
	env := map[string]string{"DATABASE_URL": "x", "APP_PORT": "8080"}
	got := UnknownKeys(env, sampleSchema)
	if len(got) != 0 {
		t.Errorf("expected no unknown keys, got %v", got)
	}
}

func TestUnknownKeys_SomeUnknown(t *testing.T) {
	env := map[string]string{"DATABASE_URL": "x", "MYSTERY_VAR": "y"}
	got := UnknownKeys(env, sampleSchema)
	if len(got) != 1 || got[0] != "MYSTERY_VAR" {
		t.Errorf("expected [MYSTERY_VAR], got %v", got)
	}
}

func TestUnknownKeys_EmptyEnv(t *testing.T) {
	got := UnknownKeys(map[string]string{}, sampleSchema)
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}
