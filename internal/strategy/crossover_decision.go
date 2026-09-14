package strategy

import (
	"sync"
	"time"
)

type ExecutionMode int

const (
	ModeLightweightHTTP ExecutionMode = iota
	ModeBrowserEmulated
	ModeFallbackProvider
)

type CrossoverDecision struct {
	mu            sync.RWMutex
	failureCount  int
	lastCrossover time.Time
	currentMode   ExecutionMode
}

func NewCrossoverDecision() *CrossoverDecision {
	return &CrossoverDecision{
		currentMode: ModeLightweightHTTP,
	}
}

// EvaluateCrossOver checks threshold limits and signals if execution mode must change
func (cd *CrossoverDecision) EvaluateCrossOver(lastErr error) ExecutionMode {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	if lastErr == nil {
		cd.failureCount = 0
		return cd.currentMode
	}

	cd.failureCount++
	if cd.failureCount >= 3 && cd.currentMode == ModeLightweightHTTP {
		cd.currentMode = ModeBrowserEmulated
		cd.lastCrossover = time.Now()
	} else if cd.failureCount >= 5 && cd.currentMode == ModeBrowserEmulated {
		cd.currentMode = ModeFallbackProvider
		cd.lastCrossover = time.Now()
	}
	return cd.currentMode
}

func (cd *CrossoverDecision) GetCurrentMode() ExecutionMode {
	cd.mu.RLock()
	defer cd.mu.RUnlock()
	return cd.currentMode
}
