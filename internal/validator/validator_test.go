package validator_test

import (
	"testing"

	"github.com/user/envlint/internal/schema"
	"github.com/user/envlint/internal/validator"
)

func buildSchema(vars map[string]schema.VarSchema) *schema.Schema {
	return &schema.Schema{Vars: vars}
}

func TestValidate_AllPresent(t *testing.T) {
	s := buildSchema(map[string]schema.VarSchema{
		"APP_ENV": {Type: "string", Required: true},
		"PORT":    {Type: "int", Required: true},
	})
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}

	report := validator.Validate(env, s)
	if report.HasErrors() {
		t.Errorf("expected no errors, got results: %+v", report.Results)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	s := buildSchema(map[string]schema.VarSchema{
		"DATABASE_URL": {Type: "url", Required: true},
	})
	env := map[string]string{}

	report := validator.Validate(env, s)
	if !report.HasErrors() {
		t.Error("expected errors for missing required variable")
	}
	if len(report.Results) != 1 || !report.Results[0].Missing {
		t.Errorf("expected result to be marked missing, got: %+v", report.Results)
	}
}

func TestValidate_OptionalMissing(t *testing.T) {
	s := buildSchema(map[string]schema.VarSchema{
		"OPTIONAL_VAR": {Type: "string", Required: false},
	})
	env := map[string]string{}

	report := validator.Validate(env, s)
	if report.HasErrors() {
		t.Errorf("expected no errors for optional missing variable, got: %+v", report.Results)
	}
}

func TestValidate_InvalidType(t *testing.T) {
	s := buildSchema(map[string]schema.VarSchema{
		"PORT": {Type: "int", Required: true},
	})
	env := map[string]string{"PORT": "not-a-number"}

	report := validator.Validate(env, s)
	if !report.HasErrors() {
		t.Error("expected validation error for invalid int")
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	s := buildSchema(map[string]schema.VarSchema{
		"API_KEY":  {Type: "string", Required: true},
		"API_URL":  {Type: "url", Required: true},
		"LOG_LEVEL": {Type: "string", Required: false},
	})
	env := map[string]string{"API_URL": "not-a-url"}

	report := validator.Validate(env, s)
	if !report.HasErrors() {
		t.Error("expected errors")
	}
	errorCount := 0
	for _, r := range report.Results {
		if r.Missing || len(r.Errors) > 0 {
			errorCount++
		}
	}
	if errorCount < 2 {
		t.Errorf("expected at least 2 error results, got %d", errorCount)
	}
}
