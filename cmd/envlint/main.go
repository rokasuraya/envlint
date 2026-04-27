package main

import (
	"fmt"
	"os"

	"github.com/user/envlint/internal/parser"
	"github.com/user/envlint/internal/schema"
	"github.com/user/envlint/internal/validator"
)

func main() {
	cfg, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		printUsage()
		os.Exit(2)
	}

	sch, err := schema.LoadFile(cfg.SchemaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load schema %q: %v\n", cfg.SchemaPath, err)
		os.Exit(2)
	}

	env, err := parser.ParseFile(cfg.EnvPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse env file %q: %v\n", cfg.EnvPath, err)
		os.Exit(2)
	}

	results := validator.Validate(sch, env)
	validator.PrintReport(os.Stdout, results)

	for _, r := range results {
		if r.Error != "" {
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: envlint -schema <schema.yaml> -env <.env>")
}
