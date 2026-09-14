package search

import (
	"sort"
	"strings"
)

type SearchOrganizer struct{}

func NewSearchOrganizer() *SearchOrganizer {
	return &SearchOrganizer{}
}

// Organize cleans, deduplicates across sources, and ranks drafts by quality score
func (so *SearchOrganizer) Organize(drafts []*SearchDraft) []SearchResult {
	seen := make(map[string]*SearchResult)

	for _, d := range drafts {
		key := normalizeKey(d.RawArtist, d.RawTitle)

		if existing, found := seen[key]; found {
			// Merge provider sources if duplicate track is found on multiple platforms
			existing.Sources = append(existing.Sources, d.Source)
			continue
		}

		seen[key] = &SearchResult{
			Title:    strings.Title(d.RawTitle),
			Artist:   strings.Title(d.RawArtist),
			Sources:  []string{d.Source},
			URL:      d.SourceURL,
			Duration: d.DurationSec,
		}
	}

	results := make([]SearchResult, 0, len(seen))
	for _, res := range seen {
		results = append(results, *res)
	}

	// Sort results by popularity or completeness
	sort.Slice(results, func(i, j int) bool {
		return len(results[i].Sources) > len(results[j].Sources)
	})

	return results
}

func normalizeKey(artist, title string) string {
	clean := strings.ToLower(artist + "-" + title)
	clean = strings.ReplaceAll(clean, " ", "")
	return clean
}
