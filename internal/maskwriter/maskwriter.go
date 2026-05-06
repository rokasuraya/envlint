// Package maskwriter provides an io.Writer wrapper that redacts sensitive
// values before they are written to an underlying writer (e.g. log output).
package maskwriter

import (
	"io"
	"strings"
)

const redacted = "[REDACTED]"

// Writer wraps an io.Writer and replaces any occurrence of a sensitive
// literal value with the placeholder [REDACTED].
type Writer struct {
	w        io.Writer
	secrets  []string
}

// New returns a Writer that masks every string in secrets before forwarding
// bytes to w. Empty secrets are silently ignored.
func New(w io.Writer, secrets []string) *Writer {
	filtered := make([]string, 0, len(secrets))
	for _, s := range secrets {
		if s != "" {
			filtered = append(filtered, s)
		}
	}
	return &Writer{w: w, secrets: filtered}
}

// Write masks all registered secrets within p and writes the result to the
// underlying writer.
func (mw *Writer) Write(p []byte) (int, error) {
	line := string(p)
	for _, secret := range mw.secrets {
		line = strings.ReplaceAll(line, secret, redacted)
	}
	_, err := io.WriteString(mw.w, line)
	if err != nil {
		return 0, err
	}
	// Report the original length so callers do not see a short-write error.
	return len(p), nil
}

// AddSecret registers an additional secret to be masked. Empty values are
// ignored.
func (mw *Writer) AddSecret(secret string) {
	if secret != "" {
		mw.secrets = append(mw.secrets, secret)
	}
}
