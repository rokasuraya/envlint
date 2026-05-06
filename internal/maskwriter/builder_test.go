package maskwriter_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/envlint/internal/maskwriter"
)

func TestFromEnv_MasksSensitiveValues(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
		"APP_NAME":    "envlint",
	}

	var buf bytes.Buffer
	mw := maskwriter.FromEnv(&buf, env)

	_, _ = mw.Write([]byte("connecting with hunter2 to envlint\n"))
	got := buf.String()

	if contains(got, "hunter2") {
		t.Errorf("sensitive value leaked in output: %q", got)
	}
	if !contains(got, "[REDACTED]") {
		t.Errorf("expected [REDACTED] in output, got: %q", got)
	}
	// Non-sensitive value should pass through unchanged.
	if !contains(got, "envlint") {
		t.Errorf("non-sensitive value should not be masked, got: %q", got)
	}
}

func TestFromEnv_EmptyEnv(t *testing.T) {
	var buf bytes.Buffer
	mw := maskwriter.FromEnv(&buf, map[string]string{})

	_, _ = mw.Write([]byte("nothing to mask\n"))
	if got := buf.String(); got != "nothing to mask\n" {
		t.Errorf("got %q", got)
	}
}

func TestFromEnv_NilEnv(t *testing.T) {
	var buf bytes.Buffer
	mw := maskwriter.FromEnv(&buf, nil)

	_, _ = mw.Write([]byte("safe output\n"))
	if got := buf.String(); got != "safe output\n" {
		t.Errorf("got %q", got)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
