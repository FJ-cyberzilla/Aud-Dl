package strategy

import (
	"sort"
	"sync"
	"time"
)

type ProviderScore struct {
	Name        string
	Latency     time.Duration
	SuccessRate float64
	RankScore   float64
}

type HybridDecisive struct {
	mu     sync.RWMutex
	scores map[string]*ProviderScore
}

func NewHybridDecisive() *HybridDecisive {
	return &HybridDecisive{
		scores: make(map[string]*ProviderScore),
	}
}

// RecordMetric updates telemetry for a given provider
func (hd *HybridDecisive) RecordMetric(name string, latency time.Duration, success bool) {
	hd.mu.Lock()
	defer hd.mu.Unlock()

	score, exists := hd.scores[name]
	if !exists {
		score = &ProviderScore{Name: name, SuccessRate: 1.0}
		hd.scores[name] = score
	}

	// Exponential moving average for latency
	score.Latency = (score.Latency*7 + latency*3) / 10
	if success {
		score.SuccessRate = (score.SuccessRate * 0.9) + 0.1
	} else {
		score.SuccessRate = (score.SuccessRate * 0.9)
	}

	// Calculate weighted composite score
	score.RankScore = score.SuccessRate*100.0 - (float64(score.Latency.Milliseconds()) * 0.05)
}

// SelectOptimalRoute returns providers ordered by best execution stability
func (hd *HybridDecisive) SelectOptimalRoute() []string {
	hd.mu.RLock()
	defer hd.mu.RUnlock()

	var list []*ProviderScore
	for _, s := range hd.scores {
		list = append(list, s)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].RankScore > list[j].RankScore
	})

	routes := make([]string, len(list))
	for i, s := range list {
		routes[i] = s.Name
	}
	return routes
}
