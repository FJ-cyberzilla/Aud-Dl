package strategy

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrAllTacticsFailed = errors.New("strategy: all fallback execution attempts failed")

type RetryConfig struct {
	MaxAttempts int
	InitialWait time.Duration
	MaxWait     time.Duration
}

type FallbackTactics struct {
	config RetryConfig
}

func NewFallbackTactics(config RetryConfig) *FallbackTactics {
	return &FallbackTactics{config: config}
}

// ExecuteWithRetry attempts an operation with exponential backoff and jitter
func (ft *FallbackTactics) ExecuteWithRetry(ctx context.Context, operation func(ctx context.Context) error) error {
	wait := ft.config.InitialWait
	for attempt := 1; attempt <= ft.config.MaxAttempts; attempt++ {
		err := operation(ctx)
		if err == nil {
			return nil
		}
		if attempt == ft.config.MaxAttempts {
			return fmt.Errorf("%w (attempts: %d, last error: %v)", ErrAllTacticsFailed, attempt, err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
		wait *= 2
		if wait > ft.config.MaxWait {
			wait = ft.config.MaxWait
		}
	}
	return ErrAllTacticsFailed
}
