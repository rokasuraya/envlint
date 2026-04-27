package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envlint/internal/formatter"
)

// formatFlag is a custom flag.Value that validates the --format option.
type formatFlag struct {
	value formatter.Format
}

var validFormats = []formatter.Format{
	formatter.FormatText,
	formatter.FormatJSON,
	formatter.FormatGitHubActions,
}

func (f *formatFlag) String() string {
	if f.value == "" {
		return string(formatter.FormatText)
	}
	return string(f.value)
}

func (f *formatFlag) Set(s string) error {
	candidate := formatter.Format(s)
	for _, v := range validFormats {
		if candidate == v {
			f.value = candidate
			return nil
		}
	}
	return fmt.Errorf("unknown format %q: choose one of text, json, github", s)
}

// registerFormatFlag registers --format on the given FlagSet and returns a
// pointer to the parsed Format value.
func registerFormatFlag(fs *flag.FlagSet) *formatter.Format {
	ff := &formatFlag{value: formatter.FormatText}
	fs.Var(ff, "format", "output format: text | json | github")
	return &ff.value
}

// resolveFormat reads ENVLINT_FORMAT from the environment as a fallback when
// the flag was not explicitly set.
func resolveFormat(flagVal formatter.Format) formatter.Format {
	if flagVal != "" {
		return flagVal
	}
	if env := os.Getenv("ENVLINT_FORMAT"); env != "" {
		return formatter.Format(env)
	}
	return formatter.FormatText
}
