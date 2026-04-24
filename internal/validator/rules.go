package validator

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/user/envlint/internal/schema"
)

// ValidationError represents a single validation failure for a variable.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Key, e.Message)
}

// validateValue checks a raw string value against the rules defined in VarSchema.
// It returns a slice of ValidationErrors (empty means valid).
func validateValue(key, value string, vs schema.VarSchema) []ValidationError {
	var errs []ValidationError

	switch strings.ToLower(vs.Type) {
	case "int", "integer":
		if _, err := strconv.Atoi(value); err != nil {
			errs = append(errs, ValidationError{Key: key, Message: fmt.Sprintf("expected integer, got %q", value)})
		}
	case "bool", "boolean":
		lower := strings.ToLower(value)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			errs = append(errs, ValidationError{Key: key, Message: fmt.Sprintf("expected boolean, got %q", value)})
		}
	case "url":
		u, err := url.ParseRequestURI(value)
		if err != nil || u.Scheme == "" || u.Host == "" {
			errs = append(errs, ValidationError{Key: key, Message: fmt.Sprintf("expected valid URL, got %q", value)})
		}
	case "email":
		if !regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`).MatchString(value) {
			errs = append(errs, ValidationError{Key: key, Message: fmt.Sprintf("expected valid email, got %q", value)})
		}
	}

	if vs.Pattern != "" {
		re, err := regexp.Compile(vs.Pattern)
		if err != nil {
			errs = append(errs, ValidationError{Key: key, Message: fmt.Sprintf("invalid pattern %q: %v", vs.Pattern, err)})
		} else if !re.MatchString(value) {
			errs = append(errs, ValidationError{Key: key, Message: fmt.Sprintf("value %q does not match pattern %q", value, vs.Pattern)})
		}
	}

	if len(vs.AllowedValues) > 0 {
		found := false
		for _, allowed := range vs.AllowedValues {
			if value == allowed {
				found = true
				break
			}
		}
		if !found {
			errs = append(errs, ValidationError{Key: key, Message: fmt.Sprintf("value %q not in allowed values %v", value, vs.AllowedValues)})
		}
	}

	return errs
}
