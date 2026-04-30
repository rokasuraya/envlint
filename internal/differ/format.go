package differ

import (
	"fmt"
	"io"
	"sort"
)

// WriteText writes a human-readable diff report to w.
// Sensitive values should be redacted by the caller before passing maps.
func WriteText(w io.Writer, changes []Change) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	sorted := make([]Change, len(changes))
	copy(sorted, changes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, c := range sorted {
		switch c.Kind {
		case Added:
			fmt.Fprintf(w, "+ %-30s = %s\n", c.Key, c.NewValue)
		case Removed:
			fmt.Fprintf(w, "- %-30s = %s\n", c.Key, c.OldValue)
		case Changed:
			fmt.Fprintf(w, "~ %-30s   %s -> %s\n", c.Key, c.OldValue, c.NewValue)
		}
	}
}

// Summary returns a one-line summary of the diff result.
func Summary(changes []Change) string {
	var added, removed, changed int
	for _, c := range changes {
		switch c.Kind {
		case Added:
			added++
		case Removed:
			removed++
		case Changed:
			changed++
		}
	}
	return fmt.Sprintf("%d added, %d removed, %d changed", added, removed, changed)
}
