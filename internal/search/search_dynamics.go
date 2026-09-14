package search

import (
	"strings"
)

type SearchDynamics struct{}

func NewSearchDynamics() *SearchDynamics {
	return &SearchDynamics{}
}

// CorrectAndExpand fixes common typos and expands short query terms
func (sd *SearchDynamics) CorrectAndExpand(input string) string {
	cleaned := strings.TrimSpace(strings.ToLower(input))
	if len(cleaned) < 3 {
		return cleaned // Retain short acronyms or single-word titles
	}

	// Dynamic typo map for common audio search misspellings
	replacements := map[string]string{
		"remixx":    "remix",
		"featt":     "feat",
		"soudtrack": "soundtrack",
		"albim":     "album",
	}

	words := strings.Fields(cleaned)
	for i, w := range words {
		if fix, ok := replacements[w]; ok {
			words[i] = fix
		}
	}
	return strings.Join(words, " ")
}

// LevenshteinDistance calculates similarity score to evaluate query candidate matches
func LevenshteinDistance(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	n, m := len(r1), len(r2)
	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
		dp[i][0] = i
	}
	for j := 0; j <= m; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}
			dp[i][j] = min(dp[i-1][j]+1, min(dp[i][j-1]+1, dp[i-1][j-1]+cost))
		}
	}
	return dp[n][m]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
