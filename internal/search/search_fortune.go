package search

import (
	"sort"
	"sync"
)

type QueryTrend struct {
	Term string `json:"term"`
	Hits int64  `json:"hits"`
}

type SearchFortune struct {
	mu     sync.RWMutex
	trends map[string]int64
}

func NewSearchFortune() *SearchFortune {
	return &SearchFortune{
		trends: make(map[string]int64),
	}
}

// RecordQuery increments frequency counters for executed queries
func (sf *SearchFortune) RecordQuery(term string) {
	if len(term) < 2 {
		return
	}
	sf.mu.Lock()
	sf.trends[term]++
	sf.mu.Unlock()
}

// GetTopSearches returns the most popular search terms sorted by frequency
func (sf *SearchFortune) GetTopSearches(limit int) []QueryTrend {
	sf.mu.RLock()
	defer sf.mu.RUnlock()

	var list []QueryTrend
	for term, count := range sf.trends {
		list = append(list, QueryTrend{Term: term, Hits: count})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Hits > list[j].Hits
	})

	if len(list) > limit {
		return list[:limit]
	}
	return list
}
