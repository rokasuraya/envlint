package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func writeProfilesFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "profiles.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestResolveProfile_NoFlag(t *testing.T) {
	env, schema, err := resolveProfile(profileFlags{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env != "" || schema != "" {
		t.Errorf("expected empty strings when no profile set")
	}
}

func TestResolveProfile_Valid(t *testing.T) {
	path := writeProfilesFile(t, `
profiles:
  - name: ci
    env_file: .env.ci
    schema_file: ci.schema.yaml
`)
	env, schema, err := resolveProfile(profileFlags{ProfileName: "ci", ProfilesFile: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env != ".env.ci" {
		t.Errorf("expected .env.ci, got %s", env)
	}
	if schema != "ci.schema.yaml" {
		t.Errorf("expected ci.schema.yaml, got %s", schema)
	}
}

func TestResolveProfile_NotFound(t *testing.T) {
	path := writeProfilesFile(t, `
profiles:
  - name: prod
    env_file: .env
    schema_file: s.yaml
`)
	_, _, err := resolveProfile(profileFlags{ProfileName: "missing", ProfilesFile: path})
	if err == nil {
		t.Fatal("expected error for unknown profile")
	}
}

func TestResolveProfile_MissingFile(t *testing.T) {
	_, _, err := resolveProfile(profileFlags{
		ProfileName:  "ci",
		ProfilesFile: filepath.Join(t.TempDir(), "nope.yaml"),
	})
	if err == nil {
		t.Fatal("expected error for missing profiles file")
	}
}

func TestRegisterProfileFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var pf profileFlags
	registerProfileFlags(fs, &pf)
	if err := fs.Parse([]string{"--profile", "staging", "--profiles-file", "custom.yaml"}); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if pf.ProfileName != "staging" {
		t.Errorf("expected staging, got %s", pf.ProfileName)
	}
	if pf.ProfilesFile != "custom.yaml" {
		t.Errorf("expected custom.yaml, got %s", pf.ProfilesFile)
	}
}
