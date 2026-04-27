package validator_test

import (
	"testing"

	"github.com/user/envlint/internal/validator"
)

func TestIssue_IsError(t *testing.T) {
	i := validator.Issue{Key: "FOO", Severity: "error", Message: "missing"}
	if !i.IsError() {
		t.Error("expected IsError() to return true")
	}
	if i.IsWarning() {
		t.Error("expected IsWarning() to return false")
	}
}

func TestIssue_IsWarning(t *testing.T) {
	i := validator.Issue{Key: "BAR", Severity: "warning", Message: "deprecated"}
	if !i.IsWarning() {
		t.Error("expected IsWarning() to return true")
	}
	if i.IsError() {
		t.Error("expected IsError() to return false")
	}
}

func TestHasErrors_WithErrors(t *testing.T) {
	issues := []validator.Issue{
		{Key: "A", Severity: "warning", Message: "w"},
		{Key: "B", Severity: "error", Message: "e"},
	}
	if !validator.HasErrors(issues) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_OnlyWarnings(t *testing.T) {
	issues := []validator.Issue{
		{Key: "A", Severity: "warning", Message: "w"},
	}
	if validator.HasErrors(issues) {
		t.Error("expected HasErrors to return false when only warnings")
	}
}

func TestHasErrors_Empty(t *testing.T) {
	if validator.HasErrors(nil) {
		t.Error("expected HasErrors to return false for empty slice")
	}
}
