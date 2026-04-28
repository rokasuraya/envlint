package main

import (
	"os"

	"github.com/your-org/envlint/internal/formatter"
	"github.com/your-org/envlint/internal/linter"
	"github.com/your-org/envlint/internal/validator"
)

// run is the testable entry-point that wires flags → linter → formatter.
// It returns an exit code (0 = success, 1 = lint errors, 2 = usage/IO error).
func run(args []string) int {
	flags, err := parseFlags(args)
	if err != nil {
		printUsage()
		return 2
	}

	fmt := resolveFormat(flags.Format)

	issues, err := linter.Run(linter.Config{
		EnvFile:    flags.EnvFile,
		SchemaFile: flags.SchemaFile,
		ExpandVars: flags.ExpandVars,
	})
	if err != nil {
		_ = formatter.Write(os.Stderr, fmt, []validator.Issue{
			{Level: validator.LevelError, Key: "", Message: err.Error()},
		})
		return 2
	}

	if err := formatter.Write(os.Stdout, fmt, issues); err != nil {
		_, _ = os.Stderr.WriteString("formatter error: " + err.Error() + "\n")
		return 2
	}

	if validator.HasErrors(issues) {
		return 1
	}
	return 0
}
