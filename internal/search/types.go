package search

import (
	"context"
	"time"
)

type DraftStatus int

const (
	DraftPending DraftStatus = iota
	DraftValidated
	DraftRejected
)

// SearchDraft represents a raw search result from a provider.
type SearchDraft struct {
	ID          string      `json:"id"`
	RawTitle    string      `json:"raw_title"`
	RawArtist   string      `json:"raw_artist"`
	Source      string      `json:"source"`
	SourceURL   string      `json:"source_url"`
	DurationSec int         `json:"duration_sec"`
	Status      DraftStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
}

// ProviderSearchResult represents a result from a search provider.
type ProviderSearchResult struct {
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Source  string `json:"source"`
	URL     string `json:"url"`
	Bitrate string `json:"bitrate"`
}

// PlatformProvider defines the interface for search providers.
type PlatformProvider interface {
	Search(ctx context.Context, query string) ([]ProviderSearchResult, error)
}

// SearchResult represents a cleaned and organized search result.
type SearchResult struct {
	Title    string
	Artist   string
	Sources  []string
	URL      string
	Duration int
}
