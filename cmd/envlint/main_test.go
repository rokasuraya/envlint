package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestRun_ValidEnv(t *testing.T) {
	dir := t.TempDir()

	schemaContent := `vars:
  PORT:
    type: int
    required: true
  APP_NAME:
    type: string
    required: true
`
	envContent := "PORT=8080\nAPP_NAME=myapp\n"

	schemaPath := writeFile(t, dir, "schema.yaml", schemaContent)
	envPath := writeFile(t, dir, ".env", envContent)

	// Verify flags parse correctly for this scenario.
	cfg, err := parseFlags([]string{"-schema", schemaPath, "-env", envPath})
	if err != nil {
		t.Fatalf("parseFlags: %v", err)
	}
	if cfg.SchemaPath != schemaPath {
		t.Errorf("schema path mismatch")
	}
}

func TestRun_MissingEnvFile(t *testing.T) {
	dir := t.TempDir()
	schemaPath := writeFile(t, dir, "schema.yaml", "vars:\n  X:\n    type: string\n")

	cfg, err := parseFlags([]string{"-schema", schemaPath, "-env", filepath.Join(dir, "nonexistent.env")})
	if err != nil {
		t.Fatalf("unexpected flag error: %v", err)
	}
	// Ensure the env path is preserved even if the file doesn't exist yet.
	if cfg.EnvPath == "" {
		t.Error("EnvPath should not be empty")
	}
}
