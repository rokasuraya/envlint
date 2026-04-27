// Package parser reads .env files into a plain key/value map.
package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envlint/internal/expander"
)

// ParseOptions controls optional behaviour of ParseFile.
type ParseOptions struct {
	// Expand resolves ${VAR} / $VAR references after parsing.
	Expand bool
	// FallbackToOS allows unresolved vars to fall back to os.Getenv.
	FallbackToOS bool
}

// ParseFile reads the .env file at path and returns a map of key→value pairs.
// Comments (# …) and blank lines are ignored. The export keyword is stripped.
// If opts.Expand is true, variable interpolation is performed after parsing.
func ParseFile(path string, opts *ParseOptions) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, err := parseLine(line)
		if err != nil {
			continue // skip malformed lines
		}
		env[k] = v
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading %q: %w", path, err)
	}

	if opts != nil && opts.Expand {
		ex := expander.New()
		ex.FallbackToOS = opts.FallbackToOS
		if err := ex.Expand(env); err != nil {
			return nil, fmt.Errorf("expanding variables: %w", err)
		}
	}
	return env, nil
}

// parseLine splits a single KEY=VALUE line into its components.
func parseLine(line string) (string, string, error) {
	line = strings.TrimPrefix(line, "export ")
	idx := strings.IndexByte(line, '=')
	if idx == -1 {
		return "", "", fmt.Errorf("no '=' in line: %q", line)
	}
	key := strings.TrimSpace(line[:idx])
	val := strings.TrimSpace(line[idx+1:])
	if key == "" {
		return "", "", fmt.Errorf("empty key in line: %q", line)
	}
	// Strip inline comment (unquoted)
	if !isQuoted(val) {
		if ci := strings.Index(val, " #"); ci != -1 {
			val = strings.TrimSpace(val[:ci])
		}
	}
	val = stripQuotes(val)
	return key, val, nil
}

func isQuoted(s string) bool {
	return (strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) ||
		(strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'"))
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
