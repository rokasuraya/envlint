package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlint/internal/profile"
)

func writeTempProfile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "profiles*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	path := writeTempProfile(t, `
profiles:
  - name: production
    env_file: .env.production
    schema_file: schema.yaml
    description: Production environment
  - name: staging
    env_file: .env.staging
    schema_file: schema.yaml
`)
	cfg, err := profile.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(cfg.Profiles))
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := profile.LoadConfig(filepath.Join(t.TempDir(), "nope.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	path := writeTempProfile(t, ": bad: yaml: [")
	_, err := profile.LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestConfig_Get_Found(t *testing.T) {
	cfg := &profile.Config{
		Profiles: []profile.Profile{
			{Name: "ci", EnvFile: ".env.ci", SchemaFile: "schema.yaml"},
		},
	}
	p, err := cfg.Get("ci")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.EnvFile != ".env.ci" {
		t.Errorf("expected .env.ci, got %s", p.EnvFile)
	}
}

func TestConfig_Get_NotFound(t *testing.T) {
	cfg := &profile.Config{}
	_, err := cfg.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestConfig_Validate_DuplicateName(t *testing.T) {
	cfg := &profile.Config{
		Profiles: []profile.Profile{
			{Name: "prod", EnvFile: ".env", SchemaFile: "s.yaml"},
			{Name: "prod", EnvFile: ".env2", SchemaFile: "s.yaml"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected duplicate name error")
	}
}

func TestConfig_Validate_MissingEnvFile(t *testing.T) {
	cfg := &profile.Config{
		Profiles: []profile.Profile{
			{Name: "prod", SchemaFile: "s.yaml"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing env_file")
	}
}

func TestConfig_Names(t *testing.T) {
	cfg := &profile.Config{
		Profiles: []profile.Profile{
			{Name: "a", EnvFile: "a", SchemaFile: "s"},
			{Name: "b", EnvFile: "b", SchemaFile: "s"},
		},
	}
	names := cfg.Names()
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Errorf("unexpected names: %v", names)
	}
}
