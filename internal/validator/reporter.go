package validator

import (
	"fmt"
	"io"
	"strings"
)

// PrintReport writes a human-readable summary of the report to w.
func PrintReport(w io.Writer, report *Report) {
	if !report.HasErrors() {
		fmt.Fprintln(w, "✓ All variables are valid.")
		return
	}

	fmt.Fprintln(w, "✗ Validation failed:")
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, result := range report.Results {
		if !result.Missing && len(result.Errors) == 0 {
			continue
		}
		for _, errMsg := range result.Errors {
			fmt.Fprintf(w, "  [ERROR] %s\n", errMsg)
		}
	}

	fmt.Fprintln(w, strings.Repeat("-", 40))
	total := countErrors(report)
	fmt.Fprintf(w, "Total issues: %d\n", total)
}

func countErrors(report *Report) int {
	count := 0
	for _, r := range report.Results {
		if r.Missing {
			count++
			continue
		}
		count += len(r.Errors)
	}
	return count
}
