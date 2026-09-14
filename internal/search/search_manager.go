package search

import (
	"context"
	"time"

	"audio-command-center/internal/cache"
)

type SearchManager struct {
	cache     *cache.CacheManager
	platforms []PlatformProvider
}

func NewSearchManager(c *cache.CacheManager) *SearchManager {
	return &SearchManager{
		cache:     c,
		platforms: []PlatformProvider{},
	}
}

// QueryFallback performs sequential provider search, returning the first successful result.
func (sm *SearchManager) QueryFallback(ctx context.Context, query string) ([]ProviderSearchResult, error) {
	// Check Cache First
	if cached, found := sm.cache.Get("search:" + query); found {
		return cached.([]ProviderSearchResult), nil
	}

	// Sequential Fallback Search
	for _, platform := range sm.platforms {
		res, err := platform.Search(ctx, query)
		if err == nil && len(res) > 0 {
			// Found results! Store in cache and return.
			sm.cache.Set("search:"+query, res, 30*time.Minute)
			return res, nil
		}
	}

	// No results from any provider
	return []ProviderSearchResult{}, nil
}

// AddPlatform allows adding new platforms to the manager.
func (sm *SearchManager) AddPlatform(p PlatformProvider) {
	sm.platforms = append(sm.platforms, p)
}
