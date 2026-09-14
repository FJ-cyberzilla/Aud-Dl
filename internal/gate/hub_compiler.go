package gate

import (
	"fmt"
	"time"
)

type GateRuntimeConfig struct {
	MaxGlobalConns int
	MaxPerDomain   int
	Timeout        time.Duration
	EnableQUIC     bool
	BanDuration    time.Duration
	MaxOffenses    int
}

type HubCompiler struct{}

func NewHubCompiler() *HubCompiler {
	return &HubCompiler{}
}

// CompileDefaults constructs an industry-standard secure baseline configuration
func (hc *HubCompiler) CompileDefaults() (*GateRuntimeConfig, error) {
	config := &GateRuntimeConfig{
		MaxGlobalConns: 64,
		MaxPerDomain:   12,
		Timeout:        30 * time.Second,
		EnableQUIC:     true,
		BanDuration:    20 * time.Minute,
		MaxOffenses:    3,
	}
	if err := hc.Validate(config); err != nil {
		return nil, fmt.Errorf("gate configuration compilation failed: %w", err)
	}
	return config, nil
}

// Validate ensures structural soundness of the compiled network policy
func (hc *HubCompiler) Validate(cfg *GateRuntimeConfig) error {
	if cfg.MaxGlobalConns <= 0 {
		return fmt.Errorf("MaxGlobalConns must be greater than zero")
	}
	if cfg.MaxPerDomain <= 0 {
		return fmt.Errorf("MaxPerDomain must be greater than zero")
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("Timeout must be a positive duration")
	}
	return nil
}
