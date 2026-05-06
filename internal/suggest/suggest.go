// Package suggest provides typo-correction hints for unknown environment
// variable names by comparing them against the known schema keys.
package suggest

import "strings"

// MaxDistance is the maximum Levenshtein distance considered a suggestion.
const MaxDistance = 3

// Closest returns up to limit schema keys that are "close" to name,
// ordered by ascending edit distance. If no key is within MaxDistance
// the returned slice is empty.
func Closest(name string, schemaKeys []string, limit int) []string {
	type candidate struct {
		key  string
		dist int
	}

	var candidates []candidate
	normName := strings.ToUpper(name)

	for _, k := range schemaKeys {
		d := levenshtein(normName, strings.ToUpper(k))
		if d <= MaxDistance {
			candidates = append(candidates, candidate{key: k, dist: d})
		}
	}

	// simple insertion sort — candidate list is typically tiny
	for i := 1; i < len(candidates); i++ {
		for j := i; j > 0 && candidates[j].dist < candidates[j-1].dist; j-- {
			candidates[j], candidates[j-1] = candidates[j-1], candidates[j]
		}
	}

	if limit <= 0 || limit > len(candidates) {
		limit = len(candidates)
	}

	out := make([]string, limit)
	for i := 0; i < limit; i++ {
		out[i] = candidates[i].key
	}
	return out
}

// levenshtein computes the edit distance between two strings.
func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)

	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)

	for j := 0; j <= lb; j++ {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			curr[j] = min3(prev[j]+1, curr[j-1]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
