package strategy

import (
	"context"
	"net/http"
)

type AntibotConductor struct {
	tactics   *AntibotTactics
	shifter   *BehaviorShifter
	limiter   *GracefulLimiter
	dissector *AntibotDissection
	defuser   *AntibotDefuser
	drifter   *AntibotDrifter
	fallback  *FallbackTactics
}

func NewAntibotConductor(
	tactics *AntibotTactics,
	shifter *BehaviorShifter,
	limiter *GracefulLimiter,
	dissector *AntibotDissection,
	defuser *AntibotDefuser,
	drifter *AntibotDrifter,
	fallback *FallbackTactics,
) *AntibotConductor {
	return &AntibotConductor{
		tactics:   tactics,
		shifter:   shifter,
		limiter:   limiter,
		dissector: dissector,
		defuser:   defuser,
		drifter:   drifter,
		fallback:  fallback,
	}
}

// ConductRequest handles pacing, header alignment, session drift, and retries in one call
func (ac *AntibotConductor) ConductRequest(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	if err := ac.limiter.AcquireSlot(ctx); err != nil {
		return nil, err
	}
	defer ac.limiter.ReleaseSlot()
	if err := ac.shifter.JitterWait(ctx); err != nil {
		return nil, err
	}
	if ac.drifter.ShouldDrift() {
		ac.drifter.ApplyDrift(req, client)
	}
	ac.tactics.InjectEvasionHeaders(req)
	ac.drifter.Increment()
	var resp *http.Response
	err := ac.fallback.ExecuteWithRetry(ctx, func(execCtx context.Context) error {
		var execErr error
		resp, execErr = client.Do(req)
		if execErr != nil {
			return execErr
		}
		analysis := ac.dissector.InspectResponse(resp)
		if analysis.Type != ChallengeNone {
			defuseRes, defuseErr := ac.defuser.DefuseChallenge(execCtx, req.URL.String(), analysis.Type)
			if defuseErr != nil || !defuseRes.Resolved {
				return execErr
			}
		}
		return nil
	})
	return resp, err
}
