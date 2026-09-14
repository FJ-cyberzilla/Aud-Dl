package strategy

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type SessionDriftState struct {
	LastDrift    time.Time
	RequestCount int64
	ActiveJa4    string
}

type AntibotDrifter struct {
	mu            sync.RWMutex
	state         SessionDriftState
	driftInterval time.Duration
	maxRequests   int64
}

func NewAntibotDrifter(interval time.Duration, maxReqs int64) *AntibotDrifter {
	return &AntibotDrifter{
		driftInterval: interval,
		maxRequests:   maxReqs,
		state: SessionDriftState{
			LastDrift: time.Now(),
		},
	}
}

// ShouldDrift checks if session parameters need a proactive mutation
func (ad *AntibotDrifter) ShouldDrift() bool {
	ad.mu.RLock()
	defer ad.mu.RUnlock()
	if time.Since(ad.state.LastDrift) > ad.driftInterval {
		return true
	}
	if ad.state.RequestCount >= ad.maxRequests {
		return true
	}
	return false
}

// ApplyDrift mutates session transport context to mimic a natural client transition
func (ad *AntibotDrifter) ApplyDrift(req *http.Request, client *http.Client) {
	ad.mu.Lock()
	defer ad.mu.Unlock()

	// Reset counter and refresh drift timestamp
	ad.state.RequestCount = 0
	ad.state.LastDrift = time.Now()

	// Flush active connection pool in client transport to reset TLS Session ID
	if tr, ok := client.Transport.(*http.Transport); ok {
		tr.CloseIdleConnections()
	}

	// Micro-shift user agent client hint versions
	req.Header.Set("Sec-Ch-Ua-Platform-Version", getRandomVersion())
}

func (ad *AntibotDrifter) Increment() {
	ad.mu.Lock()
	ad.state.RequestCount++
	ad.mu.Unlock()
}

func getRandomVersion() string {
	versions := []string{`"10.0.0"`, `"13.0.0"`, `"15.0.0"`}
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(versions))))
	if err != nil {
		return versions[0]
	}
	return versions[nBig.Int64()]
}
