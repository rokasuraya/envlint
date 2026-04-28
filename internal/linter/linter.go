// Package linter orchestrates parsing, expansion, and validation
// of a .env file against a schema, returning a unified issue list.
package linter

import (
	"fmt"

	"github.com/your-org/envlint/internal/expander"
	"github.com/your-org/envlint/internal/parser"
	"github.com/your-org/envlint/internal/schema"
	"github.com/your-org/envlint/internal/validator"
)

// Config holds the options for a single lint run.
type Config struct {
	EnvFile    string
	SchemaFile string
	// ExpandVars controls whether variable references are expanded
	// before validation (e.g. BASE_URL=${HOST}/path).
	ExpandVars bool
}

// Run loads the .env file and schema, optionally expands variable
// references, validates every entry and returns the resulting issues.
func Run(cfg Config) ([]validator.Issue, error) {
	env, err := parser.ParseFile(cfg.EnvFile)
	if err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	sc, err := schema.LoadFile(cfg.SchemaFile)
	if err != nil {
		return nil, fmt.Errorf("load schema: %w", err)
	}

	if cfg.ExpandVars {
		exp := expander.New(env)
		expanded := make(map[string]string, len(env))
		for k, v := range env {
			expanded[k] = exp.Expand(v)
		}
		env = expanded
	}

	return validator.Validate(sc, env), nil
}
