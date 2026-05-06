package suggest

import (
	"fmt"
	"io"
	"strings"
)

// WriteHints writes fuzzy-match suggestions for unknown keys to w.
// For each unknown key, it looks up the closest matches from knownKeys
// and prints them as hints. If no suggestions are found, it prints a
// generic "no similar keys found" message.
func WriteHints(w io.Writer, unknownKeys []string, knownKeys []string, limit int) {
	if len(unknownKeys) == 0 {
		return
	}

	fmt.Fprintln(w, "Suggestions for unrecognised keys:")
	for _, key := range unknownKeys {
		matches := Closest(key, knownKeys, limit)
		if len(matches) == 0 {
			fmt.Fprintf(w, "  %s — no similar keys found in schema\n", key)
			continue
		}
		fmt.Fprintf(w, "  %s — did you mean: %s?\n", key, strings.Join(matches, ", "))
	}
}

// UnknownKeys returns the keys present in env that are not defined in the schema.
func UnknownKeys(env map[string]string, schemaKeys []string) []string {
	known := make(map[string]struct{}, len(schemaKeys))
	for _, k := range schemaKeys {
		known[k] = struct{}{}
	}

	var unknown []string
	for k := range env {
		if _, ok := known[k]; !ok {
			unknown = append(unknown, k)
		}
	}
	return unknown
}
