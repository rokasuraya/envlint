package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envlint/internal/validator"
)

// Format is the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatGitHubActions Format = "github"
)

// Result is a serialisable representation of a single validation issue.
type Result struct {
	Key      string `json:"key"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// Report holds all results for JSON output.
type Report struct {
	Valid   bool     `json:"valid"`
	Results []Result `json:"results"`
}

// Write formats the validation issues to w in the requested format.
func Write(w io.Writer, issues []validator.Issue, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, issues)
	case FormatGitHubActions:
		return writeGitHubActions(w, issues)
	default:
		return writeText(w, issues)
	}
}

func writeText(w io.Writer, issues []validator.Issue) error {
	if len(issues) == 0 {
		_, err := fmt.Fprintln(w, "✓ All variables are valid.")
		return err
	}
	for _, iss := range issues {
		level := strings.ToUpper(iss.Severity)
		_, err := fmt.Fprintf(w, "[%s] %s: %s\n", level, iss.Key, iss.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, issues []validator.Issue) error {
	results := make([]Result, 0, len(issues))
	for _, iss := range issues {
		results = append(results, Result{
			Key:      iss.Key,
			Severity: iss.Severity,
			Message:  iss.Message,
		})
	}
	report := Report{
		Valid:   len(issues) == 0,
		Results: results,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func writeGitHubActions(w io.Writer, issues []validator.Issue) error {
	for _, iss := range issues {
		level := "warning"
		if iss.Severity == "error" {
			level = "error"
		}
		_, err := fmt.Fprintf(w, "::%s title=envlint::%s: %s\n", level, iss.Key, iss.Message)
		if err != nil {
			return err
		}
	}
	return nil
}
