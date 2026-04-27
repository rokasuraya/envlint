package expander_test

import (
	"testing"

	"github.com/yourorg/envlint/internal/expander"
)

func TestExpand_SimpleDollarBrace(t *testing.T) {
	e := expander.New()
	e.FallbackToOS = false
	env := map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://${HOST}:5432/db",
	}
	if err := e.Expand(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env["DSN"]; got != "postgres://localhost:5432/db" {
		t.Errorf("DSN = %q, want %q", got, "postgres://localhost:5432/db")
	}
}

func TestExpand_BareVariable(t *testing.T) {
	e := expander.New()
	e.FallbackToOS = false
	env := map[string]string{
		"PORT": "8080",
		"ADDR": "0.0.0.0:$PORT",
	}
	if err := e.Expand(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env["ADDR"]; got != "0.0.0.0:8080" {
		t.Errorf("ADDR = %q, want %q", got, "0.0.0.0:8080")
	}
}

func TestExpand_EscapedDollar(t *testing.T) {
	e := expander.New()
	e.FallbackToOS = false
	env := map[string]string{"PRICE": "$$10"}
	if err := e.Expand(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env["PRICE"]; got != "$10" {
		t.Errorf("PRICE = %q, want %q", got, "$10")
	}
}

func TestExpand_UnclosedBrace(t *testing.T) {
	e := expander.New()
	e.FallbackToOS = false
	env := map[string]string{"BAD": "${UNCLOSED"}
	if err := e.Expand(env); err == nil {
		t.Error("expected error for unclosed brace, got nil")
	}
}

func TestExpand_MissingVarNoFallback(t *testing.T) {
	e := expander.New()
	e.FallbackToOS = false
	env := map[string]string{"URL": "http://${MISSING}/path"}
	if err := e.Expand(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env["URL"]; got != "http:///path" {
		t.Errorf("URL = %q, want %q", got, "http:///path")
	}
}

func TestExpandValue_Standalone(t *testing.T) {
	e := expander.New()
	e.FallbackToOS = false
	env := map[string]string{"NAME": "world"}
	got, err := e.ExpandValue("hello ${NAME}!", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello world!" {
		t.Errorf("got %q, want %q", got, "hello world!")
	}
}

func TestExpand_NoReferences(t *testing.T) {
	e := expander.New()
	e.FallbackToOS = false
	env := map[string]string{"PLAIN": "just-a-value"}
	if err := e.Expand(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env["PLAIN"]; got != "just-a-value" {
		t.Errorf("PLAIN = %q, want %q", got, "just-a-value")
	}
}
