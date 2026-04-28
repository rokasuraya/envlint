package linter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envlint/internal/linter"
	"github.com/your-org/envlint/internal/validator"
)

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

const minimalSchema = `vars:
  PORT:
    type: int
    required: true
  HOST:
    type: string
    required: true
`

func TestRun_NoIssues(t *testing.T) {
	env := writeTemp(t, ".env", "PORT=8080\nHOST=localhost\n")
	sc := writeTemp(t, "schema.yaml", minimalSchema)

	issues, err := linter.Run(linter.Config{EnvFile: env, SchemaFile: sc})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d: %v", len(issues), issues)
	}
}

func TestRun_MissingRequired(t *testing.T) {
	env := writeTemp(t, ".env", "PORT=8080\n")
	sc := writeTemp(t, "schema.yaml", minimalSchema)

	issues, err := linter.Run(linter.Config{EnvFile: env, SchemaFile: sc})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !validator.HasErrors(issues) {
		t.Error("expected at least one error issue")
	}
}

func TestRun_ExpandVars(t *testing.T) {
	env := writeTemp(t, ".env", "PORT=80\nHOST=example.com\n")
	sc := writeTemp(t, "schema.yaml", minimalSchema)

	issues, err := linter.Run(linter.Config{EnvFile: env, SchemaFile: sc, ExpandVars: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(issues))
	}
}

func TestRun_BadEnvFile(t *testing.T) {
	sc := writeTemp(t, "schema.yaml", minimalSchema)
	_, err := linter.Run(linter.Config{EnvFile: "/no/such/file", SchemaFile: sc})
	if err == nil {
		t.Error("expected error for missing env file")
	}
}

func TestRun_BadSchemaFile(t *testing.T) {
	env := writeTemp(t, ".env", "PORT=8080\n")
	_, err := linter.Run(linter.Config{EnvFile: env, SchemaFile: "/no/such/schema.yaml"})
	if err == nil {
		t.Error("expected error for missing schema file")
	}
}
