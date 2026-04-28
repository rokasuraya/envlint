// Package redactor masks sensitive values in .env output to prevent
// accidental exposure of secrets in logs or reports.
package redactor

import "strings"

// DefaultSensitiveKeys contains common key substrings that indicate a
// variable holds a secret and should have its value redacted.
var DefaultSensitiveKeys = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

const redactedPlaceholder = "[REDACTED]"

// Redactor masks values whose keys match sensitive patterns.
type Redactor struct {
	sensitiveKeys []string
}

// New returns a Redactor using the provided sensitive key substrings.
// Pass nil to use DefaultSensitiveKeys.
func New(sensitiveKeys []string) *Redactor {
	if sensitiveKeys == nil {
		sensitiveKeys = DefaultSensitiveKeys
	}
	return &Redactor{sensitiveKeys: sensitiveKeys}
}

// IsSensitive reports whether key matches any of the configured sensitive
// key substrings (case-insensitive).
func (r *Redactor) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, s := range r.sensitiveKeys {
		if strings.Contains(upper, strings.ToUpper(s)) {
			return true
		}
	}
	return false
}

// Redact returns the original value unchanged if the key is not sensitive,
// or the redacted placeholder if it is.
func (r *Redactor) Redact(key, value string) string {
	if r.IsSensitive(key) {
		return redactedPlaceholder
	}
	return value
}

// RedactMap returns a copy of env with sensitive values replaced.
func (r *Redactor) RedactMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = r.Redact(k, v)
	}
	return out
}
