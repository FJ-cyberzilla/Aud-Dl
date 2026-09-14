package strategy

import (
	"time"

	"audio-command-center/internal/cache"
)

type AntibotRegistrator struct {
	config RegistratorConfig
	cache  *cache.CacheManager
}

func NewAntibotRegistrator(c *cache.CacheManager, cfg RegistratorConfig) *AntibotRegistrator {
	return &AntibotRegistrator{
		cache:  c,
		config: cfg,
	}
}

// BuildConductor initializes all strategy sub-components and returns a fully wired Conductor
func (ar *AntibotRegistrator) BuildConductor() *AntibotConductor {
	tactics := NewAntibotTactics()
	shifter := NewBehaviorShifter(ar.config.MinJitterMs, ar.config.MaxJitterMs)
	limiter := NewGracefulLimiter(ar.config.MaxConcurrent, ar.config.RequestsPerSec)
	dissector := NewAntibotDissection()
	defuser := NewAntibotDefuser()
	drifter := NewAntibotDrifter(ar.config.DriftInterval, ar.config.MaxReqsBeforeDrift)
	fallback := NewFallbackTactics(RetryConfig{
		MaxAttempts: 3,
		InitialWait: 1 * time.Second,
		MaxWait:     5 * time.Second,
	})
	return NewAntibotConductor(
		tactics,
		shifter,
		limiter,
		dissector,
		defuser,
		drifter,
		fallback,
	)
}
