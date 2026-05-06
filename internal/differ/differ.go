// Package differ compares two parsed .env files and reports
// added, removed, and changed variables between them.
package differ

// ChangeKind describes the type of change for a variable.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
)

// Change represents a single variable difference between two env maps.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Diff computes the differences between a baseline env map and a target env map.
// Keys present only in baseline are Removed; keys only in target are Added;
// keys in both with differing values are Changed.
func Diff(baseline, target map[string]string) []Change {
	var changes []Change

	for k, oldVal := range baseline {
		newVal, ok := target[k]
		if !ok {
			changes = append(changes, Change{
				Key:      k,
				Kind:     Removed,
				OldValue: oldVal,
			})
			continue
		}
		if newVal != oldVal {
			changes = append(changes, Change{
				Key:      k,
				Kind:     Changed,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	for k, newVal := range target {
		if _, exists := baseline[k]; !exists {
			changes = append(changes, Change{
				Key:      k,
				Kind:     Added,
				NewValue: newVal,
			})
		}
	}

	return changes
}

// HasChanges returns true if any differences exist.
func HasChanges(changes []Change) bool {
	return len(changes) > 0
}

// Filter returns only the changes that match the given kind.
// This is useful when callers want to inspect only additions,
// removals, or modifications independently.
func Filter(changes []Change, kind ChangeKind) []Change {
	var filtered []Change
	for _, c := range changes {
		if c.Kind == kind {
			filtered = append(filtered, c)
		}
	}
	return filtered
}
