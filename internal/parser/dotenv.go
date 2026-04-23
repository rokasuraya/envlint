package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair parsed from a .env file.
type Entry struct {
	Key   string
	Value string
	Line  int
}

// ParseFile reads a .env file and returns a slice of Entries.
// Lines beginning with '#' are treated as comments and skipped.
// Blank lines are also skipped.
func ParseFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: open %q: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments.
		if raw == "" || strings.HasPrefix(raw, "#") {
			continue
		}

		key, value, err := parseLine(raw)
		if err != nil {
			return nil, fmt.Errorf("parser: %q line %d: %w", path, lineNum, err)
		}

		entries = append(entries, Entry{Key: key, Value: value, Line: lineNum})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: scanning %q: %w", path, err)
	}

	return entries, nil
}

// parseLine splits a raw line into a key and value.
func parseLine(line string) (string, string, error) {
	// Support optional "export " prefix.
	line = strings.TrimPrefix(line, "export ")

	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return "", "", fmt.Errorf("missing '=' in line %q", line)
	}

	key := strings.TrimSpace(line[:idx])
	if key == "" {
		return "", "", fmt.Errorf("empty key in line %q", line)
	}

	value := strings.TrimSpace(line[idx+1:])
	value = stripQuotes(value)

	return key, value, nil
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
