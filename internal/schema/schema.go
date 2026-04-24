package schema

import (
	"fmt"
	"regexp"
)

// VarType represents the expected type of an environment variable.
type VarType string

const (
	TypeString VarType = "string"
	TypeInt    VarType = "int"
	TypeBool   VarType = "bool"
	TypeURL    VarType = "url"
	TypeEmail  VarType = "email"
)

// VarSchema defines the rules for a single environment variable.
type VarSchema struct {
	Required    bool
	Type        VarType
	Pattern     string
	AllowEmpty  bool
	Description string
}

// Schema holds the full set of variable definitions.
type Schema struct {
	Vars map[string]VarSchema
}

// Validate checks a single value against the given VarSchema.
// Returns an error if the value does not conform.
func (vs VarSchema) Validate(key, value string) error {
	if value == "" && !vs.AllowEmpty {
		if vs.Required {
			return fmt.Errorf("%s: required variable is empty", key)
		}
		return nil
	}

	switch vs.Type {
	case TypeInt:
		if !regexp.MustCompile(`^-?\d+$`).MatchString(value) {
			return fmt.Errorf("%s: expected int, got %q", key, value)
		}
	case TypeBool:
		if !regexp.MustCompile(`^(true|false|1|0|yes|no)$`).MatchString(value) {
			return fmt.Errorf("%s: expected bool, got %q", key, value)
		}
	case TypeURL:
		if !regexp.MustCompile(`^https?://[^\s]+$`).MatchString(value) {
			return fmt.Errorf("%s: expected URL, got %q", key, value)
		}
	case TypeEmail:
		if !regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`).MatchString(value) {
			return fmt.Errorf("%s: expected email, got %q", key, value)
		}
	}

	if vs.Pattern != "" {
		re, err := regexp.Compile(vs.Pattern)
		if err != nil {
			return fmt.Errorf("%s: invalid pattern %q: %w", key, vs.Pattern, err)
		}
		if !re.MatchString(value) {
			return fmt.Errorf("%s: value %q does not match pattern %q", key, value, vs.Pattern)
		}
	}

	return nil
}
