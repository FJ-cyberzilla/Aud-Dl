package gate

import (
	"fmt"
	"net/http"
	"sync"
)

type InterceptorFunc func(req *http.Request) error

type HubRegistrator struct {
	mu           sync.RWMutex
	interceptors map[string][]InterceptorFunc
	handlers     map[string]http.Handler
}

func NewHubRegistrator() *HubRegistrator {
	return &HubRegistrator{
		interceptors: make(map[string][]InterceptorFunc),
		handlers:     make(map[string]http.Handler),
	}
}

// RegisterInterceptor binds a custom middleware function to a specific domain or wildcard route
func (hr *HubRegistrator) RegisterInterceptor(domainPattern string, fn InterceptorFunc) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	hr.interceptors[domainPattern] = append(hr.interceptors[domainPattern], fn)
}

// ExecuteInterceptors runs all registered safety hooks for an outgoing request
func (hr *HubRegistrator) ExecuteInterceptors(domain string, req *http.Request) error {
	hr.mu.RLock()
	defer hr.mu.RUnlock()

	// Run wildcard/global interceptors
	if fns, exists := hr.interceptors["*"]; exists {
		for _, fn := range fns {
			if err := fn(req); err != nil {
				return fmt.Errorf("global interceptor blocked request: %w", err)
			}
		}
	}

	// Run domain-specific interceptors
	if fns, exists := hr.interceptors[domain]; exists {
		for _, fn := range fns {
			if err := fn(req); err != nil {
				return fmt.Errorf("domain [%s] interceptor blocked request: %w", domain, err)
			}
		}
	}

	return nil
}
