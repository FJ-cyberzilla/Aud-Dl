package search

import (
	"sync"
)

type DraftBuffer struct {
	mu     sync.RWMutex
	drafts map[string]*SearchDraft
}

func NewDraftBuffer() *DraftBuffer {
	return &DraftBuffer{
		drafts: make(map[string]*SearchDraft),
	}
}

func (db *DraftBuffer) AddDraft(draft *SearchDraft) {
	db.mu.Lock()
	db.drafts[draft.ID] = draft
	db.mu.Unlock()
}

func (db *DraftBuffer) GetPending() []*SearchDraft {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var pending []*SearchDraft
	for _, d := range db.drafts {
		if d.Status == DraftPending {
			pending = append(pending, d)
		}
	}
	return pending
}
