package redactor_test

import (
	"testing"

	"github.com/your-org/envlint/internal/redactor"
)

func TestIsSensitive_MatchesDefaultKeys(t *testing.T) {
	r := redactor.New(nil)
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"SECRET_TOKEN", true},
		{"PRIVATE_KEY", true},
		{"AUTH_TOKEN", true},
		{"DATABASE_URL", false},
		{"PORT", false},
		{"APP_NAME", false},
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := r.IsSensitive(tc.key)
			if got != tc.want {
				t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
			}
		})
	}
}

func TestRedact_SensitiveKey(t *testing.T) {
	r := redactor.New(nil)
	got := r.Redact("DB_PASSWORD", "supersecret")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestRedact_NonSensitiveKey(t *testing.T) {
	r := redactor.New(nil)
	got := r.Redact("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected original value, got %q", got)
	}
}

func TestRedactMap(t *testing.T) {
	r := redactor.New(nil)
	env := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc123",
	}
	out := r.RedactMap(env)
	if out["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged")
	}
	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD should be redacted")
	}
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY should be redacted")
	}
}

func TestRedactMap_OriginalUnmodified(t *testing.T) {
	r := redactor.New(nil)
	env := map[string]string{"SECRET": "mysecret"}
	r.RedactMap(env)
	if env["SECRET"] != "mysecret" {
		t.Error("original map should not be modified")
	}
}

func TestNew_CustomKeys(t *testing.T) {
	r := redactor.New([]string{"INTERNAL"})
	if !r.IsSensitive("INTERNAL_ID") {
		t.Error("expected INTERNAL_ID to be sensitive with custom keys")
	}
	if r.IsSensitive("API_KEY") {
		t.Error("expected API_KEY to be non-sensitive with custom keys")
	}
}
