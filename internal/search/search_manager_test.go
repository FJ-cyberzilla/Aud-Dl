package search

import (
	"context"
	"testing"
	"time"

	"audio-command-center/internal/cache"
)

type MockProvider struct {
	results []ProviderSearchResult
	err     error
}

func (m *MockProvider) Search(ctx context.Context, query string) ([]ProviderSearchResult, error) {
	return m.results, m.err
}

func TestQueryFallback(t *testing.T) {
	cm := cache.NewCacheManager(time.Minute)
	sm := NewSearchManager(cm)

	p1 := &MockProvider{results: []ProviderSearchResult{}, err: nil}
	p2 := &MockProvider{results: []ProviderSearchResult{{Title: "Found"}}, err: nil}

	sm.AddPlatform(p1)
	sm.AddPlatform(p2)

	res, err := sm.QueryFallback(context.Background(), "test query")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(res) != 1 || res[0].Title != "Found" {
		t.Errorf("expected result from second provider, got %v", res)
	}
}
