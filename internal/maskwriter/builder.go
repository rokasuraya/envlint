package maskwriter

import (
	"io"

	"github.com/yourorg/envlint/internal/redactor"
)

// FromEnv constructs a Writer that masks every value in env whose key is
// considered sensitive according to the default redactor heuristics.
// Non-sensitive values are left unmasked.
//
// This is the recommended entry-point when integrating maskwriter with the
// envlint linting pipeline.
func FromEnv(w io.Writer, env map[string]string) *Writer {
	r := redactor.New(nil) // nil → use default sensitive-key patterns
	redacted := r.RedactMap(env)

	secrets := make([]string, 0)
	for key, original := range env {
		masked := redacted[key]
		// If the redactor changed the value it is sensitive — collect the
		// original plaintext so we can mask it in output.
		if masked != original {
			secrets = append(secrets, original)
		}
	}
	return New(w, secrets)
}
