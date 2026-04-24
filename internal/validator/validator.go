package validator

import (
	"fmt"

	"github.com/user/envlint/internal/schema"
)

// Result holds the outcome of a single variable validation.
type Result struct {
	Key     string
	Missing bool
	Errors  []string
}

// Report aggregates all validation results.
type Report struct {
	Results []Result
}

// HasErrors returns true if any result has errors or is missing.
func (r *Report) HasErrors() bool {
	for _, res := range r.Results {
		if res.Missing || len(res.Errors) > 0 {
			return true
		}
	}
	return false
}

// Validate checks the provided env map against the given schema.
func Validate(env map[string]string, s *schema.Schema) *Report {
	report := &Report{}

	for key, varSchema := range s.Vars {
		result := Result{Key: key}

		value, exists := env[key]
		if !exists || value == "" {
			if varSchema.Required {
				result.Missing = true
				result.Errors = append(result.Errors, fmt.Sprintf("required variable %q is missing or empty", key))
			}
			report.Results = append(report.Results, result)
			continue
		}

		if err := varSchema.Validate(value); err != nil {
			result.Errors = append(result.Errors, err.Error())
		}

		report.Results = append(report.Results, result)
	}

	return report
}
