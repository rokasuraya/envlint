package validator

import (
	"testing"

	"github.com/user/envlint/internal/schema"
)

func TestValidateValue_Int(t *testing.T) {
	vs := schema.VarSchema{Type: "int"}
	if errs := validateValue("PORT", "8080", vs); len(errs) != 0 {
		t.Errorf("expected no errors for valid int, got %v", errs)
	}
	if errs := validateValue("PORT", "abc", vs); len(errs) == 0 {
		t.Error("expected error for non-int value")
	}
}

func TestValidateValue_Bool(t *testing.T) {
	vs := schema.VarSchema{Type: "bool"}
	for _, v := range []string{"true", "false", "1", "0", "True", "FALSE"} {
		if errs := validateValue("FLAG", v, vs); len(errs) != 0 {
			t.Errorf("expected no errors for %q, got %v", v, errs)
		}
	}
	if errs := validateValue("FLAG", "yes", vs); len(errs) == 0 {
		t.Error("expected error for 'yes' as bool")
	}
}

func TestValidateValue_URL(t *testing.T) {
	vs := schema.VarSchema{Type: "url"}
	if errs := validateValue("SITE", "https://example.com", vs); len(errs) != 0 {
		t.Errorf("expected no errors for valid URL, got %v", errs)
	}
	if errs := validateValue("SITE", "not-a-url", vs); len(errs) == 0 {
		t.Error("expected error for invalid URL")
	}
}

func TestValidateValue_Email(t *testing.T) {
	vs := schema.VarSchema{Type: "email"}
	if errs := validateValue("ADMIN", "admin@example.com", vs); len(errs) != 0 {
		t.Errorf("expected no errors for valid email, got %v", errs)
	}
	if errs := validateValue("ADMIN", "not-an-email", vs); len(errs) == 0 {
		t.Error("expected error for invalid email")
	}
}

func TestValidateValue_Pattern(t *testing.T) {
	vs := schema.VarSchema{Type: "string", Pattern: `^v\d+\.\d+\.\d+$`}
	if errs := validateValue("VERSION", "v1.2.3", vs); len(errs) != 0 {
		t.Errorf("expected no errors for matching pattern, got %v", errs)
	}
	if errs := validateValue("VERSION", "1.2.3", vs); len(errs) == 0 {
		t.Error("expected error for non-matching pattern")
	}
}

func TestValidateValue_AllowedValues(t *testing.T) {
	vs := schema.VarSchema{Type: "string", AllowedValues: []string{"dev", "staging", "prod"}}
	if errs := validateValue("ENV", "prod", vs); len(errs) != 0 {
		t.Errorf("expected no errors for allowed value, got %v", errs)
	}
	if errs := validateValue("ENV", "local", vs); len(errs) == 0 {
		t.Error("expected error for disallowed value")
	}
}

func TestValidateValue_InvalidPattern(t *testing.T) {
	vs := schema.VarSchema{Type: "string", Pattern: `[invalid`}
	if errs := validateValue("KEY", "value", vs); len(errs) == 0 {
		t.Error("expected error for invalid regex pattern")
	}
}
