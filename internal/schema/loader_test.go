package schema

import (
	"testing"
)

const validYAML = `
vars:
  PORT:
    required: true
    type: int
    description: HTTP port
  DEBUG:
    type: bool
    allow_empty: true
  API_URL:
    required: true
    type: url
  APP_NAME:
    type: string
    pattern: "^[a-z]+$"
`

func TestLoad_ValidYAML(t *testing.T) {
	s, err := Load([]byte(validYAML))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Vars) != 4 {
		t.Errorf("expected 4 vars, got %d", len(s.Vars))
	}

	port, ok := s.Vars["PORT"]
	if !ok {
		t.Fatal("PORT not found in schema")
	}
	if !port.Required {
		t.Error("PORT should be required")
	}
	if port.Type != TypeInt {
		t.Errorf("PORT type: expected int, got %s", port.Type)
	}

	debug := s.Vars["DEBUG"]
	if !debug.AllowEmpty {
		t.Error("DEBUG should allow empty")
	}

	app := s.Vars["APP_NAME"]
	if app.Pattern != "^[a-z]+$" {
		t.Errorf("unexpected pattern: %q", app.Pattern)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	_, err := Load([]byte(":::invalid yaml:::"))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoad_EmptyVars(t *testing.T) {
	s, err := Load([]byte("vars: {}"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Vars) != 0 {
		t.Errorf("expected 0 vars, got %d", len(s.Vars))
	}
}

func TestLoad_DefaultTypeIsString(t *testing.T) {
	yamlData := "vars:\n  FOO:\n    required: true\n"
	s, err := Load([]byte(yamlData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Vars["FOO"].Type != TypeString {
		t.Errorf("expected default type string, got %s", s.Vars["FOO"].Type)
	}
}
