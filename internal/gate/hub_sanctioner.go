package gate

import (
	"fmt"
	"sync"
	"time"
)

type SanctionLevel int

const (
	SanctionNone SanctionLevel = iota
	SanctionWarning
	SanctionTempBan
	SanctionPermanentBan
)

type SanctionRecord struct {
	Domain       string
	Level        SanctionLevel
	OffenseCount int
	ExpiresAt    time.Time
	Reason       string
}

type HubSanctioner struct {
	mu          sync.RWMutex
	sanctions   map[string]*SanctionRecord
	maxOffenses int
	banDuration time.Duration
}

func NewHubSanctioner(maxOffenses int, banDuration time.Duration) *HubSanctioner {
	if maxOffenses <= 0 {
		maxOffenses = 3
	}
	if banDuration <= 0 {
		banDuration = 15 * time.Minute
	}
	return &HubSanctioner{
		sanctions:   make(map[string]*SanctionRecord),
		maxOffenses: maxOffenses,
		banDuration: banDuration,
	}
}

// InspectSanction checks if a domain or provider is currently under an active sanction lock
func (hs *HubSanctioner) InspectSanction(domain string) error {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	record, exists := hs.sanctions[domain]
	if !exists {
		return nil
	}

	// Check if a temporary sanction has expired
	if record.Level == SanctionTempBan && time.Now().After(record.ExpiresAt) {
		delete(hs.sanctions, domain)
		return nil
	}

	if record.Level >= SanctionTempBan {
		return fmt.Errorf("domain [%s] is sanctioned due to %s (Expires in %v)",
			domain, record.Reason, time.Until(record.ExpiresAt).Round(time.Second))
	}

	return nil
}

// RecordViolation logs an anti-bot block or security trigger and escalates penalties if needed
func (hs *HubSanctioner) RecordViolation(domain, reason string) SanctionLevel {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	record, exists := hs.sanctions[domain]
	if !exists {
		record = &SanctionRecord{
			Domain: domain,
			Level:  SanctionWarning,
		}
		hs.sanctions[domain] = record
	}

	record.OffenseCount++
	record.Reason = reason

	if record.OffenseCount >= hs.maxOffenses {
		record.Level = SanctionTempBan
		record.ExpiresAt = time.Now().Add(hs.banDuration)
	} else {
		record.Level = SanctionWarning
	}

	return record.Level
}

// ClearSanction manually resets penalties for a specific target domain
func (hs *HubSanctioner) ClearSanction(domain string) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	delete(hs.sanctions, domain)
}

// GetSanctionedCount returns the number of currently sanctioned domains
func (hs *HubSanctioner) GetSanctionedCount() int {
	hs.mu.RLock()
	defer hs.mu.RUnlock()
	count := 0
	for _, record := range hs.sanctions {
		if record.Level >= SanctionTempBan {
			count++
		}
	}
	return count
}
