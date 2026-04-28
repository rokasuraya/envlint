// Package redactor provides utilities for detecting and masking sensitive
// environment variable values before they are included in reports, logs, or
// any other output that could be exposed to unintended audiences.
//
// A variable is considered sensitive when its key contains one of a
// configurable set of substrings such as "PASSWORD", "SECRET", or "TOKEN".
// Matching is case-insensitive so both DB_PASSWORD and db_password are
// treated the same way.
//
// Usage:
//
//	r := redactor.New(nil) // uses DefaultSensitiveKeys
//	safe := r.RedactMap(parsedEnv)
//
// Custom sensitive substrings can be supplied:
//
//	r := redactor.New([]string{"INTERNAL", "PRIVATE"})
package redactor
