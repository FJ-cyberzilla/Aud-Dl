package strategy

import (
	"context"
	"fmt"
	"net/http"
)

type AntibotTacticsInterface interface {
	InjectEvasionHeaders(req *http.Request)
}

type BehaviorShifterInterface interface {
	JitterWait(ctx context.Context) error
}

type GracefulLimiterInterface interface {
	AcquireSlot(ctx context.Context) error
	ReleaseSlot()
}

type AntibotDissectionInterface interface {
	InspectResponse(resp *http.Response) DissectionResult
}

type FallbackTacticsInterface interface {
	ExecuteWithRetry(ctx context.Context, operation func(context.Context) error) error
}

type AntibotCommander struct {
	tactics   AntibotTacticsInterface
	shifter   BehaviorShifterInterface
	limiter   GracefulLimiterInterface
	dissector AntibotDissectionInterface
	fallback  FallbackTacticsInterface
}

func NewAntibotCommander(
	tactics AntibotTacticsInterface,
	shifter BehaviorShifterInterface,
	limiter GracefulLimiterInterface,
	dissector AntibotDissectionInterface,
	fallback FallbackTacticsInterface,
) *AntibotCommander {
	return &AntibotCommander{
		tactics:   tactics,
		shifter:   shifter,
		limiter:   limiter,
		dissector: dissector,
		fallback:  fallback,
	}
}

// PrepareAndExecute applies rate slots, headers, and jitter before executing the network request
func (ac *AntibotCommander) PrepareAndExecute(ctx context.Context, req *http.Request, client *http.Client) (*http.Response, error) {
	if err := ac.limiter.AcquireSlot(ctx); err != nil {
		return nil, err
	}
	defer ac.limiter.ReleaseSlot()

	if err := ac.shifter.JitterWait(ctx); err != nil {
		return nil, err
	}

	ac.tactics.InjectEvasionHeaders(req)

	var resp *http.Response
	err := ac.fallback.ExecuteWithRetry(ctx, func(execCtx context.Context) error {
		var execErr error
		resp, execErr = client.Do(req)
		if execErr != nil {
			return execErr
		}

		// Inspect response via dissection engine
		analysis := ac.dissector.InspectResponse(resp)
		if analysis.Type != ChallengeNone {
			return fmt.Errorf("antibot challenge detected: status=%d type=%d", analysis.StatusCode, analysis.Type)
		}
		return nil
	})
	return resp, err
}

