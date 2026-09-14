package search

import (
	"strings"
)

type QuestStrategy int

const (
	QuestExact QuestStrategy = iota
	QuestFuzzyArtist
	QuestGenreFallback
)

type SearchQuest struct {
	Query          string
	NormalizedTerm string
	Strategy       QuestStrategy
	AttemptCount   int
}

func NewSearchQuest(rawQuery string) *SearchQuest {
	return &SearchQuest{
		Query:          rawQuery,
		NormalizedTerm: strings.TrimSpace(strings.ToLower(rawQuery)),
		Strategy:       QuestExact,
	}
}

// NextMutation generates fallback search terms if initial exact matches yield low results
func (sq *SearchQuest) NextMutation() (string, bool) {
	sq.AttemptCount++
	switch sq.Strategy {
	case QuestExact:
		sq.Strategy = QuestFuzzyArtist
		// Strip parenthetical descriptors like "(Original Mix)" or "[HQ]" for fuzzy retry
		cleanTerm := stripModifiers(sq.Query)
		if cleanTerm != sq.Query {
			return cleanTerm, true
		}
		fallthrough
	case QuestFuzzyArtist:
		sq.Strategy = QuestGenreFallback
		parts := strings.Split(sq.NormalizedTerm, "-")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[0]), true // Retry with artist name only
		}
		return "", false
	default:
		return "", false
	}
}

func stripModifiers(term string) string {
	r := strings.NewReplacer("(", "", ")", "", "[", "", "]", "")
	return strings.TrimSpace(r.Replace(term))
}
