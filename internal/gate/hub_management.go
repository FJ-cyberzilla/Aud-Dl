package gate

import (
	"context"
	"fmt"
	"sync"
)

type RouteStatus int

const (
	RouteHealthy RouteStatus = iota
	RouteDegraded
	RouteBlocked
)

type HubManagement struct {
	mu             sync.RWMutex
	activeGateways map[string]RouteStatus
	egressTLS      *EgressTLSClient
	maxConnections int
	activeCount    int
}

func NewHubManagement(egress *EgressTLSClient, maxConns int) *HubManagement {
	if maxConns <= 0 {
		maxConns = 50
	}
	return &HubManagement{
		activeGateways: make(map[string]RouteStatus),
		egressTLS:      egress,
		maxConnections: maxConns,
	}
}

// AcquireConnection acquires an egress slot under strict rate/concurrency limits
func (hm *HubManagement) AcquireConnection(ctx context.Context) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	if hm.activeCount >= hm.maxConnections {
		return fmt.Errorf("gate network pool exhausted: %d/%d active connections", hm.activeCount, hm.maxConnections)
	}
	hm.activeCount++
	return nil
}

// ReleaseConnection releases an active transport slot
func (hm *HubManagement) ReleaseConnection() {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	if hm.activeCount > 0 {
		hm.activeCount--
	}
}

// RecordRouteStatus updates health state for a target CDN/Domain
func (hm *HubManagement) RecordRouteStatus(domain string, status RouteStatus) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.activeGateways[domain] = status
}

func (hm *HubManagement) GetRouteStatus(domain string) RouteStatus {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	if status, exists := hm.activeGateways[domain]; exists {
		return status
	}
	return RouteHealthy
}
