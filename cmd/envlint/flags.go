package main

import (
	"errors"
	"flag"
	"io"
)

// Config holds the parsed CLI flags.
type Config struct {
	SchemaPath string
	EnvPath    string
}

// parseFlags parses the provided argument slice and returns a Config.
func parseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("envlint", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	schemaPath := fs.String("schema", "", "path to the schema YAML file (required)")
	envPath := fs.String("env", ".env", "path to the .env file (default: .env)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *schemaPath == "" {
		return nil, errors.New("-schema flag is required")
	}

	return &Config{
		SchemaPath: *schemaPath,
		EnvPath:    *envPath,
	}, nil
}
