package main

import (
	"testing"
)

func TestParseFlags_Valid(t *testing.T) {
	cfg, err := parseFlags([]string{"-schema", "schema.yaml", "-env", ".env.test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SchemaPath != "schema.yaml" {
		t.Errorf("SchemaPath = %q, want %q", cfg.SchemaPath, "schema.yaml")
	}
	if cfg.EnvPath != ".env.test" {
		t.Errorf("EnvPath = %q, want %q", cfg.EnvPath, ".env.test")
	}
}

func TestParseFlags_DefaultEnv(t *testing.T) {
	cfg, err := parseFlags([]string{"-schema", "schema.yaml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.EnvPath != ".env" {
		t.Errorf("EnvPath = %q, want default %q", cfg.EnvPath, ".env")
	}
}

func TestParseFlags_MissingSchema(t *testing.T) {
	_, err := parseFlags([]string{"-env", ".env"})
	if err == nil {
		t.Fatal("expected error when -schema is missing")
	}
}

func TestParseFlags_UnknownFlag(t *testing.T) {
	_, err := parseFlags([]string{"-unknown", "value"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestParseFlags_EmptyArgs(t *testing.T) {
	_, err := parseFlags([]string{})
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}
