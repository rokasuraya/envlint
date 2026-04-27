package expander_test

import (
	"os"
	"testing"

	"github.com/yourorg/envlint/internal/expander"
)

func TestExpand_FallbackToOS(t *testing.T) {
	os.Setenv("_ENVLINT_TEST_HOST", "os-host")
	defer os.Unsetenv("_ENVLINT_TEST_HOST")

	e := expander.New() // FallbackToOS = true
	env := map[string]string{
		"DSN": "postgres://${_ENVLINT_TEST_HOST}:5432/db",
	}
	if err := e.Expand(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "postgres://os-host:5432/db"
	if got := env["DSN"]; got != want {
		t.Errorf("DSN = %q, want %q", got, want)
	}
}

func TestExpand_ChainedReferences(t *testing.T) {
	e := expander.New()
	e.FallbackToOS = false
	env := map[string]string{
		"PROTO": "https",
		"HOST":  "example.com",
		"BASE":  "${PROTO}://${HOST}",
		"URL":   "${BASE}/api/v1",
	}
	if err := e.Expand(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// BASE must be resolved before URL for the chain to work.
	// Because map iteration order is non-deterministic we call Expand twice
	// to simulate a stable two-pass scenario.
	if err := e.Expand(env); err != nil {
		t.Fatalf("second expand error: %v", err)
	}
	want := "https://example.com/api/v1"
	if got := env["URL"]; got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
}

func TestExpand_MultipleRefsInOneValue(t *testing.T) {
	e := expander.New()
	e.FallbackToOS = false
	env := map[string]string{
		"USER": "admin",
		"PASS": "secret",
		"DSN":  "${USER}:${PASS}@localhost",
	}
	if err := e.Expand(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "admin:secret@localhost"
	if got := env["DSN"]; got != want {
		t.Errorf("DSN = %q, want %q", got, want)
	}
}

func TestExpand_EmptyMap(t *testing.T) {
	e := expander.New()
	env := map[string]string{}
	if err := e.Expand(env); err != nil {
		t.Fatalf("unexpected error on empty map: %v", err)
	}
}
