package strategy

import (
	"time"
)

type RegistratorConfig struct {
	MaxConcurrent      int
	RequestsPerSec     int
	MinJitterMs        int64
	MaxJitterMs        int64
	DriftInterval      time.Duration
	MaxReqsBeforeDrift int64
}

func DefaultConfig() RegistratorConfig {
	return RegistratorConfig{
		MaxConcurrent:      4,
		RequestsPerSec:     2,
		MinJitterMs:        200,
		MaxJitterMs:        800,
		DriftInterval:      10 * time.Minute,
		MaxReqsBeforeDrift: 50,
	}
}
