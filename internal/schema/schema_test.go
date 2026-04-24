package schema

import (
	"testing"
)

func TestVarSchema_ValidateInt(t *testing.T) {
	vs := VarSchema{Required: true, Type: TypeInt}
	if err := vs.Validate("PORT", "8080"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := vs.Validate("PORT", "abc"); err == nil {
		t.Error("expected error for non-int value")
	}
}

func TestVarSchema_ValidateBool(t *testing.T) {
	vs := VarSchema{Type: TypeBool}
	for _, v := range []string{"true", "false", "1", "0", "yes", "no"} {
		if err := vs.Validate("FLAG", v); err != nil {
			t.Errorf("unexpected error for %q: %v", v, err)
		}
	}
	if err := vs.Validate("FLAG", "maybe"); err == nil {
		t.Error("expected error for invalid bool")
	}
}

func TestVarSchema_ValidateURL(t *testing.T) {
	vs := VarSchema{Type: TypeURL}
	if err := vs.Validate("ENDPOINT", "https://example.com"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := vs.Validate("ENDPOINT", "not-a-url"); err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestVarSchema_ValidateEmail(t *testing.T) {
	vs := VarSchema{Type: TypeEmail}
	if err := vs.Validate("EMAIL", "user@example.com"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := vs.Validate("EMAIL", "notanemail"); err == nil {
		t.Error("expected error for invalid email")
	}
}

func TestVarSchema_ValidatePattern(t *testing.T) {
	vs := VarSchema{Type: TypeString, Pattern: `^v\d+\.\d+\.\d+$`}
	if err := vs.Validate("VERSION", "v1.2.3"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := vs.Validate("VERSION", "1.2.3"); err == nil {
		t.Error("expected error for pattern mismatch")
	}
}

func TestVarSchema_RequiredEmpty(t *testing.T) {
	vs := VarSchema{Required: true, Type: TypeString}
	if err := vs.Validate("KEY", ""); err == nil {
		t.Error("expected error for required empty variable")
	}
}

func TestVarSchema_AllowEmpty(t *testing.T) {
	vs := VarSchema{Required: true, Type: TypeInt, AllowEmpty: true}
	if err := vs.Validate("KEY", ""); err != nil {
		t.Errorf("unexpected error when AllowEmpty=true: %v", err)
	}
}

func TestVarSchema_ValidateEnum(t *testing.T) {
	vs := VarSchema{Type: TypeString, Enum: []string{"debug", "info", "warn", "error"}}
	for _, v := range []string{"debug", "info", "warn", "error"} {
		if err := vs.Validate("LOG_LEVEL", v); err != nil {
			t.Errorf("unexpected error for valid enum value %q: %v", v, err)
		}
	}
	if err := vs.Validate("LOG_LEVEL", "trace"); err == nil {
		t.Error("expected error for value not in enum")
	}
}
