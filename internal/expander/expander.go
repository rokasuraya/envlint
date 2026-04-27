// Package expander resolves variable interpolation within .env values.
// It supports ${VAR} and $VAR syntax, expanding references from the
// already-parsed environment map or the process environment.
package expander

import (
	"fmt"
	"os"
	"strings"
)

// Expander resolves variable references in env values.
type Expander struct {
	// FallbackToOS controls whether unresolved vars fall back to os.Getenv.
	FallbackToOS bool
}

// New returns an Expander with OS fallback enabled by default.
func New() *Expander {
	return &Expander{FallbackToOS: true}
}

// Expand resolves all variable references in the provided env map in-place.
// References are resolved in insertion order; forward references to keys not
// yet expanded will fall back to the OS environment if FallbackToOS is set.
func (e *Expander) Expand(env map[string]string) error {
	for k, v := range env {
		expanded, err := e.expandValue(v, env)
		if err != nil {
			return fmt.Errorf("expanding %q: %w", k, err)
		}
		env[k] = expanded
	}
	return nil
}

// ExpandValue resolves variable references in a single value string.
func (e *Expander) ExpandValue(value string, env map[string]string) (string, error) {
	return e.expandValue(value, env)
}

func (e *Expander) expandValue(value string, env map[string]string) (string, error) {
	var sb strings.Builder
	i := 0
	for i < len(value) {
		if value[i] != '$' {
			sb.WriteByte(value[i])
			i++
			continue
		}
		// escape: $$
		if i+1 < len(value) && value[i+1] == '$' {
			sb.WriteByte('$')
			i += 2
			continue
		}
		var name string
		var advance int
		if i+1 < len(value) && value[i+1] == '{' {
			end := strings.Index(value[i+2:], "}")
			if end == -1 {
				return "", fmt.Errorf("unclosed ${ in value")
			}
			name = value[i+2 : i+2+end]
			advance = 3 + len(name)
		} else {
			j := i + 1
			for j < len(value) && isVarChar(value[j]) {
				j++
			}
			name = value[i+1 : j]
			advance = 1 + len(name)
		}
		if name == "" {
			sb.WriteByte('$')
			i++
			continue
		}
		sb.WriteString(e.resolve(name, env))
		i += advance
	}
	return sb.String(), nil
}

func (e *Expander) resolve(name string, env map[string]string) string {
	if v, ok := env[name]; ok {
		return v
	}
	if e.FallbackToOS {
		return os.Getenv(name)
	}
	return ""
}

func isVarChar(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') || c == '_'
}
