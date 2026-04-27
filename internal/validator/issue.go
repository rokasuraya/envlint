package validator

// Issue represents a single validation problem found in the .env file.
type Issue struct {
	// Key is the environment variable name that triggered the issue.
	Key string
	// Severity is either "error" or "warning".
	Severity string
	// Message is a human-readable description of the problem.
	Message string
}

// IsError returns true when the issue is considered a hard error.
func (i Issue) IsError() bool {
	return i.Severity == "error"
}

// IsWarning returns true when the issue is a non-fatal warning.
func (i Issue) IsWarning() bool {
	return i.Severity == "warning"
}

// HasErrors returns true if any issue in the slice is an error.
func HasErrors(issues []Issue) bool {
	for _, iss := range issues {
		if iss.IsError() {
			return true
		}
	}
	return false
}
